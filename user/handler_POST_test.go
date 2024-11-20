package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
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

	request := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(string(requestDTO)))
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	context := e.NewContext(request, recorder)

	assert.NoError(t, handler.Create(context))
	assert.Equal(t, http.StatusCreated, recorder.Code)
}
