package route

import (
	"github.com/Asylann/Auth/internal/config"
	"github.com/Asylann/Auth/internal/handler"
	"github.com/Asylann/Auth/internal/middleware"
	"github.com/Asylann/Auth/internal/repository"
	"github.com/Asylann/Auth/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type Route struct {
	Cfg    config.Config
	Logger *logrus.Logger
	Srv    *http.Server
}

func NewRoute(cfg config.Config, logger *logrus.Logger) Route {
	return Route{Cfg: cfg, Logger: logger}
}

func (Route *Route) Run() {
	repo, err := repository.NewRepository(Route.Cfg)
	if err != nil {
		Route.Logger.Fatalf("Cant connect database: %s", err.Error())
		return
	}
	Route.Logger.Infof("Databse connected")

	route := gin.Default()
	/*route.Use(middleware.Logger(), gin.Recovery())*/

	serv := service.NewService(repo, Route.Logger)
	hd := handler.New(Route.Logger, Route.Cfg, serv)

	route.POST("/login", hd.Login)
	route.POST("/signup", hd.RegisterUser)

	auth := route.Group("/auth", middleware.Auth(Route.Cfg.JWTSecret, Route.Logger))
	adminRoute := auth.Group("/", middleware.RequireRole("admin", Route.Logger))
	adminRoute.GET("/users/ByID/:id", hd.GetUserById)
	adminRoute.GET("/users", hd.GetListOfUsers)
	adminRoute.GET("/users/:email", hd.GetUserByEmail)

	Route.Srv = &http.Server{
		Addr:         ":" + Route.Cfg.Port,
		Handler:      route,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	Route.Logger.Infof("Server is running on %v port", Route.Cfg.Port)
	if err := Route.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		Route.Logger.Errorf("Couldnt run server: %s", err.Error())
		return
	}
}

func (Route *Route) GracefullyShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Route.Srv.Shutdown(ctx); err != nil {
		Route.Logger.Fatalf("Server is forced to shut down: %s", err.Error())
	}

	Route.Logger.Info("Server is exited")
}
