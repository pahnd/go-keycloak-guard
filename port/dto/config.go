package dto

type ConfigDTO struct {
	KeycloakURL            string
	Realm                  string
	ClientID               string
	ClientSecret           string
	EnableAuth             bool
	EnableUMAAuthorization bool
	EnableRPTAuthorization bool
	Permissions            []string
	Strategy               string
	ResourceIDs            []string
}
