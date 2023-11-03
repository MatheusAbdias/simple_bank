package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type responseRenewAccessTokenDTO struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

var InvalidUserSession = errors.New("invalid user session")

var MismatchSessionToken = errors.New("mismatched session and token")

var ErrExpiredToken = errors.New("token has expired")

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, InvalidUserSession)
		return
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, MismatchSessionToken)
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrExpiredToken))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := responseRenewAccessTokenDTO{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}
