package contract

import (
	"keycloak-guard/infrastructure/client/keycloak/permission"
	"keycloak-guard/port/dto"
)

const TokenTypeHintAccess = "access_token"
const TokenTypeHintRPT = "requesting_party_token"

type IAMInterface interface {
	Introspect(token, tokenTypeHint string) (*dto.Introspect, error)
	GetUMA(token string, permissions string) (*permission.PermissionCollection, error)
	GetClientCredentialsToken() (string, error)
	RequestPermissionTicket(resourceIDs []string) (string, error)
}
