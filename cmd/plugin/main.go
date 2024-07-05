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
	kong.Log.Info("Introspecting AUTH Token")
	introspectedToken, err := a.VerifyAuth()
	if err != nil && !c.EnableRPTAuthorization {
		r := response.NewErrorResponse(err.Error(), 401)
		kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}
	kong.Log.Info("DONE Introspecting AUTH Token")
	if c.EnableUMAAuthorization && !c.EnableRPTAuthorization {
		kong.Log.Info("Verifying UMA Authorization")
		err := a.VerifyUMA()
		if err != nil {
			r := response.NewErrorResponse(err.Error(), 401)
			kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
			return
		}
	}

	if !c.EnableUMAAuthorization && c.EnableRPTAuthorization {
		kong.Log.Info("Verifying RPT Authorization")
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
		isRPT := false
		kong.Log.Info("Introspecting RPT Token")
		accessToken, err := a.GetAccessTokenFromHeader()
		if err == nil {
			introspectedRPT, err := a.IAMClient.Introspect(accessToken, contract.TokenTypeHintRPT)
			kong.Log.Info("Finished introspecting RPT Token")
			if err == nil && introspectedRPT.IsRPT() {
				if introspectedToken == nil {
					kong.Log.Info("Replaced introspected token")
					introspectedToken = introspectedRPT
				}
				isRPT = true
			}
		}
		if isRPT {
			kong.Log.Info("Verifying RPT")
			r, err := a.VerifyRPT()
			if r != nil {
				kong.Log.Info("Verifying RPT FAILED")
				kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
				return
			}
			if err != nil {
				kong.Log.Info(fmt.Sprintf("Verifying RPT FAILED. Error: %s", err.Error()))
				r := response.NewErrorResponse(err.Error(), 401)
				kong.Response.Exit(r.Code, r.ToJson(), map[string][]string{"Content-Type": {"application/json"}})
				return
			}
			kong.Log.Info("Done verifying RPT")
		}
		if !isRPT {
			kong.Log.Info("Verifying UMA")
			err = a.VerifyUMA()
			if err != nil {
				kong.Log.Info(fmt.Sprintf("Verifying UMA FAILED. Error: %s", err.Error()))
				kong.Log.Info("Verifying RPT")
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
				kong.Log.Info("Done verifying RPT")
			}

			kong.Log.Info("Done verifying UMA")

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
