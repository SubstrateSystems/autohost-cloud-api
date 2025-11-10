package auth

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const (
	CtxUserID  ctxKey = "userID"
	CtxAgentID ctxKey = "agentID"
)

type TokenResolver interface {
	// Valida JWT de usuario → userID
	ValidateUserJWT(r *http.Request, jwt string) (int64, error)
	// Busca por token_prefix y valida hash → agentID
	ValidateAgentToken(r *http.Request, token string) (string, error)
}

func Bearer(token string) string {
	if !strings.HasPrefix(token, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
}

func AuthMiddleware(resolver TokenResolver, allowAgent, allowUser bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			btok := Bearer(r.Header.Get("Authorization"))
			if btok == "" {
				http.Error(w, "missing bearer", http.StatusUnauthorized)
				return
			}
			// agente
			if strings.HasPrefix(btok, "agt_sk_") {
				if !allowAgent {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
				agentID, err := resolver.ValidateAgentToken(r, btok)
				if err != nil {
					http.Error(w, "invalid agent token", http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), CtxAgentID, agentID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// usuario
			if allowUser {
				userID, err := resolver.ValidateUserJWT(r, btok)
				if err == nil {
					ctx := context.WithValue(r.Context(), CtxUserID, userID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			http.Error(w, "invalid token", http.StatusUnauthorized)
		})
	}
}
