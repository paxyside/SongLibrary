package application

import (
	"context"
	"effictiveMobile/internal/domain/service"
	"effictiveMobile/internal/infrastrtucture/external_api"
	"effictiveMobile/internal/infrastrtucture/http_controller"
	"effictiveMobile/internal/infrastrtucture/persistence"
	"effictiveMobile/pkg/config"
	"effictiveMobile/pkg/database"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}))
	logger.Info("Starting application")

	DBURI := config.Config.DatabaseURI()

	db, err := database.Init(DBURI)
	if err != nil {
		logger.Error("error initializing database", "error", err, "dbURL", DBURI)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error closing database", "error", err)
		}
	}()

	logger.Info("init database")

	// init repository
	songRepo := persistence.NewSongRepository(db, logger)

	// init services
	apiCli := external_api.NewClient()
	songService := service.NewSongService(songRepo, logger, apiCli)

	// init controllers
	songController := http_controller.NewSongController(songService, logger)

	r := mux.NewRouter()

	// add subprefix to routes
	route := r.PathPrefix("/api/v1").Subrouter()

	// init routes for songs
	songsRouter := route.PathPrefix("/songs").Subrouter()
	songsRouter.Use(http_controller.Auth)
	songsRouter.HandleFunc("", songController.GetSongsHandler).Methods("GET")
	songsRouter.HandleFunc("/{id:[0-9]+}", songController.GetSongByIDHandler).Methods("GET")
	songsRouter.HandleFunc("/create", songController.CreateSongHandler).Methods("POST")
	songsRouter.HandleFunc("/update/{id:[0-9]+}", songController.UpdateSongHandler).Methods("PUT")
	songsRouter.HandleFunc("/delete/{id:[0-9]+}", songController.DeleteSongHandler).Methods("DELETE")

	// ping endpoint
	route.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}).Methods("GET")

	// swagger docs
	route.HandleFunc("/docs/spec", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	}).Methods("GET")

	// swagger UI endpoint
	specUrl := config.Config.ServerURI() + "/api/v1/docs/spec"

	route.PathPrefix("/docs/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(specUrl),
	))

	// server config
	address := fmt.Sprintf("%s:%s", config.Config.ServerHost(), config.Config.ServerPort())

	server := http.Server{
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      route,
	}

	logger.Info("starting server", "address", address)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen", "address", address, "error", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case d := <-quit:
		logger.Info("servers stopped", "reason", d.String())
	}

	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		logger.Warn("timeout of 5 seconds.")
	}

	logger.Info("Server exiting")
}
