package gapi

import (
	"context"
	"database/sql"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
	"github.com/MatheusAbdias/go_simple_bank/pb"
	"github.com/MatheusAbdias/go_simple_bank/util"
	"github.com/MatheusAbdias/go_simple_bank/validators"
)

func (server *Server) UpdateUser(
	ctx context.Context,
	request *pb.UpdateUserRequest,
) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := ValidateUpdateUserRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if authPayload.Username != request.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user")
	}
	arg := db.UpdateUserParams{
		Username: request.GetUsername(),
		FullName: sql.NullString{
			String: request.GetFullName(),
			Valid:  true,
		},
		Email: sql.NullString{
			String: request.GetEmail(),
			Valid:  true,
		},
	}

	if request.Password != nil {
		hashedPassword, err := util.HashPassword(request.GetPassword())

		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	response := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return response, nil
}

func ValidateUpdateUserRequest(
	req *pb.UpdateUserRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validators.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Password != nil {
		if err := validators.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if req.FullName != nil {
		if err := validators.ValidateFullName(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := validators.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
