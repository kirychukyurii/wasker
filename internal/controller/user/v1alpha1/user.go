package v1alpha1

import (
	"context"
	"fmt"
	pbl "github.com/kirychukyurii/wasker/gen/go/lookup/v1alpha1"
	pb "github.com/kirychukyurii/wasker/gen/go/user/v1alpha1"
	"github.com/kirychukyurii/wasker/internal/pkg/log"
	"github.com/kirychukyurii/wasker/internal/service/user/v1alpha1"
)

// UserController will implement the service defined in protocol buffer definitions
type UserController struct {
	// userSvc.UnsafeUserServiceServer
	pb.UnimplementedUserServiceServer

	userService v1alpha1.UserService
	logger      log.Logger
}

func NewUserController(userService v1alpha1.UserService, logger log.Logger) UserController {
	fmt.Println("test from new controller")
	logger.Log.Info().Msg("test from new controller")
	return UserController{
		userService: userService,
		logger:      logger,
	}
}

func (a UserController) ReadUser(ctx context.Context, request *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	u, err := a.userService.ReadUser(ctx, request.Id)
	if err != nil {
		a.logger.Log.Warn().Err(err).Msg("Failed UserService method ReadUser()")
		return nil, err
	}

	r := &pb.ReadUserResponse{
		User: &pb.User{
			Id:       u.Id,
			Name:     u.Name,
			Email:    u.Email,
			Username: u.UserName,
			Password: "",
			Role: &pbl.ObjectId{
				Id:   u.Role.Id,
				Name: u.Role.Name,
			},
			CreatedAt: 0,
			CreatedBy: &pbl.ObjectId{
				Id:   u.CreatedBy.Id,
				Name: u.CreatedBy.Name,
			},
			UpdatedAt: 0,
			UpdatedBy: &pbl.ObjectId{
				Id:   u.UpdatedBy.Id,
				Name: u.UpdatedBy.Name,
			},
		},
	}

	return r, nil
}
