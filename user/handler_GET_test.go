package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var requestBody string = `{"username": "ncardozo","password":"ncardozo"}`

func TestLoginOK(t *testing.T) {
	// setup the mocks into the SUT
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	foundUser := Entity{
		Id:       "asdf",
		Email:    "ncardozo@gapef.com.ar",
		Username: "ncardozo",
		Password: "$2a$12$8HKZFQTtifYRXmiguKAO2OPp3IxtsnZPEV7f7MnQdl5uzCJwsttci",
		Role:     "ATLETHE",
	}

	// we define the spected vehabior of the mock
	mockUserRepository.EXPECT().FindByUsername("ncardozo").Return(foundUser, nil, false)

	// setup the application
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	recorder := httptest.NewRecorder() // The recorder records the response of the handler
	context := e.NewContext(request, recorder)

	// Assertions
	if assert.NoError(t, handler.Login(context)) {
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Regexp(t, regexp.MustCompile(`^\n{0,1}.{1,}[\w\-]{1,}\.{1,1}[\w\-]{1,}\.{1,1}[\w\-]{1,}.{1,}\n{0,1}$`),
			recorder.Body)
	}
}

func TestLoginUserNotFound(t *testing.T) {
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	mockUserRepository.EXPECT().FindByUsername("ncardozo").Return(Entity{}, mongo.ErrNoDocuments, true)

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder() // The recorder records the response of the handler
	context := e.NewContext(request, recorder)

	if assert.NoError(t, handler.Login(context)) {
		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), MESSAGE_USER_NOT_FOUND)
	}

}

func TestLoginBadRequest(t *testing.T) {
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	mockUserRepository.EXPECT().FindByUsername("ncardozo").Times(0)

	badRequests := []string{`{"uname": "ncardozo","password":"ncardozo"}`,
		`{"username": "ncardozo","pass":"ncardozo"}`}

	for _, badRequest := range badRequests {
		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(badRequest))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder() // The recorder records the response of the handler
		context := e.NewContext(request, recorder)

		if assert.NoError(t, handler.Login(context)) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			assert.Contains(t, recorder.Body.String(), MESSAGE_BINDING_ERROR)
		}
	}
}

func TestGetAllUsers(t *testing.T) {
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	// what the repository will return
	users := []Entity{
		{
			Username: "ncardozo",
			Password: "ANITALAVALATINA",
			Id:       "1",
			Email:    "ncardozo@gapef.com.ar",
			Role:     "ATLETHE",
		},
		{
			Username: "joan",
			Password: "ANITALAVALATINA",
			Id:       "2",
			Email:    "joan@gapef.com.ar",
			Role:     "COACH",
		},
	}

	usersDTOs := []DTO{
		{Id: "1", Username: "ncardozo", Email: "ncardozo@gapef.com.ar", Role: "ATLETHE", Password: ""},
		{Id: "2", Username: "joan", Email: "joan@gapef.com.ar", Role: "COACH", Password: ""},
	}

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/users", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// We create the query string and set it in the request
	requestQueryString := request.URL.Query()
	requestQueryString.Add("page", "1")
	requestQueryString.Add("size", "4")

	request.URL.RawQuery = requestQueryString.Encode()

	recorder := httptest.NewRecorder() // The recorder records the response of the handler
	context := e.NewContext(request, recorder)

	mockUserRepository.EXPECT().GetUsers(int64(1), int64(4)).Return(users, nil)

	if assert.NoError(t, handler.GetAllUsers(context)) {
		responseJson, _ := json.Marshal(usersDTOs)
		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, string(responseJson), strings.Replace(recorder.Body.String(), "\n", "", 1))
	}
}

func TestGetAllUsersError(t *testing.T) {

	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/users", strings.NewReader(""))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// We create the query string and set it in the request
	requestQueryString := request.URL.Query()
	requestQueryString.Add("page", "1")
	requestQueryString.Add("size", "4")

	request.URL.RawQuery = requestQueryString.Encode()

	recorder := httptest.NewRecorder() // The recorder records the response of the handler
	context := e.NewContext(request, recorder)

	mockUserRepository.EXPECT().GetUsers(int64(1), int64(4)).Return([]Entity{}, errors.New("Error"))

	assert.NoError(t, handler.GetAllUsers(context))
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
