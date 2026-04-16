package api

import (
// 	"fmt"
    "net/http"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	// "github.com/go-playground/validator/v10"
	db "github.com/Oyinoye/bank_mini/db/sqlc"
// 	"github.com/Oyinoye/bank_mini/token"
// 	"github.com/Oyinoye/bank_mini/util"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  "0",
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		// errCode := db.ErrorCode(err)
		// if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
		// 	ctx.JSON(http.StatusForbidden, errorResponse(err))
		// 	return
		// }
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
