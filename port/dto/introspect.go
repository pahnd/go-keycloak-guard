package dto

import "strings"

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Roles []string `json:"roles"`
}

type Permission struct {
	ResourceID string   `json:"resource_id"`
	Scopes     []string `json:"scopes"`
}

type Introspect struct {
	Active            bool                      `json:"active"`
	Username          string                    `json:"username,omitempty"`
	Email             string                    `json:"email,omitempty"`
	FirstName         string                    `json:"given_name,omitempty"`
	LastName          string                    `json:"family_name,omitempty"`
	ClientID          string                    `json:"client_id,omitempty"`
	UserID            string                    `json:"sub,omitempty"`
	ExpiresAt         int64                     `json:"exp,omitempty"`
	IssuedAt          int64                     `json:"iat,omitempty"`
	AuthTime          int64                     `json:"auth_time,omitempty"`
	Jti               string                    `json:"jti,omitempty"`
	Iss               string                    `json:"iss,omitempty"`
	Typ               string                    `json:"typ,omitempty"`
	Azp               string                    `json:"azp,omitempty"`
	SessionState      string                    `json:"session_state,omitempty"`
	Name              string                    `json:"name,omitempty"`
	PreferredUsername string                    `json:"preferred_username,omitempty"`
	EmailVerified     bool                      `json:"email_verified,omitempty"`
	Acr               string                    `json:"acr,omitempty"`
	AllowedOrigins    []string                  `json:"allowed-origins,omitempty"`
	RealmAccess       RealmAccess               `json:"realm_access,omitempty"`
	ResourceAccess    map[string]ResourceAccess `json:"resource_access,omitempty"`
	Scope             string                    `json:"scope,omitempty"`
	Sid               string                    `json:"sid,omitempty"`
	AccessToken       string                    `json:"accessToken,omitempty"`
	Permissions       []Permission              `json:"permissions,omitempty"`
}

func (i *Introspect) IsRPT() bool {
	if strings.Contains(i.Scope, "uma_protection") || len(i.Permissions) > 0 {
		return true
	}
	return false
}

func (i *Introspect) HasRole(clientID, role string) bool {
	if roles, ok := i.ResourceAccess[clientID]; ok {
		for _, r := range roles.Roles {
			if r == role {
				return true
			}
		}
	}
	return false

}
