package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/shubhangcs/agromart-server/internal/blob"
	"github.com/shubhangcs/agromart-server/internal/handlers"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/migrations"
)

type Application struct {
	Logger          *log.Logger
	DB              *sql.DB
	Blob            *blob.AWSS3
	UserHandler     *handlers.UserHandler
	BlobHandler     *handlers.BlobHandler
	TokenHandler    *handlers.TokenHandler
	BusinessHandler *handlers.BusinessHandler
	CategoryHandler *handlers.CategoryHandler
	FollowHandler   *handlers.FollowerHandler
	RFQHandler      *handlers.RFQHandler
	ProductHandler  *handlers.ProductHandler
}

func NewApplication() (*Application, error) {
	// Creating a new logger
	logger := log.New(os.Stdout, "AGROMART ", log.Ldate|log.Ltime)

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

	// Handlers
	userHandler := handlers.NewUserHandler(userStore, logger)
	tokenHandler := handlers.NewTokenHandler(userStore, businessStore, logger)
	blobHandler := handlers.NewBlobHandler(logger, as3, blobStore)
	businessHandler := handlers.NewBusinessHandler(businessStore, logger)
	categoryHandler := handlers.NewCategoryHandler(categoryStore, logger)
	followerHandler := handlers.NewFollowerHandler(followerStore, logger)
	rfqHandler := handlers.NewRFQHandler(rfqStore, logger)
	productHandler := handlers.NewProductHandler(productStore, logger)

	// Creating a object of application struct
	app := &Application{
		Logger:          logger,
		DB:              pgdb,
		Blob:            as3,
		UserHandler:     userHandler,
		TokenHandler:    tokenHandler,
		BlobHandler:     blobHandler,
		BusinessHandler: businessHandler,
		CategoryHandler: categoryHandler,
		FollowHandler:   followerHandler,
		RFQHandler:      rfqHandler,
		ProductHandler:  productHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := a.DB.Ping(); err != nil {
		a.Logger.Printf("ERROR: healthCheck: %s\n", err.Error())
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"status": fmt.Sprintf("error: %s", err.Error())})
		return
	}
	a.Logger.Println("INFO: healthCheck: system is healthy")
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "system is healthy"})
}
