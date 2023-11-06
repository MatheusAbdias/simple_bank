package gapi

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/MatheusAbdias/go_simple_bank/config"
	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
	"github.com/MatheusAbdias/go_simple_bank/pb"
	"github.com/MatheusAbdias/go_simple_bank/token"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     config.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store, config config.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker, %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	return server, nil
}
