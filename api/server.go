package api

import (
	"database/sql"
	"fmt"
	db "github.com/Bakhram74/small_bank/db/sqlc"
	"github.com/Bakhram74/small_bank/token"
	"github.com/Bakhram74/small_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
func (s *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := s.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if currency != account.Currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return false
	}
	return true
}

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func (s *Server) setupRouter() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	router.POST("/users/", s.createUser)
	router.POST("/users/login", s.loginUser)
	router.GET("/users", s.getUser)

	router.POST("/accounts", s.createAccount)
	router.GET("/accounts/:id", s.getAccount)
	router.GET("/accounts", s.listAccounts)

	router.POST("/transfers", s.createTransfer)

	s.router = router
}
func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
