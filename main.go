package main

import (
	"context"
	"database/sql"
	"github.com/Bakhram74/small_bank/api"
	db "github.com/Bakhram74/small_bank/db/sqlc"
	_ "github.com/Bakhram74/small_bank/doc/statik"
	"github.com/Bakhram74/small_bank/gapi"
	"github.com/Bakhram74/small_bank/pb"
	"github.com/Bakhram74/small_bank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	runDBMigration(config.MigrationURL, config.DBSource)
	store := db.NewStore(conn)

	go runGatewayServer(store, config)

	runGrpcServer(store, config)

}
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {

		log.Fatal("failed to run migrate up ", err)
	}

	log.Println("db migrated successfully")
}
func runGrpcServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSmallBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener %s", err)
	}

	log.Printf("start GRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("cannot start GRPC server %s", err)
	}

}

func runGatewayServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	ctx, cansel := context.WithCancel(context.Background())
	defer cansel()
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		}})
	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterSmallBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)

	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("cannot create statik file %s", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener %s", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("cannot start HTTP gateway server %s", err)
	}

}

func runGinServer(store db.Store, config util.Config) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
