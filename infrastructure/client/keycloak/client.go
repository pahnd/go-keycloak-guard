package keycloak

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"keycloak-guard/infrastructure/client/keycloak/permission"
	"keycloak-guard/port/dto"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Url          string
	Realm        string
	ClientID     string
	ClientSecret string
}

func New(url, realm, clientID, clientSecret string) *Client {
	return &Client{
		Url:          url,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (c *Client) formatToken(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		return token[7:]
	}
	return token
}

func (c *Client) Introspect(token, tokenTypeHint string) (*dto.Introspect, error) {
	formattedToken := c.formatToken(token)

	formData := url.Values{
		"client_id":       {c.ClientID},
		"client_secret":   {c.ClientSecret},
		"token":           {formattedToken},
		"token_type_hint": {tokenTypeHint},
	}
	req, err := http.NewRequest("POST", c.Url+"/realms/"+c.Realm+"/protocol/openid-connect/token/introspect", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, fmt.Errorf("error requesting token introspection: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if resp == nil {
			return nil, errors.New("failed to introspect token")
		}
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to introspect token. Reason: %s", string(body))
	}

	defer resp.Body.Close()
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.New("failed to introspect token")
		}
		return nil, fmt.Errorf("failed to introspect token. Reason: %s", string(body))
	}
	var introspect dto.Introspect
	err = json.NewDecoder(resp.Body).Decode(&introspect)
	if err != nil {
		return nil, fmt.Errorf("error decoding introspection response: %w", err)
	}
	if introspect.Active == false {
		return nil, fmt.Errorf("token is not active")
	}
	introspect.AccessToken = token
	return &introspect, nil
}
func (c *Client) GetUMA(token string, permissions string) (*permission.PermissionCollection, error) {
	endpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.Url, c.Realm)

	formData := url.Values{
		"grant_type":    {"urn:ietf:params:oauth:grant-type:uma-ticket"},
		"response_mode": {"permissions"},
		"audience":      {c.ClientID},
	}

	if permissions != "" {
		formData.Set("permission", permissions)
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if resp == nil {
			return nil, errors.New("failed to fetch uma permissions")
		}
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch uma permissions. Reason %s", string(body))
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("keycloak error: %s, response body: %s", resp.Status, string(bodyBytes))
	}

	var permissionCollection permission.PermissionCollection
	err = json.NewDecoder(resp.Body).Decode(&permissionCollection)
	if err != nil {
		return nil, fmt.Errorf("error decoding permissions response: %w", err)
	}

	return &permissionCollection, nil
}

func (c *Client) GetClientCredentialsToken() (string, error) {
	endpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.Url, c.Realm)

	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if resp == nil {
			return "", fmt.Errorf("failed to fetch client credentials token")
		}
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch client credentials token. Reason: %s", string(body))
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("keycloak error: %s, response body: %s", resp.Status, string(bodyBytes))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding token response: %w", err)
	}

	return tokenResponse.AccessToken, nil
}

func (c *Client) RequestPermissionTicket(resourceIDs []string) (string, error) {
	endpoint := fmt.Sprintf("%s/realms/%s/authz/protection/permission", c.Url, c.Realm)

	var requestBody []map[string]string
	for _, resourceID := range resourceIDs {
		requestBody = append(requestBody, map[string]string{"resource_id": resourceID})
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(string(jsonBody)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	token, err := c.GetClientCredentialsToken()
	if err != nil {
		return "", fmt.Errorf("failed to fetch client credentials token. Reason: %s", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+c.formatToken(token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("keycloak error: %s, response body: %s", resp.Status, string(bodyBytes))
	}

	var permissionTicketResponse struct {
		Ticket string `json:"ticket"`
	}
	err = json.NewDecoder(resp.Body).Decode(&permissionTicketResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return permissionTicketResponse.Ticket, nil
}
