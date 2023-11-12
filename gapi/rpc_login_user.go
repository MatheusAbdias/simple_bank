package gapi

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
	"github.com/MatheusAbdias/go_simple_bank/pb"
	"github.com/MatheusAbdias/go_simple_bank/util"
	"github.com/MatheusAbdias/go_simple_bank/validators"
)

func (server *Server) LoginUser(
	ctx context.Context,
	request *pb.LoginUserRequest,
) (*pb.LoginUserResponse, error) {
	violations := ValidateLoginUserRequest(request)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, request.GetUsername())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "fail to fetch user")
	}

	err = util.CheckPassword(request.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials provided")
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to create access token")
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	metaData := server.extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
		UserAgent:    metaData.UserAgent,
		ClientIp:     metaData.ClientIP,
		IsBlocked:    false,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	response := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		User:                  convertUser(user),
	}
	return response, nil
}

func ValidateLoginUserRequest(
	req *pb.LoginUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername((req.GetUsername())); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validators.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
