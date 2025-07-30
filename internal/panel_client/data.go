package panelclient

import (
	contract "github.com/luckyComet55/marzban-proto-contract/gen/go/contract"
)

type marzbanPanelAuthPair struct {
	Username string
	Password string
}

type marzbanUsersResponse struct {
	Users []*contract.UserInfo `json:"users"`
	Total uint64               `json:"total"`
}

type marzbanJwtData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
