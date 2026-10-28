package server

import (
	"log/slog"
	"net/http"

	"github.com/callmhejerry/booking-api/internal/bookings"
	"github.com/callmhejerry/booking-api/internal/config"
	"github.com/callmhejerry/booking-api/internal/db"
	"github.com/callmhejerry/booking-api/internal/middleware"
	"github.com/callmhejerry/booking-api/internal/slots"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config   *config.Config
	database *db.Queries
	dbConn   *pgxpool.Pool

	slotService    *slots.SlotService
	bookingService *bookings.BookingService
}

func NewServer(
	config *config.Config,
	database *db.Queries,
	logger *slog.Logger,
	dbConn *pgxpool.Pool,
) *Server {
	slotRepository := slots.NewRepositoryImpl(database)
	slotService := slots.NewSlotService(slotRepository, logger)

	bookingRepository := bookings.NewBookingRepositoryImpl(database, dbConn)
	bookingCache := bookings.NewBookingCache()
	bookingService := bookings.NewBookingService(bookingRepository, slotRepository, bookingCache, logger)

	return &Server{
		config:         config,
		database:       database,
		slotService:    slotService,
		bookingService: bookingService,
		dbConn:         dbConn,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	gin.SetMode(s.config.AppCongfig.Mode)
	router := gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(s.config.AppCongfig.CORSAllowedOrigins))
	router.Use(middleware.AuthMiddleware(s.config.AppCongfig.ApiKey))

	slots.NewSlotController(s.slotService, router)
	bookings.NewBookingController(s.bookingService, router)
	return router
}
