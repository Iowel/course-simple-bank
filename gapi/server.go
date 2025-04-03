package gapi

import (
	"fmt"

	db "github.com/Iowel/course-simple-bank/db/sqlc"
	"github.com/Iowel/course-simple-bank/pb"
	"github.com/Iowel/course-simple-bank/token"
	"github.com/Iowel/course-simple-bank/util"
	"github.com/Iowel/course-simple-bank/worker"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	// Изменить NewPasetoMaker на NewJWTMaker
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
