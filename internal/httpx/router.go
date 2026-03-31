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

//сбор всего HTTP слоя
func NewRouter(cfg config.Config, db *pgxpool.Pool, logger *slog.Logger) *gin.Engine {
	r := gin.New()

	// базовые middleware
	r.Use(gin.Recovery()) // чтобы паника не уронила сервер (вернёт 500)
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(logger))


	// SQL слой
	userRepo := postgres.NewUserRepo(db)
	roomRepo := postgres.NewRoomRepo(db)
	reservationRepo := postgres.NewReservationRepo(db)
	blockRepo := postgres.NewBlockRepo(db)

	// services
	authSvc := service.NewAuthService(userRepo, cfg.Auth.JWTSecret)
	roomSvc := service.NewRoomService(roomRepo)
	reservationSvc := service.NewReservationService(reservationRepo, roomRepo)
	blockSvc := service.NewBlockService(blockRepo, roomRepo)
	availabilitySvc := service.NewAvailabilityService(reservationRepo, blockRepo)

	// handlers 
	authH := handlers.NewAuthHandler(authSvc)
	healthH := handlers.NewHealthHandler(db)
	meH := handlers.NewMeHandler()
	roomsH := handlers.NewRoomsHandler(roomSvc)
	reservationsH := handlers.NewReservationsHandler(reservationSvc)
	blocksH := handlers.NewBlocksHandler(blockSvc)
	availabilityH := handlers.NewAvailabilityHandler(availabilitySvc)

	// PUBLIC 
	r.GET("/health", healthH.Health)

	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	// публичные комнаты
	r.GET("/rooms", roomsH.List)
	r.GET("/rooms/:id", roomsH.Get)

	// доступность 
	r.GET("/rooms/:id/availability", availabilityH.Get)

	// PROTECTED (любой залогиненный)
	authMW := middleware.Auth(cfg.Auth.JWTSecret)

	api := r.Group("/")
	api.Use(authMW)

	api.GET("/me", meH.Me)

	api.POST("/reservations", reservationsH.Create)
	api.GET("/reservations/my", reservationsH.MyList)
	api.GET("/reservations/:id", reservationsH.Get)
	api.POST("/reservations/:id/cancel", reservationsH.Cancel)

	//только admin

	admin := r.Group("/admin")
	admin.Use(authMW)                      
	admin.Use(middleware.RequireRole("admin"))

	// admin rooms
	admin.POST("/rooms", roomsH.Create)
	admin.PUT("/rooms/:id", roomsH.Update)
	admin.DELETE("/rooms/:id", roomsH.Deactivate)

	// admin blocks
	admin.POST("/rooms/:id/blocks", blocksH.Create)
	admin.DELETE("/blocks/:block_id", blocksH.Delete)

	return r
}