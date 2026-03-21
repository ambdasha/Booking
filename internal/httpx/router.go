package httpx

import (
	"log/slog"

	"booking/internal/config"
	"booking/internal/httpx/handlers"
	"booking/internal/httpx/middleware"
	"booking/internal/repository/postgres"
	"booking/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)


//собирает весь HTTP-слой: создает роутер,подключает middleware, создает зависимости и регистрирует маршруты

func NewRouter(cfg config.Config, db *pgxpool.Pool, logger *slog.Logger) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(logger))

	userRepo := postgres.NewUserRepo(db)
	authSvc := service.NewAuthService(userRepo, cfg.Auth.JWTSecret)
	authH := handlers.NewAuthHandler(authSvc)
	healthH := handlers.NewHealthHandler(db)
	meH := handlers.NewMeHandler()

	r.GET("/health", healthH.Health)

	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)
	
	authMW := middleware.Auth(cfg.Auth.JWTSecret)	
	r.GET("/me", authMW, meH.Me)

	return r
}