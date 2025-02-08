package main

import (
	"flag"
	"log"
	"net"

	"github.com/kingxl111/url-shortener/internal/config"
	serv "github.com/kingxl111/url-shortener/internal/handlers"
	pg "github.com/kingxl111/url-shortener/internal/repository/url/postgres"
	urlSrv "github.com/kingxl111/url-shortener/internal/service/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	desc "github.com/kingxl111/url-shortener/pkg/shortener"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	db, err := pg.NewDB(
		pgConfig.Username,
		pgConfig.Password,
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.DBName,
		pgConfig.SSLMode)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}
	defer db.Close()

	repo := pg.NewRepository(db)
	service := urlSrv.New(repo)

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterURLShortenerServer(s, &serv.Server{Services: service})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
