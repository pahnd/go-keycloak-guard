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
	a.Kong.Log.Info("Trying to obtain access token from header")
	headers, err := a.Kong.Request.GetHeaders(-1)
	if err != nil {
		return "", err
	}
	normalizedHeaders := make(map[string]string)
	for k, v := range headers {
		normalizedHeaders[strings.ToLower(k)] = strings.Join(v, ", ")
	}
	accessToken, ok := normalizedHeaders["authorization"]
	if !ok {
		a.Kong.Log.Err("Access token not found in header.")
		return "", errors.New("the authorization header is missing")
	}
	a.Kong.Log.Info("Access token found in header.")
	return accessToken, nil
}

func (a *Auth) VerifyAuth() (*dto.Introspect, error) {
	a.Kong.Log.Info("Verifying access token.")
	accessToken, err := a.GetAccessTokenFromHeader()
	if err != nil {
		return nil, err
	}
	a.Kong.Log.Info("Introspecting token..")
	introspectedToken, err := a.IAMClient.Introspect(accessToken, contract.TokenTypeHintAccess)
	if err != nil {
		a.Kong.Log.Err("Failed to introspect token")
		return nil, errors.New(fmt.Sprintf("Failed to fetch instrospected token from Keycloak. Reason: %s", err.Error()))
	}
	a.Kong.Log.Info("Token introspection completed successfully.")
	return introspectedToken, nil
}

func (a *Auth) VerifyUMA() error {
	a.Kong.Log.Info("Verifying UMA Permissions")
	accessToken, err := a.GetAccessTokenFromHeader()
	if err != nil {
		return err
	}
	a.Kong.Log.Info("Fetching UMA permissions.")
	umaPermissions, err := a.IAMClient.GetUMA(accessToken, "")
	if err != nil {
		a.Kong.Log.Err(fmt.Sprintf("Failed to fetch UMA permissions. Reason %s", err.Error()))
		return err
	}
	a.Kong.Log.Info(fmt.Sprintf("Checking if user permissions are valid. Strategy: `%a`.  Permissions: %v", a.PluginConfig.Strategy, a.PluginConfig.ResourceIDs))
	hasPermissions, err := umaPermissions.HasPermission(a.PluginConfig.Permissions, a.PluginConfig.Strategy)
	if err != nil {
		a.Kong.Log.Err(fmt.Sprintf("Failed to verify permissions. Reason: %s", err.Error()))
		return err
	}

	if !hasPermissions {
		a.Kong.Log.Err(fmt.Sprintf("Insufficient permissions. Strategy: `%a`.  Permissions: %v", a.PluginConfig.Strategy, a.PluginConfig.ResourceIDs))
		return errors.New("insufficient permissions")
	}

	return nil
}

func (a *Auth) VerifyRPT() (*response.MissingPermissionTicketResponse, error) {
	a.Kong.Log.Info("Verifying RPT.")
	accessToken, err := a.GetAccessTokenFromHeader()
	if err != nil {
		a.Kong.Log.Info(fmt.Sprintf("Failed to read access token. Requesting permission ticket. Reason: %s", err.Error()))
		permissionTicket, err := a.IAMClient.RequestPermissionTicket(a.PluginConfig.ResourceIDs)
		if err != nil {
			a.Kong.Log.Err(fmt.Sprintf("Failed to request a new permission ticket. Reason: %s", err.Error()))
			return nil, err
		}
		a.Kong.Log.Info("Permission ticket successfully generated.")
		return response.NewRPTResponse(permissionTicket), nil
	}
	a.Kong.Log.Info("Introspecting RPT token.")
	introspectedToken, err := a.IAMClient.Introspect(accessToken, contract.TokenTypeHintRPT)
	if err != nil {
		a.Kong.Log.Info(fmt.Sprintf("Failed to introspect RPT. Requesting permission ticket. Reason: %s", err.Error()))
		permissionTicket, err := a.IAMClient.RequestPermissionTicket(a.PluginConfig.ResourceIDs)
		if err != nil {
			a.Kong.Log.Err(fmt.Sprintf("Failed to request a new permission ticket. Reason: %s", err.Error()))
			return nil, err
		}
		a.Kong.Log.Info("Permission ticket successfully generated.")
		return response.NewRPTResponse(permissionTicket), nil
	}
	if !introspectedToken.IsRPT() {
		a.Kong.Log.Err("Introspected token is not a valid RPT. Proceeding to request a new permission ticket.")
		permissionTicket, err := a.IAMClient.RequestPermissionTicket(a.PluginConfig.ResourceIDs)
		if err != nil {
			a.Kong.Log.Err(fmt.Sprintf("Failed to request a new permission ticket. Reason: %s", err.Error()))
			return nil, err
		}
		a.Kong.Log.Info("Permission ticket successfully generated.")
		return response.NewRPTResponse(permissionTicket), nil
	}
	return nil, nil
}
