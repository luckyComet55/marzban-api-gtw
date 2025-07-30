package panelclient

import (
	pcl "github.com/luckyComet55/marzban-api-gtw/infra/panel_client"
)

type marzbanPanelAuthPair struct {
	Username string
	Password string
}

type marzbanUsersResponse struct {
	Users []*pcl.MarzbanUserInfo `json:"users"`
	Total uint64                 `json:"total"`
}

type marzbanJwtData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
