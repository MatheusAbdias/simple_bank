package gapi

import (
	"context"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	db "github.com/MatheusAbdias/go_simple_bank/db/sqlc"
	"github.com/MatheusAbdias/go_simple_bank/pb"
	"github.com/MatheusAbdias/go_simple_bank/util"
)

func (server *Server) CreateUser(
	ctx context.Context,
	request *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       request.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       request.GetFullName(),
		Email:          request.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return response, nil
}
