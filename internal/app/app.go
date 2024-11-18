package app

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "clicker/internal/application/usecase"
    "clicker/internal/config"
    "clicker/internal/domain/repository"
    "clicker/internal/interfaces/grpc/handler"
    "clicker/pkg/counter"
    "clicker/pkg/stats"

    "github.com/gorilla/mux"
    "google.golang.org/grpc"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc/credentials/insecure"
    _ "github.com/lib/pq"
)

type App struct {
    cfg    *config.Config
    router *mux.Router
    grpc   *grpc.Server
    db     *pgxpool.Pool
}

func New(cfg *config.Config) *App {
    dbConfig, err := pgxpool.ParseConfig(cfg.GetPostgresDSN())
    if err != nil {
        log.Fatalf("Unable to parse PostgreSQL DSN: %v", err)
    }

    db, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    grpcServer := grpc.NewServer()

    clickRepo := repository.NewPostgresClickRepository(db)
    statsRepo := repository.NewPostgresStatsRepository(db)

    clickUseCase := usecase.NewClickUseCase(clickRepo)
    statsUseCase := usecase.NewStatsUseCase(statsRepo)

    clickHandler := handler.NewClickHandler(clickUseCase)
    statsHandler := handler.NewStatsHandler(statsUseCase)

    grpcHandler := handler.NewHandler(clickHandler, statsHandler)
    grpcHandler.Register(grpcServer)

    router := mux.NewRouter()

    gwmux := runtime.NewServeMux()

    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

    if err := counter.RegisterCounterServiceHandlerFromEndpoint(context.Background(), 
        gwmux, cfg.GetGrpcAddress(), opts); err != nil {
        log.Fatalf("Не удалось зарегистрировать gateway для CounterService: %v", err)
    }

    if err := stats.RegisterStatsServiceHandlerFromEndpoint(context.Background(), 
        gwmux, cfg.GetGrpcAddress(), opts); err != nil {
        log.Fatalf("Не удалось зарегистрировать gateway для StatsService: %v", err)
    }

    router.PathPrefix("/").Handler(gwmux)

    return &App{
        cfg:    cfg,
        router: router,
        grpc:   grpcServer,
        db:     db,
    }
}

func (a *App) Run() error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    a.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "REST API работает")
    })

    httpServer := &http.Server{
        Addr:    a.cfg.GetRestAddress(),
        Handler: a.router,
    }

    go func() {
        log.Printf("Запуск REST сервера на %s", a.cfg.GetRestAddress())
        if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
            log.Printf("Ошибка REST сервера: %v", err)
        }
    }()

    go func() {
        lis, err := net.Listen("tcp", a.cfg.GetGrpcAddress())
        if err != nil {
            log.Printf("Не удалось слушать gRPC: %v", err)
            return
        }
        log.Printf("Запуск gRPC сервера на %s", a.cfg.GetGrpcAddress())
        if err := a.grpc.Serve(lis); err != nil {
            log.Printf("Ошибка gRPC сервера: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Остановка серверов...")

    shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
    defer shutdownCancel()

    if err := httpServer.Shutdown(shutdownCtx); err != nil {
        log.Printf("Ошибка при остановке HTTP сервера: %v", err)
    }

    a.grpc.GracefulStop()
    a.db.Close()

    return nil
}
