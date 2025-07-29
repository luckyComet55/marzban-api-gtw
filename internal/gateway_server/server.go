package gatewayserver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/golang/protobuf/ptypes/empty"
	cli "github.com/luckyComet55/marzban-api-gtw/infra/panel_client"
	"github.com/luckyComet55/marzban-proto-contract/gen/go/contract"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type marzbanManagementPanelServer struct {
	logger             *slog.Logger
	marzbanPanelClient cli.MarzbanPanelClient
	contract.UnimplementedMarzbanManagementPanelServer
}

func Register(gRPCServer *grpc.Server, panelClient cli.MarzbanPanelClient, logger *slog.Logger) {
	contract.RegisterMarzbanManagementPanelServer(gRPCServer, &marzbanManagementPanelServer{
		logger:             logger,
		marzbanPanelClient: panelClient,
	})
}

func (s *marzbanManagementPanelServer) ListUsers(
	in *empty.Empty,
	out grpc.ServerStreamingServer[contract.UserInfo],
) error {
	userInfo, err := s.marzbanPanelClient.GetUsers()
	if err != nil {
		errMsg := fmt.Sprintf("error on ListUsers call: %s", err.Error())
		s.logger.Error(errMsg)
		return status.Error(codes.Internal, errMsg)
	}

	for _, user := range userInfo {
		if err := out.Send(user); err != nil {
			errMsg := fmt.Sprintf("error on send in ListUsers call: %s", err.Error())
			s.logger.Error(errMsg)
			return status.Error(codes.Internal, errMsg)
		}
	}
	return nil
}

func (s *marzbanManagementPanelServer) CreateUser(
	ctx context.Context,
	createUserUnfo *contract.CreateUserInfo,
) (*contract.UserInfo, error) {
	userInfo, err := s.marzbanPanelClient.CreateUser(createUserUnfo)
	if err != nil {
		errMsg := fmt.Sprintf("error on CreateUser call: %s", err.Error())
		s.logger.Error(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}
	return userInfo, nil
}
