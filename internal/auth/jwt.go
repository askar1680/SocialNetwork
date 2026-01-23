package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return nil, nil
}
