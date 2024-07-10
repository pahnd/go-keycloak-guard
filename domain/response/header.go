package response

import (
	"keycloak-guard/port/dto"
	"strings"
)

type Headers struct {
	AccessToken string
	Username    string
}

func formatToken(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		return token[7:]
	}
	return token
}

func FromIntrospectedToken(introspectedToken dto.Introspect) *Headers {
	return &Headers{
		AccessToken: "Bearer " + formatToken(introspectedToken.AccessToken),
		Username:    introspectedToken.Username,
	}
}
