package user

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ncardozo92/gapef_swimming_metrics/constants"
	"github.com/ncardozo92/gapef_swimming_metrics/custom_error"
	"github.com/ncardozo92/gapef_swimming_metrics/logging"
	"golang.org/x/crypto/bcrypt"
)

const (
	MESSAGE_USER_NOT_FOUND        = "ususario no encontrado"
	MESSAGE_INCORRECT_PASSWORD    = "la contraseña es incorrecta"
	MESSAGE_INTERNAL_ERROR        = "no pudimos autenticar al usuario"
	MESSAGE_BINDING_ERROR         = "el formato del cuerpo de la solicitud no es válido"
	MESSAGE_VALIDATION_ERROR      = "la solicitud posee datos inválidos"
	MESSAGE_USER_ALREADY_EXISTS   = "Ya existe un usuario con el mismo username o email"
	MESSAGE_USER_CREATION_ERROR   = "No se ha podido crear al usuario"
	MESSAGE_JWT_NOT_CREATED       = "No se pudo iniciar la sesión"
	MESSAGE_CANNOT_RETRIEVE_USERS = "No se pudo recuperar los usuarios"
	DETAIL_INVALID_EMAIL          = "El email no es válido"
	DETAIL_INVALID_USERNAME       = "El username no puede ser un string vacío"
	DETAIL_INVALID_PASSWORD       = "La password no puede ser un string vacío"
	DETAIL_INVALID_ROLE           = "El rol suministrado no es válido"
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
	if err := context.Bind(dto); err != nil || !isLoginValidRequest(dto) {
		logging.LogError(MESSAGE_BINDING_ERROR)
		return context.JSON(http.StatusBadRequest, custom_error.DTO{Message: MESSAGE_BINDING_ERROR})
	}

	// finding the user by his username
	user, findUserErr, userNotFound := handler.userRepository.FindByUsername(dto.Username)

	if findUserErr != nil {
		if userNotFound {
			logging.LogError("User not found")
			return context.JSON(http.StatusNotFound, custom_error.DTO{Message: MESSAGE_USER_NOT_FOUND})
		}
		return context.JSON(http.StatusNotFound, nil)
	}

	// Comparing the received password with the storaged password
	passwordValidationErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))

	if passwordValidationErr != nil {
		logging.LogError("the password is incorrect")
		return context.JSON(http.StatusUnauthorized, nil)
	}

	// Generating the jwt
	jwt, jwtGenerationErr := generateJWT(user)

	if jwtGenerationErr != nil {
		logging.LogError("Cannot generate JWT %v", jwtGenerationErr)
		return context.JSON(http.StatusInternalServerError, custom_error.DTO{Message: MESSAGE_JWT_NOT_CREATED})
	}

	return context.JSON(http.StatusOK, LoginDTO{Token: jwt})
}

// Validates the login DTO has not blank username and password
func isLoginValidRequest(dto *DTO) bool {
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
		logging.LogError("could not retrieve users from collection")
		context.JSON(http.StatusInternalServerError, custom_error.DTO{Message: MESSAGE_CANNOT_RETRIEVE_USERS})
	}

	for _, user := range users {
		usersDTOs = append(usersDTOs, toDTO(user))
	}

	return context.JSON(http.StatusOK, usersDTOs)
}

func (handler UserHandler) Create(context echo.Context) error {

	dto := DTO{}

	if bindingErr := context.Bind(&dto); bindingErr != nil {
		logging.LogError("DTO is not valid")
		return context.JSON(http.StatusBadRequest, custom_error.DTO{Message: "Los datos no son válidos"})
	}

	var validationErrorDetails []string = validateDTO(dto)

	if len(validationErrorDetails) > 0 {
		logging.LogError("there are DTO validation errors, %v", validationErrorDetails)
		return context.JSON(http.StatusBadRequest,
			custom_error.DTO{Message: MESSAGE_VALIDATION_ERROR, Details: validationErrorDetails})
	}

	entity := fromDTO(dto)

	userExists, findingUserErr := handler.userRepository.Exists(entity)

	if findingUserErr != nil {
		logging.LogError("User could not be created")
		return context.JSON(http.StatusInternalServerError, custom_error.DTO{Message: MESSAGE_USER_CREATION_ERROR})
	}

	if userExists {
		logging.LogWarning("User already exists")
		return context.JSON(http.StatusConflict, custom_error.DTO{Message: MESSAGE_USER_ALREADY_EXISTS})
	}

	// Storing hashed password
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)

	if hashingErr != nil {
		logging.LogError("Error hashing the password, %v", hashingErr)
		return hashingErr
	}

	entity.Password = string(hashedPassword)

	// Storing the entity into the Database collection
	saveErr := handler.userRepository.Create(entity)

	if saveErr != nil {
		logging.LogError("could not save user, %v", saveErr)
		return context.JSON(http.StatusInternalServerError, custom_error.DTO{Message: MESSAGE_USER_CREATION_ERROR})
	}

	return context.NoContent(http.StatusCreated)
}

// Validates that user creation/updating DTO has the right fields
func validateDTO(dto DTO) []string {
	details := []string{}
	emailRegex := regexp.MustCompile(constants.REGEX_EMAIL_VALIDATION)

	if !emailRegex.Match([]byte(dto.Email)) {
		details = append(details, DETAIL_INVALID_EMAIL)
	}

	if len(dto.Username) == 0 {
		details = append(details, DETAIL_INVALID_USERNAME)
	}

	if len(dto.Password) == 0 {
		details = append(details, DETAIL_INVALID_PASSWORD)
	}

	if dto.Role != constants.ROLE_ADMIN && dto.Role != constants.ROLE_ATLETHE && dto.Role != constants.ROLE_COACH {
		details = append(details, DETAIL_INVALID_ROLE)
	}

	return details
}

// Returns a new instance of UserHandler
func NewUserHandler(userRepository Repository) *UserHandler {
	return &UserHandler{userRepository: userRepository}
}
