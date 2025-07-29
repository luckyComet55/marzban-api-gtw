package panelclient

import "github.com/luckyComet55/marzban-proto-contract/gen/go/contract"

type MarzbanPanelClient interface {
	GetUsers() ([]*contract.UserInfo, error)
	CreateUser(*contract.CreateUserInfo) (*contract.UserInfo, error)
}
