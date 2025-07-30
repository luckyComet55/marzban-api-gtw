package panelclient

type MarzbanPanelClient interface {
	GetUsers() ([]*MarzbanUserInfo, error)
	CreateUser(MarzbanUserConf) (*MarzbanUserInfo, error)
}
