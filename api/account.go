package api

import (
	"github.com/gin-gonic/gin"
	db "goBank/db/sqlc"
	"net/http"
)

type createAccountRequest struct {
	// binding: required added to validate the json via Gin
	Owner string `json:"owner" binding:"required"`
	// Removing this param as it needs to default to 0 when setup
	//Balance  int64  `json:"balance"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	//Handling bad data responses
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
