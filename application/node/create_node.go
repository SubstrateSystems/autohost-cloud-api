package node

import (
	"context"
	"fmt"

	"github.com/arturo/autohost-cloud-api/internal/domain"
)

type CreateNodeCommand struct {
	HostName     string
	IPLocal      string
	OS           string
	Arch         string
	VersionAgent string
	OwnerID      *string
}

type CreateNodeUseCase struct {
	repository domain.NodeRepository
}

func NewCreateNodeUseCase(repo domain.NodeRepository) *CreateNodeUseCase {
	return &CreateNodeUseCase{
		repository: repo,
	}
}
func (uc *CreateNodeUseCase) Handle(ctx context.Context, cmd CreateNodeCommand) error {
	// 1) Validaciones del caso de uso (no del dominio)
	if cmd.HostName == "" {
		return fmt.Errorf("host name cannot be empty")
	}

	if cmd.IPLocal == "" {
		return fmt.Errorf("ip_local cannot be empty")
	}

	// 2) Crear la entidad de dominio
	node := domain.NewNode(
		cmd.HostName,
		cmd.IPLocal,
		cmd.OS,
		cmd.Arch,
		cmd.VersionAgent,
		cmd.OwnerID,
	)

	// 3) Guardar en el repositorio
	if err := uc.repository.Save(ctx, node); err != nil {
		return err
	}

	return nil
}
