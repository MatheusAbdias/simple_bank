package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validator "github.com/go-playground/validator/v10"

	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.GET("/accounts/:id", server.getAccount)
	router.PATCH("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/accounts", server.createAccount)

	router.POST("/transfers", server.CreateTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error { return server.router.Run(address) }

func errorResponse(err error) gin.H { return gin.H{"error": err.Error()} }
