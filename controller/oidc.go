package controller

import "github.com/dgrijalva/jwt-go"

var idTokenKey = []byte("idtoken秘钥")

type IdTokenClaims struct {
	jwt.StandardClaims
}

func OidcGenIdToken(scope, clientId string) (string, error) {
	return "", nil
}
