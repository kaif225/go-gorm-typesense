package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignedDetail struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func SignedToken(firstName, lastName, email, username, role string) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")
	claims := &SignedDetail{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Username:  username,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "kaif",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}
	return signedToken, nil
}
