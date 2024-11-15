package app

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "clicker/internal/config"
    "clicker/pkg/clicker"
    "github.com/gorilla/mux"
    "github.com/redis/go-redis/v9"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type App struct {
	cfg    *config.Config
	router *mux.Router
	grpc   *grpc.Server
	redis  *redis.Client
}

func New(cfg *config.Config) *App {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddress(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &App{
		cfg:    cfg,
		router: mux.NewRouter(),
		grpc:   grpc.NewServer(),
		redis:  rdb,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем REST маршруты
	a.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "REST API is working")
	})

	// HTTP сервер
	httpServer := &http.Server{
		Addr:    a.cfg.GetRestAddress(),
		Handler: a.router,
	}

	// Запускаем REST сервер в горутине
	go func() {
		log.Printf("Starting REST server on %s", a.cfg.GetRestAddress())
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("REST server error: %v", err)
		}
	}()

	// Запускаем gRPC сервер в горутине
	go func() {
		lis, err := net.Listen("tcp", a.cfg.GetGrpcAddress())
		if err != nil {
			log.Printf("Failed to listen gRPC: %v", err)
			return
		}
		log.Printf("Starting gRPC server on %s", a.cfg.GetGrpcAddress())
		if err := a.grpc.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown
	log.Println("Shutting down servers...")

	// Останавливаем HTTP сервер
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Останавливаем gRPC сервер
	a.grpc.GracefulStop()

	return nil
}
