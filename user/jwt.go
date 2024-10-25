package user

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWT_FIELD_ROLE    = "role"
	JWT_FIELD_ID      = "id"
	JWT_BEARER_PREFIX = "Bearer "
)

var JWT_SECRET string = os.Getenv("JWT_SECRET")

func generateJWT(user User) (string, error) {
	now := time.Now().Unix()
	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss":        "golang-blog",
			"sub":        user.Username,
			JWT_FIELD_ID: user.Id,
			"iat":        now,
			"exp":        now + (5 * 60), // the token is valid for 5 minutes
		})

	return tokenGenerator.SignedString([]byte(JWT_SECRET))
}
