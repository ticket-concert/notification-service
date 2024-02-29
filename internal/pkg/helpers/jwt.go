package helpers

import (
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	"notification-service/internal/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	verifyKey    *rsa.PublicKey
	signKey      *rsa.PrivateKey
	verifyKeyRef *rsa.PublicKey
	signKeyRef   *rsa.PrivateKey
)

type JwtImpl struct{}

type ConfigInitializer interface {
	InitConfig(privateKeyConf string, publicKeyConf string, privateKeyRefConf string, publicKeyRefConf string)
}

func (j *JwtImpl) InitConfig(privateKeyConf string, publicKeyConf string, privateKeyRefConf string, publicKeyRefConf string) {
	var err error

	privateKey := strings.ReplaceAll(privateKeyConf, "\\n", "\n")
	signBytes := []byte(privateKey)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic("Private Key cannot Verify")
	}

	publicKey := strings.ReplaceAll(publicKeyConf, "\\n", "\n")
	verifyBytes := []byte(publicKey)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic("Public Key cannot Verify")
	}

	// Sign Refresh Key
	privateKeyRefresh := strings.ReplaceAll(privateKeyRefConf, "\\n", "\n")
	signBytesRef := []byte(privateKeyRefresh)

	signKeyRef, err = jwt.ParseRSAPrivateKeyFromPEM(signBytesRef)
	if err != nil {
		panic("Private Key cannot Verify")
	}

	publicKeyRefresh := strings.ReplaceAll(publicKeyRefConf, "\\n", "\n")
	verifyBytesRef := []byte(publicKeyRefresh)

	verifyKeyRef, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytesRef)
	if err != nil {
		panic("Public Key cannot Verify")
	}
}

type PayloadJWT struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
	Role   string `json:"role"`
}

const leeway = -120

type MyClaims struct {
	PayloadJWT
	*jwt.StandardClaims
}

func (c *MyClaims) Valid() error {
	c.StandardClaims.ExpiresAt += leeway
	c.StandardClaims.IssuedAt += leeway
	err := c.StandardClaims.Valid()
	c.StandardClaims.ExpiresAt -= leeway
	c.StandardClaims.IssuedAt -= leeway
	return err
}

func (j *JwtImpl) JWTAuthorization(request *http.Request) (*PayloadJWT, error) {
	token := strings.Split(request.Header.Get("authorization"), " ")
	if token[0] == "Bearer" || token[0] == "bearer" {
		authToken := token[1]
		if len(authToken) == 0 {
			return nil, errors.ForbiddenError("Invalid token!")
		}
		var parsedTokenClaims = new(MyClaims)

		_, err := jwt.ParseWithClaims(authToken, parsedTokenClaims, func(authToken *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

		if err != nil {
			return nil, errors.UnauthorizedError("Access token expired!")
		}

		claim := parsedTokenClaims.PayloadJWT

		// Handle Payload Token
		if claim.UserId == "" {
			return nil, errors.ForbiddenError("Invalid token userId")
		}
		return &PayloadJWT{
			UserId: claim.UserId,
			Token:  token[1],
		}, nil
	} else {
		return nil, errors.ForbiddenError("Invalid token")
	}

}

func (j *JwtImpl) JWTRefreshAuthorization(authToken string) (*PayloadJWT, error) {
	var parsedTokenClaims = new(MyClaims)
	_, err := jwt.ParseWithClaims(authToken, parsedTokenClaims, func(authToken *jwt.Token) (interface{}, error) {
		return verifyKeyRef, nil
	})

	if err != nil {
		return nil, errors.UnauthorizedError("Access token expired!")
	}

	claim := parsedTokenClaims.PayloadJWT

	// Handle Payload Token
	if claim.UserId == "" {
		return nil, errors.ForbiddenError("Invalid token userId")
	}

	return &PayloadJWT{
		UserId: claim.UserId,
		Token:  authToken,
	}, nil
}

func (j *JwtImpl) GenerateToken(ttl time.Duration, payload map[string]interface{}) (string, string, error) {
	now := time.Now().UTC()
	expiredAt := now.Add(ttl).Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = expiredAt
	claims["iat"] = now.Unix()
	claims["userId"] = payload["userId"]
	for key, value := range payload {
		claims[key] = value
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(signKey)
	if err != nil {
		return "", "", errors.InternalServerError(err.Error())
	}

	return token, time.Unix(expiredAt, 0).Format(time.RFC3339), nil
}

func (j *JwtImpl) GenerateTokenRefresh(ttl time.Duration, payload map[string]interface{}) (string, error) {
	now := time.Now().UTC()
	expiredAt := now.Add(ttl).Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = expiredAt
	claims["iat"] = now.Unix()
	claims["userId"] = payload["userId"]
	for key, value := range payload {
		claims[key] = value
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(signKeyRef)
	if err != nil {
		return "", errors.InternalServerError(err.Error())
	}

	return token, nil
}

type TokenGenerator interface {
	GenerateToken(ttl time.Duration, payload map[string]interface{}) (string, string, error)
	GenerateTokenRefresh(ttl time.Duration, payload map[string]interface{}) (string, error)
	JWTAuthorization(request *http.Request) (*PayloadJWT, error)
	JWTRefreshAuthorization(authToken string) (*PayloadJWT, error)
}
