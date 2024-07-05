package response

import "keycloak-guard/port/dto"

type Headers struct {
	AccessToken string
	Username    string
}

func FromIntrospectedToken(introspectedToken dto.Introspect) *Headers {
	return &Headers{
		AccessToken: "Bearer " + introspectedToken.AccessToken,
		Username:    introspectedToken.Username,
	}
}
