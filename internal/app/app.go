package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/shubhangcs/agromart-server/internal/blob"
	"github.com/shubhangcs/agromart-server/internal/handlers"
	"github.com/shubhangcs/agromart-server/internal/hub"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/migrations"
)

type Application struct {
	Logger          *slog.Logger
	DB              *sql.DB
	Blob            *blob.AWSS3
	Hub             *hub.Hub
	UserHandler     *handlers.UserHandler
	BlobHandler     *handlers.BlobHandler
	TokenHandler    *handlers.TokenHandler
	BusinessHandler *handlers.BusinessHandler
	CategoryHandler *handlers.CategoryHandler
	FollowHandler   *handlers.FollowerHandler
	RFQHandler      *handlers.RFQHandler
	ProductHandler  *handlers.ProductHandler
	RatingHandler   *handlers.RatingHandler
	ReviewHandler   *handlers.ReviewHandler
	ChatHandler     *handlers.ChatHandler
}

func NewApplication() (*Application, error) {
	// Creating a new structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Connecting to database
	pgdb, err := store.Open()
	if err != nil {
		return nil, err
	}

	// Migrate
	err = store.MigrateFS(pgdb, migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	// Connecting to Blob
	as3, err := blob.Connect()
	if err != nil {
		return nil, err
	}

	// Store
	userStore := store.NewPostgresUserStore(pgdb)
	businessStore := store.NewPostgresBusinessStore(pgdb)
	categoryStore := store.NewPostgresCategoryStore(pgdb)
	followerStore := store.NewPostgresFollowerStore(pgdb)
	rfqStore := store.NewPostgresRFQStore(pgdb)
	productStore := store.NewPostgresProductStore(pgdb)
	blobStore := store.NewPostgresBlobStore(pgdb)
	ratingStore := store.NewPostgresRatingStore(pgdb)
	reviewStore := store.NewPostgresReviewStore(pgdb)
	chatStore := store.NewPostgresChatStore(pgdb)

	// Handlers
	userHandler := handlers.NewUserHandler(userStore, logger)
	tokenHandler := handlers.NewTokenHandler(userStore, businessStore, logger)
	blobHandler := handlers.NewBlobHandler(logger, as3, blobStore)
	businessHandler := handlers.NewBusinessHandler(businessStore, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryStore, logger)
	followerHandler := handlers.NewFollowerHandler(followerStore, logger)
	rfqHandler := handlers.NewRFQHandler(rfqStore, logger)
	productHandler := handlers.NewProductHandler(productStore, logger)
	ratingHandler := handlers.NewRatingHandler(ratingStore, logger)
	reviewHandler := handlers.NewReviewHandler(reviewStore, logger)
	wsHub := hub.NewHub()
	chatHandler := handlers.NewChatHandler(chatStore, wsHub, logger)

	// Creating a object of application struct
	app := &Application{
		Logger:          logger,
		DB:              pgdb,
		Blob:            as3,
		Hub:             wsHub,
		UserHandler:     userHandler,
		TokenHandler:    tokenHandler,
		BlobHandler:     blobHandler,
		BusinessHandler: businessHandler,
		CategoryHandler: categoryHandler,
		FollowHandler:   followerHandler,
		RFQHandler:      rfqHandler,
		ProductHandler:  productHandler,
		RatingHandler:   ratingHandler,

		ReviewHandler:   reviewHandler,
		ChatHandler:     chatHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := a.DB.Ping(); err != nil {
		a.Logger.Error("healthCheck", "error", err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"status": fmt.Sprintf("error: %s", err.Error())})
		return
	}
	a.Logger.Info("healthCheck: system is healthy")
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "system is healthy"})
}
