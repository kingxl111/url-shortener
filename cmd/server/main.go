package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	serv "github.com/kingxl111/url-shortener/internal/gates/grpc"
	"github.com/kingxl111/url-shortener/internal/repository/factory"
	urlSrv "github.com/kingxl111/url-shortener/internal/url/service"
	"golang.org/x/sync/errgroup"

	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kingxl111/url-shortener/internal/config"
	"google.golang.org/grpc/reflection"

	en "github.com/kingxl111/url-shortener/internal/environment"
	desc "github.com/kingxl111/url-shortener/pkg/shortener"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	defaultLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(defaultLogger)

	if err := runMain(ctx); err != nil {
		defaultLogger.Error("run main", slog.Any("err", err))
		return
	}
}

func runMain(ctx context.Context) error {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		return fmt.Errorf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		return fmt.Errorf("failed to get pg config: %v", err)
	}

	loggerConfig, err := config.NewLoggerConfig()
	if err != nil {
		return fmt.Errorf("failed to get logger config: %v", err)
	}

	handleOpts := &slog.HandlerOptions{
		Level: loggerConfig.Level(),
	}
	var h slog.Handler = slog.NewTextHandler(os.Stdout, handleOpts)
	logger := slog.New(h)

	// waiting for db init
	time.Sleep(time.Second * 3)
	repo, err := factory.NewURLRepository(
		pgConfig.Username,
		pgConfig.Password,
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.DBName,
		pgConfig.SSLMode)
	if err != nil {
		return fmt.Errorf("failed to create repository: %v", err)
	}
	service := urlSrv.New(repo)

	var opts en.ServerOptions
	opts.WithLogger(logger)

	grpcServ, err := opts.NewServer()
	if err != nil {
		return fmt.Errorf("failed to create grpc server: %v", err)
	}

	reflection.Register(grpcServ)
	desc.RegisterURLShortenerServer(grpcServ, &serv.Server{Services: service})

	eg, gctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		logger.Info("grpc starting", slog.String("addr", grpcConfig.Address()))
		return en.ListenAndServeContext(gctx, grpcConfig.Address(), grpcServ)
	})

	return eg.Wait()
}
