package user

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ncardozo92/gapef_swimming_metrics/constants"
	"github.com/ncardozo92/gapef_swimming_metrics/custom_error"
	"github.com/ncardozo92/gapef_swimming_metrics/logging"
)

const (
	JWT_FIELD_ROLE          = "role"
	JWT_FIELD_ID            = "id"
	JWT_BEARER_PREFIX       = "Bearer "
	AUTHORIZATION_HEADER    = "Authorization"
	ISSUER                  = "GAPEF"
	ROLE_TRAINNER           = "TRAINER"
	BEARER_PREFIX           = "Bearer "
	MESSAGE_JWT_NOT_PRESENT = "Debe enviarse un JWT válido"
)

var JWT_SECRET string = os.Getenv("JWT_SECRET")

type AuthenticationClaims struct {
	Role      string `json:"role"`
	UserId    string `json:"user_id"`
	JwtClaims jwt.RegisteredClaims
}

func generateJWT(user Entity) (string, error) {

	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":           user.Id,
		JWT_FIELD_ROLE: user.Role,
		"iss":          ISSUER,
		"sub":          user.Username,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(3 * time.Minute).Unix(),
	})

	return tokenGenerator.SignedString([]byte(JWT_SECRET))
}

func CustomJwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Path() != PATH_LOGIN {
			authenticationHeader := c.Request().Header.Get(AUTHORIZATION_HEADER)

			if authenticationHeader == "" {
				logging.LogError("JWT not present at the request")
				return c.JSON(http.StatusUnauthorized, custom_error.DTO{Message: MESSAGE_JWT_NOT_PRESENT})
			}

			validationErr := validateJWT(strings.Replace(authenticationHeader, JWT_BEARER_PREFIX, "", 1))

			if validationErr != nil {
				logging.LogError("Error validating the JWT: %v", validationErr)
				return c.JSON(http.StatusForbidden, custom_error.DTO{Message: "Debe enviarse un JWT válido"})
			}
		}

		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

func CoachAccessMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorizationHeader := c.Request().Header.Get(AUTHORIZATION_HEADER)
		if authorizationHeader == "" {
			logging.LogError("JWT not present at the request")
			return c.JSON(http.StatusUnauthorized, custom_error.DTO{Message: MESSAGE_JWT_NOT_PRESENT})
		}

		if validationTokenErr := validateJWT(authorizationHeader); validationTokenErr != nil {
			logging.LogError("Could not validate JWT, %s", validationTokenErr)
			return c.JSON(http.StatusUnauthorized, custom_error.DTO{Message: "The JWT is not valid"})
		}

		// we must deserialize the token to retrieve user data
		encodedPayload := strings.Split(authorizationHeader, ".")[1]
		rawPayload := make(map[string]any)

		// Now we decode the payload
		decodedPayload, decodingErr := base64.RawURLEncoding.DecodeString(encodedPayload)

		if decodingErr != nil {
			logging.LogError("Could not decode JWT payload")
			return c.JSON(http.StatusInternalServerError, custom_error.DTO{Message: "Cannot authorize user"})
		}

		unmarshallErr := json.Unmarshal(decodedPayload, &rawPayload)

		if unmarshallErr != nil {
			logging.LogError("Could not unmarshall JWT payload")
			return c.JSON(http.StatusInternalServerError, custom_error.DTO{Message: "Cannot authorize user"})
		}

		if userRole, userRolePresent := rawPayload[JWT_FIELD_ROLE]; userRolePresent && userRole == constants.ROLE_COACH {
			if err := next(c); err != nil {
				c.Error(err)
			}
		} else {
			logging.LogError("Role not present at JWT")
			return c.JSON(http.StatusInternalServerError, custom_error.DTO{Message: "User must be Coach"})
		}
		return nil
	}
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(strings.Replace(tokenString, BEARER_PREFIX, "", 1), func(token *jwt.Token) (interface{}, error) {
		if _, okToken := token.Method.(*jwt.SigningMethodHMAC); !okToken {
			return nil, errors.New("the JWT provided haven't got the right signing method")
		}

		return []byte(JWT_SECRET), nil
	})
}

func validateJWT(tokenString string) error {

	parsedToken, parsingErr := parseJWT(tokenString)

	if parsingErr != nil {
		return parsingErr
	}

	issuer, getIssuerErr := parsedToken.Claims.GetIssuer()

	if !parsedToken.Valid || issuer != ISSUER || getIssuerErr != nil {
		return errors.New("the JWT provided is not valid")
	}

	return nil
}
