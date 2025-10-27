package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Asylann/Auth/internal/config"
	"github.com/Asylann/Auth/internal/handler"
	"github.com/Asylann/Auth/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MOcking the service to make it unit teest and do not get in touch with DB
type MockSvc struct {
	mock.Mock
}

func (m *MockSvc) RegisterUser(ctx context.Context, user model.User) (string, error) {
	a := m.Called(ctx, user)
	return a.String(0), a.Error(1)
}

func (m *MockSvc) GetUserByEmail(ctx context.Context, e string) (model.User, error) {
	a := m.Called(ctx, e)
	return a.Get(0).(model.User), a.Error(1)
}

func (m *MockSvc) GetUserById(ctx context.Context, id int) (model.User, error) {
	a := m.Called(ctx, id)
	return a.Get(0).(model.User), a.Error(1)
}
func (m *MockSvc) GetListOfUsers(ctx context.Context) ([]model.User, error) {
	a := m.Called(ctx)
	return a.Get(0).([]model.User), a.Error(1)
}

// main test flow
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockServ := new(MockSvc)
	hashed, _ := bcrypt.GenerateFromPassword([]byte("1234"), 12)
	mockUser := model.User{
		ID:       1,
		Email:    "firstuser@gmail.com",
		Password: string(hashed),
	}

	mockServ.On("GetUserByEmail", mock.Anything, "firstuser@gmail.com").Return(mockUser, nil)

	hd := handler.New(logrus.New(), config.Config{}, mockServ)
	router := gin.Default()
	router.POST("/login", hd.Login)

	reqBody := map[string]string{
		"email":    "firstuser@gmail.com",
		"password": "1234",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// main method to make gin route serve as ordinary http router
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
