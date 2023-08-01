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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config:")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db:")
	}
	runDBMigration(config.MigrationURL, config.DBSource)
	store := db.NewStore(conn)

	go runGatewayServer(store, config)

	runGrpcServer(store, config)

}
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance ")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {

		log.Fatal().Err(err).Msg("failed to run migrate up ")
	}

	log.Print("db migrated successfully")
}
func runGrpcServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSmallBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener %s")
	}

	log.Info().Msgf("start GRPC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start GRPC server %s")
	}

}

func runGatewayServer(store db.Store, config util.Config) {
	server, err := gapi.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
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
		log.Fatal().Err(err).Msg("cannot register handler server:")

	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik file %s")
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener %s")
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP gateway server %s")
	}

}

func runGinServer(store db.Store, config util.Config) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server:")
	}
}
