package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/ncardozo92/gapef_swimming_metrics/constants"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserSuccess(t *testing.T) {
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	e := echo.New()

	requestDTO, _ := json.Marshal(DTO{Email: "ncardozo@gapef.com.ar", Username: "ncardozo", Password: "anitaLAVAlaTina", Role: "ATLETHE"})

	mockUserRepository.EXPECT().Create(gomock.Any()).Return(nil)
	mockUserRepository.EXPECT().Exists(gomock.Any()).Return(false, nil)

	request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(requestDTO)))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	assert.NoError(t, handler.Create(context))
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestCreateUserInvalidDTOFails(t *testing.T) {
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	// test table
	testCases := []DTO{
		{Email: "NCARDOZO", Username: "ncardozo", Password: "1234asdf", Role: constants.ROLE_ATLETHE},
		{Email: "ncardozo@gapef.com.ar", Username: "", Password: "1234asdf", Role: constants.ROLE_ATLETHE},
		{Email: "ncardozo@gapef.com.ar", Username: "ncardozo", Password: "", Role: constants.ROLE_ATLETHE},
		{Email: "ncardozo@gapef.com.ar", Username: "ncardozo", Password: "1234asdf", Role: "undefined"},
	}

	e := echo.New()

	for _, testCase := range testCases {

		jsonMarshall, _ := json.Marshal(testCase)

		request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(jsonMarshall)))
		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)

		assert.NoError(t, handler.Create(context))
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	}

}

func TestCreateDuplicatedUserFails(t *testing.T) {
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	e := echo.New()

	testCasesDtos := []DTO{
		{Email: "ncardozo@gapef.com.ar", Username: "ncardozo92", Password: "1234asdf", Role: constants.ROLE_ATLETHE},
		{Email: "nc92030@gapef.com.ar", Username: "ncardozo", Password: "1234asdf", Role: constants.ROLE_ATLETHE},
	}

	for _, dto := range testCasesDtos {

		mockUserRepository.EXPECT().Exists(gomock.Any()).Return(true, nil)
		mockUserRepository.EXPECT().Create(gomock.Any()).Times(0)

		json, _ := json.Marshal(dto)
		request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(json)))
		request.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)

		assert.NoError(t, handler.Create(context))
		assert.Equal(t, http.StatusConflict, recorder.Code)
	}
}

func TestCreateUserFindExistingFails(t *testing.T) {
	controller := gomock.NewController(t)
	mockUserRepository := NewMockRepository(controller)
	handler := NewUserHandler(mockUserRepository)
	defer controller.Finish()

	e := echo.New()

	dto := DTO{Email: "nc92030@gapef.com.ar", Username: "ncardozo", Password: "1234asdf", Role: constants.ROLE_ATLETHE}

	mockUserRepository.EXPECT().Exists(gomock.Any()).Return(false, errors.New("Opps..."))
	mockUserRepository.EXPECT().Create(gomock.Any()).Times(0)

	json, _ := json.Marshal(dto)
	request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(json)))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	assert.NoError(t, handler.Create(context))
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
