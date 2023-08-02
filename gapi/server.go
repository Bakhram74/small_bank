package gapi

import (
	"fmt"
	db "github.com/Bakhram74/small_bank/db/sqlc"
	"github.com/Bakhram74/small_bank/pb"
	"github.com/Bakhram74/small_bank/token"
	"github.com/Bakhram74/small_bank/util"
	"github.com/Bakhram74/small_bank/worker"
)

type Server struct {
	pb.UnimplementedSmallBankServer
	config          util.Config
	tokenMaker      token.Maker
	store           db.Store
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, config util.Config, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
