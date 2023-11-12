package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/MatheusAbdias/go_simple_bank/api"
	"github.com/MatheusAbdias/go_simple_bank/cmd"
	"github.com/MatheusAbdias/go_simple_bank/config"
	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
	_ "github.com/MatheusAbdias/go_simple_bank/docs/statik"
	"github.com/MatheusAbdias/go_simple_bank/gapi"
	"github.com/MatheusAbdias/go_simple_bank/pb"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	conn, err := sql.Open(config.Driver, config.Source)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	cmd.Migrate()

	store := db.NewSQLStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGatewayServer(config config.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs:", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start gateway server", err)
	}
}

func runGrpcServer(config config.Config, store db.Store) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config config.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

}
