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

// NewRouter собирает весь HTTP-слой: роутер, middleware, зависимости и маршруты
func NewRouter(cfg config.Config, db *pgxpool.Pool, logger *slog.Logger) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(logger))

	// deps
	userRepo := postgres.NewUserRepo(db)
	authSvc := service.NewAuthService(userRepo, cfg.Auth.JWTSecret)

	roomRepo := postgres.NewRoomRepo(db)
	roomSvc := service.NewRoomService(roomRepo)

	// handlers
	authH := handlers.NewAuthHandler(authSvc)
	healthH := handlers.NewHealthHandler(db)
	meH := handlers.NewMeHandler()
	roomsH := handlers.NewRoomsHandler(roomSvc)

	// public
	r.GET("/health", healthH.Health)
	
	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	// protected (any logged in)
	authMW := middleware.Auth(cfg.Auth.JWTSecret)
	r.GET("/me", authMW, meH.Me)

	// public rooms
	r.GET("/rooms", roomsH.List)
	r.GET("/rooms/:id", roomsH.Get)

	// admin rooms (IMPORTANT: Auth -> RequireRole)
	admin := r.Group("/")
	admin.Use(authMW)
	admin.Use(middleware.RequireRole("admin"))
	admin.POST("/rooms", roomsH.Create)
	admin.PUT("/rooms/:id", roomsH.Update)
	admin.DELETE("/rooms/:id", roomsH.Deactivate)

	return r
}