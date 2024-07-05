package auth

import (
	"errors"
	"fmt"
	"github.com/Kong/go-pdk"
	"keycloak-guard/domain/response"
	"keycloak-guard/port/contract"
	"keycloak-guard/port/dto"
	"strings"
)

type Auth struct {
	IAMClient    contract.IAMInterface
	Kong         *pdk.PDK
	PluginConfig *dto.ConfigDTO
}

func New(iamClient contract.IAMInterface, kong *pdk.PDK, pluginConfig *dto.ConfigDTO) Auth {
	return Auth{
		IAMClient:    iamClient,
		Kong:         kong,
		PluginConfig: pluginConfig,
	}
}

func (a *Auth) GetAccessTokenFromHeader() (string, error) {
	headers, err := a.Kong.Request.GetHeaders(-1)
	if err != nil {
		return "", err
	}
	normalizedHeaders := make(map[string]string)
	for k, v := range headers {
		// In case you were wondering and I know some of you reading this are.
		// The http/2 protocol requires the header keys to be lower-case
		// Also Kong normalizes them to lowercase by default.
		// This is just a measure of precaution in case Kong every decides to change this.
		normalizedHeaders[strings.ToLower(k)] = strings.Join(v, ", ")
	}
	accessToken, ok := normalizedHeaders["authorization"]
	if !ok {
		return "", errors.New("the authorization header is missing")
	}
	return accessToken, nil
}

func (a *Auth) VerifyAuth() (*dto.Introspect, error) {
	accessToken, err := a.GetAccessTokenFromHeader()
	if err != nil {
		return nil, err
	}
	introspectedToken, err := a.IAMClient.Introspect(accessToken, contract.TokenTypeHintAccess)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to fetch instrospected token from Keycloak. Reason: %s", err.Error()))
	}
	return introspectedToken, nil
}

func (a *Auth) VerifyUMA() error {
	accessToken, err := a.GetAccessTokenFromHeader()
	if err != nil {
		return err
	}
	umaPermissions, err := a.IAMClient.GetUMA(accessToken, "")
	if err != nil {
		return err
	}

	hasPermissions, err := umaPermissions.HasPermission(a.PluginConfig.Permissions, a.PluginConfig.Strategy)
	if err != nil {
		return err
	}

	if !hasPermissions {
		return errors.New("insufficient permissions")
	}

	return nil
}

func (a *Auth) VerifyRPT() (*response.MissingPermissionTicketResponse, error) {
	accessToken, err := a.GetAccessTokenFromHeader()
	if err != nil {
		permissionTicket, err := a.IAMClient.RequestPermissionTicket(a.PluginConfig.ResourceIDs)
		if err != nil {
			return nil, err
		}
		return response.NewRPTResponse(permissionTicket), nil
	}
	introspectedToken, err := a.IAMClient.Introspect(accessToken, contract.TokenTypeHintRPT)
	if err != nil {
		permissionTicket, err := a.IAMClient.RequestPermissionTicket(a.PluginConfig.ResourceIDs)
		if err != nil {
			return nil, err
		}
		return response.NewRPTResponse(permissionTicket), nil
	}
	if !introspectedToken.IsRPT() {
		permissionTicket, err := a.IAMClient.RequestPermissionTicket(a.PluginConfig.ResourceIDs)
		if err != nil {
			return nil, err
		}
		return response.NewRPTResponse(permissionTicket), nil
	}
	return nil, nil
}
