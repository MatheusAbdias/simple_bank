package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"

	"github.com/MatheusAbdias/go_simple_bank/config"
	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
	"github.com/MatheusAbdias/go_simple_bank/token"
)

type Server struct {
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.createUser)
	router.POST("/users/refresh", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.PATCH("/accounts/:id", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/accounts", server.createAccount)

	authRoutes.POST("/transfers", server.CreateTransfer)

	server.router = router

}

func (server *Server) Start(address string) error { return server.router.Run(address) }

func errorResponse(err error) gin.H { return gin.H{"error": err.Error()} }
