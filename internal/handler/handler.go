package handler

import (
	"errors"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/Asylann/Auth/internal/config"
	"github.com/Asylann/Auth/internal/model"
	"github.com/Asylann/Auth/internal/response"
	"github.com/Asylann/Auth/internal/service"
	"github.com/Asylann/Auth/lib/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type Handler struct {
	Logger        *logrus.Logger
	Cfg           config.Config
	Service       service.Service
	emailVerifier *emailverifier.Verifier
}

func New(logger *logrus.Logger, cfg config.Config, service service.Service) Handler {
	return Handler{Logger: logger, Cfg: cfg, Service: service}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (hd *Handler) Login(c *gin.Context) {
	var lr LoginRequest
	if err := c.BindJSON(&lr); err != nil {
		hd.Logger.Errorf("Coundt parse given json: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, response.Response{Err: "Cant convert properly"})
		return
	}

	userInDB, err := hd.Service.GetUserByEmail(c.Request.Context(), lr.Email)
	if err != nil {
		hd.Logger.Errorf("Coundt find such user by email: %s", err.Error())
		c.IndentedJSON(http.StatusNotFound, response.Response{Err: "Not found user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInDB.Password), []byte(lr.Password)); err != nil {
		hd.Logger.Error("Invalid password")
		c.IndentedJSON(http.StatusNotAcceptable, response.Response{Err: "Invalid password"})
		return
	}

	tokenString, err := utils.GenerateJWT(lr.Email, userInDB.Role, hd.Cfg.JWTSecret, 24)
	if err != nil {
		hd.Logger.Errorf("Couldnt convert to JWT: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, response.Response{Err: "smt went wrong"})
		return
	}

	c.SetCookie("auth_token", tokenString, 3600, "", "localhost", true, true)
	c.IndentedJSON(http.StatusOK, response.Response{Data: "Logged in"})
}

func (hd *Handler) RegisterUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		hd.Logger.Errorf("Couldnt convert json %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, response.Response{Err: "Couldn't convert json"})
		return
	}

	err = hd.ValidateEmail(user.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, response.Response{Err: err.Error()})
		return
	}

	id, err := hd.Service.RegisterUser(c.Request.Context(), user)
	if err != nil {
		hd.Logger.Errorf("Couldnt create a user: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, response.Response{Err: "smt went wrong"})
		return
	}

	hd.Logger.Infof("User was created with id=%s and email=%s", id, user.Email)
	c.IndentedJSON(http.StatusOK, response.Response{Data: id})
}

func (hd *Handler) ValidateEmail(email string) error {
	hd.emailVerifier = emailverifier.NewVerifier()
	/*emailVerifier = emailVerifier.EnableSMTPCheck()*/
	hd.emailVerifier = hd.emailVerifier.EnableDomainSuggest()
	hd.emailVerifier = hd.emailVerifier.AddDisposableDomains([]string{"fucku.com"})

	res, err := hd.emailVerifier.Verify(email)
	if err != nil {
		hd.Logger.Errorf("Invalid email: %s", err.Error())
		return errors.New("Invalid email text, try different one")
	}

	if !res.Syntax.Valid {
		hd.Logger.Errorf("Invalid Syntax to email")
		return errors.New("Invalid syntax to email text, try different one")
	}

	if res.Disposable {
		hd.Logger.Errorf("Disposable email")
		return errors.New("Sorry, we dont accept disposable emails")
	}

	if res.Suggestion != "" {
		hd.Logger.Errorf("Invalid email suggestion %s:", res.Suggestion)
		return errors.New("Invalid email text, maybe you mean" + res.Suggestion + " ?")
	}

	if res.Reachable == "no" {
		hd.Logger.Errorf("Invalid email unreachable")
		return errors.New("Invalid email text, unreachable")
	}

	return nil
}

func (hd *Handler) GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		hd.Logger.Errorf("Couldnt convert string to int %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, response.Response{Err: "Couldn't find id in params"})
		return
	}

	user, err := hd.Service.GetUserById(c.Request.Context(), id)
	if err != nil {
		hd.Logger.Errorf("Couldnt get user by id: %s", err.Error())
		c.IndentedJSON(http.StatusNotFound, response.Response{Err: "Couldn't find User"})
		return
	}

	hd.Logger.Infof("User by email=%s was received", user.Email)
	c.IndentedJSON(http.StatusOK, response.Response{Data: user})
}

func (hd *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := hd.Service.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		hd.Logger.Errorf("Couldnt get user by email: %s", err.Error())
		c.IndentedJSON(http.StatusNotFound, response.Response{Err: "Couldn't find User"})
		return
	}

	hd.Logger.Infof("User by email=%s was received", user.Email)
	c.IndentedJSON(http.StatusOK, response.Response{Data: user})
}

func (hd *Handler) GetListOfUsers(c *gin.Context) {
	users, err := hd.Service.GetListOfUsers(c.Request.Context())
	if err != nil {
		hd.Logger.Errorf("Couldnt get all user: %s", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, response.Response{Err: "Smt went wrong"})
		return
	}

	hd.Logger.Infof("Users  were received")
	c.IndentedJSON(http.StatusOK, response.Response{Data: users})
}
