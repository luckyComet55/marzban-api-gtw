package gatewayserver

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
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
		userProtoInfo := &contract.UserInfo{
			Username:    user.Username,
			Status:      user.Status.String(),
			UsedTraffic: uint64(user.UsedTraffic),
			ConfigUrls:  user.ConfigLinks,
		}
		if err := out.Send(userProtoInfo); err != nil {
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
	vlessProxySetting := cli.MarzbanProxySettings{
		Id:   uuid.New(),
		Flow: "",
	}
	nowTime := time.Now()
	userConf := cli.MarzbanUserConf{
		Username: createUserUnfo.Username,
		Proxies: map[cli.MarzbanProtocolType]cli.MarzbanProxySettings{
			cli.VlessProtocolType: vlessProxySetting,
		},
		Inbounds: map[cli.MarzbanProtocolType][]string{
			cli.VlessProtocolType: []string{createUserUnfo.ProxyProtocol},
		},
		DataLimit:              0,
		DataLimitResetStrategy: cli.NoResetStrategy,
		Expire:                 0,
		NextPlan:               nil,
		Note:                   "",
		Status:                 cli.ActiveStatus,
		OnHoldExpireDuration:   0,
		OnHoldTimeout:          cli.MarzbanDateTime(nowTime),
	}
	userInfo, err := s.marzbanPanelClient.CreateUser(userConf)
	if err != nil {
		errMsg := fmt.Sprintf("error on CreateUser call: %s", err.Error())
		s.logger.Error(errMsg)
		return nil, status.Error(codes.Internal, errMsg)
	}
	userResponse := &contract.UserInfo{
		Username:    userInfo.Username,
		Status:      userInfo.Status.String(),
		UsedTraffic: uint64(userInfo.UsedTraffic),
		ConfigUrls:  userInfo.ConfigLinks,
	}
	return userResponse, nil
}
