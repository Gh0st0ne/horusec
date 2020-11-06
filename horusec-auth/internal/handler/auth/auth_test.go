package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	authEntities "github.com/ZupIT/horusec/development-kit/pkg/entities/auth"
	authEnums "github.com/ZupIT/horusec/development-kit/pkg/enums/auth"
	authUseCases "github.com/ZupIT/horusec/development-kit/pkg/usecases/auth"
	"github.com/ZupIT/horusec/horusec-auth/config/app"
	authController "github.com/ZupIT/horusec/horusec-auth/internal/controller/auth"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAuthController(t *testing.T) {
	t.Run("should success create new controller", func(t *testing.T) {
		appConfig := &app.Config{}
		handler := NewAuthHandler(nil, appConfig)
		assert.NotEmpty(t, handler)
	})
}

func TestOptions(t *testing.T) {
	t.Run("should return 204 when options", func(t *testing.T) {
		appConfig := &app.Config{}
		handler := NewAuthHandler(nil, appConfig)
		r, _ := http.NewRequest(http.MethodOptions, "test", nil)
		w := httptest.NewRecorder()

		handler.Options(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestAuthByType(t *testing.T) {
	t.Run("should return 200 when successful login", func(t *testing.T) {
		controllerMock := &authController.MockAuthController{}

		controllerMock.On("AuthByType").Return(map[string]interface{}{"test": "test"}, nil)

		handler := Handler{
			appConfig:      &app.Config{},
			authUseCases:   authUseCases.NewAuthUseCases(),
			authController: controllerMock,
		}

		credentialsBytes, _ := json.Marshal(authEntities.Credentials{Username: "test", Password: "test"})

		r, _ := http.NewRequest(http.MethodPost, "test", bytes.NewReader(credentialsBytes))
		w := httptest.NewRecorder()

		r.Header.Add("X_AUTH_TYPE", "horusec")

		handler.AuthByType(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return 500 when something went wrong", func(t *testing.T) {
		controllerMock := &authController.MockAuthController{}

		controllerMock.On("AuthByType").Return(map[string]interface{}{"test": "test"}, errors.New("test"))

		handler := Handler{
			appConfig:      &app.Config{},
			authUseCases:   authUseCases.NewAuthUseCases(),
			authController: controllerMock,
		}

		credentialsBytes, _ := json.Marshal(authEntities.Credentials{Username: "test", Password: "test"})

		r, _ := http.NewRequest(http.MethodPost, "test", bytes.NewReader(credentialsBytes))
		w := httptest.NewRecorder()

		r.Header.Add("X_AUTH_TYPE", "horusec")

		handler.AuthByType(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return 400 when invalid credentials", func(t *testing.T) {
		controllerMock := &authController.MockAuthController{}

		handler := Handler{
			appConfig:      &app.Config{},
			authUseCases:   authUseCases.NewAuthUseCases(),
			authController: controllerMock,
		}

		credentialsBytes, _ := json.Marshal(authEntities.Credentials{})

		r, _ := http.NewRequest(http.MethodPost, "test", bytes.NewReader(credentialsBytes))
		w := httptest.NewRecorder()

		r.Header.Add("X_AUTH_TYPE", "horusec")

		handler.AuthByType(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_AuthTypes(t *testing.T) {
	t.Run("should return 200 when get auth types", func(t *testing.T) {
		handler := NewAuthHandler(nil, &app.Config{
			AuthType: authEnums.Horusec.ToString(),
		})

		r, _ := http.NewRequest(http.MethodGet, "test", nil)
		w := httptest.NewRecorder()

		handler.Config(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("should return 200 when get auth types mocked", func(t *testing.T) {
		controllerMock := &authController.MockAuthController{}
		controllerMock.On("GetAuthType").Return(authEnums.Horusec, nil)
		handler := Handler{
			appConfig: &app.Config{
				AuthType: "test",
			},
			authUseCases:   authUseCases.NewAuthUseCases(),
			authController: controllerMock,
		}

		r, _ := http.NewRequest(http.MethodGet, "test", nil)
		w := httptest.NewRecorder()

		handler.Config(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}