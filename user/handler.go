package user

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ncardozo92/gapef_swimming_metrics/custom_error"
	"golang.org/x/crypto/bcrypt"
)

const (
	MESSAGE_USER_NOT_FOUND     = "ususario no encontrado"
	MESSAGE_INCORRECT_PASSWORD = "la contraseña es incorrecta"
	MESSAGE_INTERNAL_ERROR     = "no pudimos autenticar al usuario"
	MESSAGE_BINDING_ERROR      = "el cuerpo de la respuesta no es válido"
)

type Handler interface {
	Login(context echo.Context) error
}

type UserHandler struct {
	userRepository Repository
}

// Login finds the user and generates the jwt for authorization
func (handler UserHandler) Login(context echo.Context) error {

	// binding the json body
	dto := new(DTO)
	if err := context.Bind(dto); err != nil || !isValidRequest(dto) {
		log.Println(MESSAGE_BINDING_ERROR)
		return context.JSON(http.StatusBadRequest, custom_error.DTO{Message: MESSAGE_BINDING_ERROR})
	}

	// finding the user by his username
	user, findUserErr, userNotFound := handler.userRepository.FindByUsername(dto.Username)

	if findUserErr != nil {
		if userNotFound {
			log.Println(MESSAGE_USER_NOT_FOUND)
			return context.JSON(http.StatusNotFound, custom_error.DTO{Message: MESSAGE_USER_NOT_FOUND})
		}
		return context.JSON(http.StatusNotFound, nil)
	}

	// Comparing the received password with the storaged password
	passwordValidationErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))

	if passwordValidationErr != nil {
		log.Println("the password is incorrect")
		return context.JSON(http.StatusUnauthorized, nil)
	}

	// Generating the jwt
	jwt, jwtGenerationErr := generateJWT(user)

	if jwtGenerationErr != nil {
		log.Println("Cannot generate JWT", jwtGenerationErr)
	}

	return context.JSON(http.StatusOK, LoginDTO{Token: jwt})
}

func isValidRequest(dto *DTO) bool {
	var result bool
	if len(dto.Password) > 0 && len(dto.Username) > 0 {
		result = true
	}

	return result
}

// Returns a new instance of UserHandler
func NewUserHandler(userRepository Repository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}
