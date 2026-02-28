package grpcserver

import (
	"context"
	"log"
	"strings"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/arturo/autohost-cloud-api/internal/domain/job"
	nodecommand "github.com/arturo/autohost-cloud-api/internal/domain/node_command"
	nodetoken "github.com/arturo/autohost-cloud-api/internal/domain/node_token"
	pb "github.com/arturo/autohost-cloud-api/internal/grpc/nodepb"
	"github.com/arturo/autohost-cloud-api/internal/platform"
)

// nodeStream holds the send side of an active gRPC Connect stream.
type nodeStream struct {
	send chan *pb.ServerMessage
}

// NodeAgentServer implements pb.NodeAgentServiceServer.
type NodeAgentServer struct {
	pb.UnimplementedNodeAgentServiceServer

	commandSvc *nodecommand.Service
	jobSvc     *job.Service
	tokenSvc   *nodetoken.Service

	streamsMu sync.RWMutex
	streams   map[string]*nodeStream // nodeID -> active stream
}

// NewNodeAgentServer creates a ready-to-register gRPC server.
func NewNodeAgentServer(
	commandSvc *nodecommand.Service,
	jobSvc *job.Service,
	tokenSvc *nodetoken.Service,
) *NodeAgentServer {
	return &NodeAgentServer{
		commandSvc: commandSvc,
		jobSvc:     jobSvc,
		tokenSvc:   tokenSvc,
		streams:    make(map[string]*nodeStream),
	}
}

// ---- Auth helper ------------------------------------------------------------

// nodeIDFromCtx validates the "authorization" metadata and returns the node ID.
func (s *NodeAgentServer) nodeIDFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	raw := values[0]
	if strings.HasPrefix(raw, "Bearer ") {
		raw = raw[7:]
	}

	if !strings.HasPrefix(raw, platform.TokenApiPrefix) {
		return "", status.Error(codes.Unauthenticated, "invalid node token format")
	}

	tokenHash := platform.HashTokenApi(raw)
	tok, err := s.tokenSvc.FindNodeTokenByHash(tokenHash)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid node token")
	}
	if tok.RevokedAt != nil {
		return "", status.Error(codes.Unauthenticated, "node token revoked")
	}
	return tok.NodeID, nil
}

// ---- NodeDispatcher (used by JobHandler) ------------------------------------

// SendToNode pushes a ServerMessage to a connected node.
// Returns an error if the node has no live gRPC Connect stream.
func (s *NodeAgentServer) SendToNode(nodeID string, msg *pb.ServerMessage) error {
	s.streamsMu.RLock()
	ns, ok := s.streams[nodeID]
	s.streamsMu.RUnlock()
	if !ok {
		return status.Errorf(codes.NotFound, "node %s not connected via gRPC", nodeID)
	}
	select {
	case ns.send <- msg:
		return nil
	default:
		return status.Errorf(codes.Unavailable, "node %s send buffer full", nodeID)
	}
}

// ---- RegisterCommands (client-side stream) ----------------------------------

func (s *NodeAgentServer) RegisterCommands(
	stream pb.NodeAgentService_RegisterCommandsServer,
) error {
	nodeID, err := s.nodeIDFromCtx(stream.Context())
	if err != nil {
		return err
	}

	var count int32
	for {
		req, recvErr := stream.Recv()
		if recvErr != nil {
			break // io.EOF or disconnect
		}
		cmd := &nodecommand.NodeCommand{
			NodeID:      nodeID,
			Name:        req.GetName(),
			Description: req.GetDescription(),
			Type:        pbCommandType(req.GetType()),
			ScriptPath:  req.GetScriptPath(),
		}
		if _, err := s.commandSvc.Register(cmd); err != nil {
			log.Printf("[gRPC] register command %s for node %s: %v", cmd.Name, nodeID, err)
			continue
		}
		count++
	}

	log.Printf("[gRPC] node %s registered %d commands", nodeID, count)
	return stream.SendAndClose(&pb.RegisterCommandsResponse{Registered: count})
}

// ---- Connect (bidirectional stream) -----------------------------------------

func (s *NodeAgentServer) Connect(
	stream pb.NodeAgentService_ConnectServer,
) error {
	nodeID, err := s.nodeIDFromCtx(stream.Context())
	if err != nil {
		return err
	}

	ns := &nodeStream{send: make(chan *pb.ServerMessage, 64)}
	s.register(nodeID, ns)
	defer func() {
		s.unregister(nodeID)
		close(ns.send)
		log.Printf("[gRPC] node %s disconnected", nodeID)
	}()

	log.Printf("[gRPC] node %s connected", nodeID)

	// Forward queued ServerMessages to the stream.
	sendErr := make(chan error, 1)
	go func() {
		for msg := range ns.send {
			if err := stream.Send(msg); err != nil {
				sendErr <- err
				return
			}
		}
		sendErr <- nil
	}()

	for {
		// Check if the send goroutine hit an error.
		select {
		case err := <-sendErr:
			return err
		default:
		}

		in, recvErr := stream.Recv()
		if recvErr != nil {
			break
		}

		switch p := in.Payload.(type) {
		case *pb.NodeMessage_JobResult:
			r := p.JobResult
			if err := s.jobSvc.UpdateResult(
				r.GetJobId(),
				job.JobStatus(pbJobStatus(r.GetStatus())),
				r.GetOutput(),
				r.GetError(),
			); err != nil {
				log.Printf("[gRPC] update job %s: %v", r.GetJobId(), err)
			} else {
				log.Printf("[gRPC] job %s -> %s (node %s)", r.GetJobId(), r.GetStatus(), nodeID)
			}

		case *pb.NodeMessage_Heartbeat:
			log.Printf("[gRPC] heartbeat from node %s", p.Heartbeat.GetNodeId())

		default:
			log.Printf("[gRPC] unknown payload from node %s", nodeID)
		}
	}

	return nil
}

// ---- Registry ---------------------------------------------------------------

func (s *NodeAgentServer) register(nodeID string, ns *nodeStream) {
	s.streamsMu.Lock()
	defer s.streamsMu.Unlock()
	s.streams[nodeID] = ns
}

func (s *NodeAgentServer) unregister(nodeID string) {
	s.streamsMu.Lock()
	defer s.streamsMu.Unlock()
	delete(s.streams, nodeID)
}

// ---- Proto enum converters --------------------------------------------------

func pbCommandType(t pb.CommandType) nodecommand.CommandType {
	if t == pb.CommandType_COMMAND_TYPE_CUSTOM {
		return nodecommand.CommandTypeCustom
	}
	return nodecommand.CommandTypeDefault
}

func pbJobStatus(s pb.JobStatus) string {
	switch s {
	case pb.JobStatus_JOB_STATUS_RUNNING:
		return string(job.StatusRunning)
	case pb.JobStatus_JOB_STATUS_FAILED:
		return string(job.StatusFailed)
	default:
		return string(job.StatusCompleted)
	}
}

// ---- DispatchJob (implements handler.NodeDispatcher) -----------------------

// DispatchJob builds a gRPC ServerMessage and pushes it to the connected node.
func (s *NodeAgentServer) DispatchJob(
	nodeID, jobID, commandName string,
	commandType nodecommand.CommandType,
) error {
	ct := pb.CommandType_COMMAND_TYPE_DEFAULT
	if commandType == nodecommand.CommandTypeCustom {
		ct = pb.CommandType_COMMAND_TYPE_CUSTOM
	}
	msg := &pb.ServerMessage{
		Payload: &pb.ServerMessage_ExecuteJob{
			ExecuteJob: &pb.ExecuteJobPayload{
				JobId:       jobID,
				CommandName: commandName,
				CommandType: ct,
			},
		},
	}
	return s.SendToNode(nodeID, msg)
}
