package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	"github.com/Oyinoye/bank_mini/gapi"
	"github.com/Oyinoye/bank_mini/pb"
	"github.com/Oyinoye/bank_mini/util"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Oyinoye/bank_mini/api"
	db "github.com/Oyinoye/bank_mini/db/sqlc"
	_ "github.com/lib/pq"
)


func main() {

    config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	// if config.Environment == "development" {
	// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// }


	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

    runGinServer(config, store)


    // server, err := api.NewServer(config, store)
    // if err != nil {
    //     log.Fatal("cannot create server:", err)
    // }

	// err = server.Start(config.ServerAddress)
    // if err != nil {
    //     log.Fatal("cannot start server:", err)
    // }
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
// 	taskDistributor worker.TaskDistributor,
) {
	server, err := gapi.NewServer(config, store,
        // taskDistributor
    )
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create server")
		log.Fatal("cannot create server")
	}

	// gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	// grpcServer := grpc.NewServer(gprcLogger)
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create listener")
		log.Fatal("cannot create listener")
	}

    log.Printf("start gRPC server at %s", listener.Addr().String())

    err = grpcServer.Serve(listener)
    if err != nil {
        log.Fatal("cannot start gRPC server")
    }

	// waitGroup.Go(func() error {
	// 	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	// 	err = grpcServer.Serve(listener)
	// 	if err != nil {
	// 		if errors.Is(err, grpc.ErrServerStopped) {
	// 			return nil
	// 		}
	// 		log.Error().Err(err).Msg("gRPC server failed to serve")
	// 		return err
	// 	}

	// 	return nil
	// })

	// waitGroup.Go(func() error {
	// 	<-ctx.Done()
	// 	log.Info().Msg("graceful shutdown gRPC server")

	// 	grpcServer.GracefulStop()
	// 	log.Info().Msg("gRPC server is stopped")

	// 	return nil
	// })
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot create server")
		log.Fatal("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		// log.Fatal().Err(err).Msg("cannot start server")
		log.Fatal("cannot start server")
	}
}
