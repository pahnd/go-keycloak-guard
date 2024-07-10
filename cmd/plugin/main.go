package plugin

import (
	"fmt"
	"github.com/Kong/go-pdk"
	"keycloak-guard/domain/auth"
	"keycloak-guard/domain/response"
	"keycloak-guard/infrastructure/client/keycloak"
	"keycloak-guard/port/contract"
	"keycloak-guard/port/dto"
)

type Config struct {
	KeycloakURL                  string
	Realm                        string
	ClientID                     string
	ClientSecret                 string
	EnableAuth                   bool
	EnableUMAAuthorization       bool
	EnableRPTAuthorization       bool
	EnableRoleBasedAuthorization bool
	Permissions                  []string
	Strategy                     string
	ResourceIDs                  []string
	Role                         string
}

func New() interface{} {
	return &Config{}
}

func (c *Config) Access(kong *pdk.PDK) {
	kong.Log.Info("Plugin keycloak-guard started!")
	kc := keycloak.New(c.KeycloakURL, c.Realm, c.ClientID, c.ClientSecret)
	a := auth.New(kc, kong, c.ToConfigDTO())
	if !c.EnableAuth {
		return
	}
	if c.EnableRoleBasedAuthorization && (c.EnableRPTAuthorization || c.EnableUMAAuthorization) {
		r := response.NewErrorResponse("[KeycloakGuard] Conflict in configuration: RoleCheck cannot be used together with RPT or UMA authorization workflows.", 400)
		kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}
	introspectedToken, err := a.VerifyAuth()
	if err != nil && !c.EnableRPTAuthorization {
		r := response.NewErrorResponse(err.Error(), 401)
		kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}
	if c.EnableRoleBasedAuthorization {
		kong.Log.Info(fmt.Sprintf("Verifying if the user has role %s", c.Role))
		if !introspectedToken.HasRole(c.ClientID, c.Role) {
			msg := fmt.Sprintf("User does not have the required role - %s", c.Role)
			kong.Log.Err(msg)
			r := response.NewErrorResponse(msg, 400)
			kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
			return
		}
	}
	if c.EnableUMAAuthorization && !c.EnableRPTAuthorization {
		kong.Log.Info("WORKFLOW: UMA Authorization selected.")
		err := a.VerifyUMA()
		if err != nil {
			r := response.NewErrorResponse(err.Error(), 401)
			kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
			return
		}
	}

	if !c.EnableUMAAuthorization && c.EnableRPTAuthorization {
		kong.Log.Info("WORKFLOW: RPT Authorization selected.")
		r, err := a.VerifyRPT()
		if r != nil {
			kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
			return
		}
		if err != nil {
			r := response.NewErrorResponse(err.Error(), 401)
			kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
			return
		}
	}

	if c.EnableUMAAuthorization && c.EnableRPTAuthorization {
		kong.Log.Info("WORKFLOW: RPT and UMA Authorization selected.")
		isRPT := false
		accessToken, err := a.GetAccessTokenFromHeader()
		if err == nil {
			introspectedRPT, err := a.IAMClient.Introspect(accessToken, contract.TokenTypeHintRPT)
			if err == nil && introspectedRPT.IsRPT() {
				if introspectedToken == nil {
					kong.Log.Info("Replaced introspected token")
					introspectedToken = introspectedRPT
				}
				isRPT = true
			}
		}
		if isRPT {
			kong.Log.Info("Access token is of type RPT.")
			kong.Log.Info("WORKFLOW: RPT Authorization selected.")
			r, err := a.VerifyRPT()
			if r != nil {
				kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
				return
			}
			if err != nil {
				r := response.NewErrorResponse(err.Error(), 401)
				kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
				return
			}
		}
		if !isRPT {
			kong.Log.Info("Access token is NOT of type RPT.")
			kong.Log.Info("WORKFLOW: UMA Authorization selected.")
			err = a.VerifyUMA()
			if err != nil {
				kong.Log.Info(fmt.Sprintf("Verifying UMA FAILED. Error: %s", err.Error()))
				kong.Log.Info("WORKFLOW: SWITCHING to RPT Authorization Workflow")
				r, err := a.VerifyRPT()
				if r != nil {
					kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
					return
				}
				if err != nil {
					r := response.NewErrorResponse(err.Error(), 401)
					kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
					return
				}
			}

		}
	}

	if introspectedToken == nil {
		r := response.NewErrorResponse("Access token does not exist or couldn't be introspected.", 401)
		kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}

	headers := response.FromIntrospectedToken(*introspectedToken)
	kong.ServiceRequest.SetHeader("Authorization", headers.AccessToken)
	kong.ServiceRequest.SetHeader("X-Username", headers.Username)
	kong.Log.Info("Headers set successfully, continuing to upstream service.")
}

func (c *Config) ToConfigDTO() *dto.ConfigDTO {
	return &dto.ConfigDTO{
		KeycloakURL:            c.KeycloakURL,
		Realm:                  c.Realm,
		ClientID:               c.ClientID,
		ClientSecret:           c.ClientSecret,
		EnableAuth:             c.EnableAuth,
		EnableUMAAuthorization: c.EnableUMAAuthorization,
		EnableRPTAuthorization: c.EnableRPTAuthorization,
		Permissions:            c.Permissions,
		Strategy:               c.Strategy,
		ResourceIDs:            c.ResourceIDs,
	}
}
