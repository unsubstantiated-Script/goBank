package gapi

import (
	"fmt"
	db "goBank/db/sqlc"
	"goBank/pb"
	"goBank/token"
	"goBank/util"
	"goBank/worker"
)

// This server serves up all gRPC requests for our banking service.
type Server struct {
	// Enables forward compatibility and the ability to gradually add real implementations later.
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server and setup routing.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	if tokenMaker == nil {
		return nil, fmt.Errorf("cannot create token maker")
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
