package customsturctures

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type AuthClaimers struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

type TokenPair struct {
	Access_Token  string `json:access_token`
	Refresh_Token string `json:access_token`
}
