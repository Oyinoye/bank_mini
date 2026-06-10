package gapi

import (
	"context"

	db "github.com/Oyinoye/bank_mini/db/sqlc"
	"github.com/Oyinoye/bank_mini/pb"
	"github.com/Oyinoye/bank_mini/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// violations := validateCreateUserRequest(req)
	// if violations != nil {
	// 	return nil, invalidArgumentError(violations)
	// }

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
        Username:       req.GetUsername(),
        HashedPassword: hashedPassword,
        FullName:       req.GetFullName(),
        Email:          req.GetEmail(),
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


    // arg := db.CreateUserTxParams{
	// 	CreateUserParams: db.CreateUserParams{
	// 		Username:       req.GetUsername(),
	// 		HashedPassword: hashedPassword,
	// 		FullName:       req.GetFullName(),
	// 		Email:          req.GetEmail(),
	// 	},
	// 	AfterCreate: func(user db.User) error {
	// 		taskPayload := &worker.PayloadSendVerifyEmail{
	// 			Username: user.Username,
	// 		}
	// 		opts := []asynq.Option{
	// 			asynq.MaxRetry(10),
	// 			asynq.ProcessIn(10 * time.Second),
	// 			asynq.Queue(worker.QueueCritical),
	// 		}

	// 		return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
	// 	},
	// }

	// txResult, err := server.store.CreateUserTx(ctx, arg)
	// if err != nil {
	// 	if db.ErrorCode(err) == db.UniqueViolation {
	// 		return nil, status.Error(codes.AlreadyExists, err.Error())
	// 	}
	// 	return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	// }

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
