package lib

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

var (
	// ErrTokenIsMalformed is returned when the provided token for
	// validation is malformed
	ErrTokenIsMalformed = errors.New("Token is malformed")

	// ErrTokenHasExpired is returned when the provided token for
	// validation has expired
	ErrTokenHasExpired = errors.New("Token has expired")

	// ErrTokenIsInvalid is returned when the provided token for
	// validation is invalid
	ErrTokenIsInvalid = errors.New("Token is invalid")
)

// Encode is used to encode JWT tokens
func Encode(sub, email string) string {
	now := time.Now().UTC()
	expires := now.Add(time.Hour * 168) // Expires in one week

	encodeToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   "Wishlist",
		"isa":   now.Unix(),
		"exp":   expires.Unix(),
		"sub":   sub,
		"email": email,
	})

	token, err := encodeToken.SignedString(privateKey)
	if err != nil {
		LogError(LPanic, "Could not encode token", err)
	}

	return token
}

// Decode is used to decode JWT tokens
func Decode(tokenString string) (jwt.MapClaims, bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, token.Valid, nil
	}

	return nil, false, errors.New("error: Token is invalid")
}

// IsMalformed reports whether err represents a token is malformed error
func IsMalformed(err error) bool {
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&jwt.ValidationErrorMalformed != 0
	}

	return false
}

// HasExpired reports whether err represents a token has expired error
func HasExpired(err error) bool {
	if ve, ok := err.(*jwt.ValidationError); ok {
		return ve.Errors&jwt.ValidationErrorExpired != 0
	}

	return false
}

func init() {
	prv, err := ioutil.ReadFile("./secrets/private.pem")
	if err != nil {
		LogError(LFatal, "Could not read './secrets/private.pem' file", err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(prv)
	if err != nil {
		LogError(LFatal, "Could not parse './secrets/private.pem'", err)
	}

	pub, err := ioutil.ReadFile("./secrets/public.pem")
	if err != nil {
		LogError(LFatal, "Could not read './secrets/public.pem' file", err)
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		LogError(LFatal, "Could not parse './secrets/public.pem'", err)
	}
}
