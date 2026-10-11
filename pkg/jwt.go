package pkg

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	Id   int
	Role string
	jwt.RegisteredClaims
}

func NewJWTClaims(userID int, role string) *JWTClaims {
	ID, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
	}
	return &JWTClaims{
		Id:   userID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			ID: ID.String(),
		},
	}
}

func (j *JWTClaims) GenToken() (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		return "", fmt.Errorf("jwt key not found")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)

	return token.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func (j *JWTClaims) DecodeToken(token string) error {
	jwtToken, err := jwt.ParseWithClaims(token, j, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return err
	}
	if !jwtToken.Valid {
		return jwt.ErrTokenExpired
	}
	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}
	return nil
}
