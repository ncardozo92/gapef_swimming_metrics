package user

import (
	"log"
	"net/http"
	"strconv"

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
	GetAllUsers(context echo.Context) error
	Create(context echo.Context) error
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

// Get all users in the database
func (handler UserHandler) GetAllUsers(context echo.Context) error {
	page, _ := strconv.Atoi(context.QueryParams().Get("page"))
	size, _ := strconv.Atoi(context.QueryParams().Get("size"))

	usersDTOs := []DTO{}

	users, getUsersErr := handler.userRepository.GetUsers(int64(page), int64(size))

	if getUsersErr != nil {
		context.JSON(http.StatusInternalServerError, custom_error.DTO{Message: "No se pudo recuperar los usuarios"})
	}

	for _, user := range users {
		usersDTOs = append(usersDTOs, toDTO(user))
	}

	return context.JSON(http.StatusOK, usersDTOs)
}

func (handler UserHandler) Create(context echo.Context) error {

	dto := DTO{}

	if bindingErr := context.Bind(&dto); bindingErr != nil {
		return context.JSON(http.StatusBadRequest, custom_error.DTO{Message: "El DTO no es válido"})
	}

	entity := fromDTO(dto)

	// Storing hashed password
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)

	if hashingErr != nil {
		return hashingErr
	}

	entity.Password = string(hashedPassword)

	// Storing the entity into the Database collection
	saveErr := handler.userRepository.Create(entity)

	if saveErr != nil {
		return context.JSON(http.StatusInternalServerError, custom_error.DTO{Message: "No se pudo guardar el usuario en la DB"})
	}

	return context.JSON(http.StatusCreated, "")
}

// Returns a new instance of UserHandler
func NewUserHandler(userRepository Repository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}
