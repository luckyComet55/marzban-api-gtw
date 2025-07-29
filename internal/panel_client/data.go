package panelclient

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	contract "github.com/luckyComet55/marzban-proto-contract/gen/go/contract"
)

type marzbanLimitResetStrategy int
type marzbanUserStatus int
type marzbanProtocolType int

const (
	noResetStrategy marzbanLimitResetStrategy = iota
	dailyResetStrategy
	weekyResetStrategy
	monthlyResetStrategy
	yearlyResetStrategy
)

const (
	activeStatus marzbanUserStatus = iota
	onHoldStatus
)

const (
	vmessProtocolType marzbanProtocolType = iota
	vlessProtocolType
	trojanProtocolType
	shadowsocksProtocolType
)

var strategyName = map[marzbanLimitResetStrategy]string{
	noResetStrategy:      "no_reset",
	dailyResetStrategy:   "day",
	weekyResetStrategy:   "week",
	monthlyResetStrategy: "month",
	yearlyResetStrategy:  "year",
}

var userStatusName = map[marzbanUserStatus]string{
	activeStatus: "active",
	onHoldStatus: "on_hold",
}

var protocolTypeName = map[marzbanProtocolType]string{
	vmessProtocolType:       "vmess",
	vlessProtocolType:       "vless",
	trojanProtocolType:      "trojan",
	shadowsocksProtocolType: "shadowsocks",
}

func (strategy marzbanLimitResetStrategy) String() string {
	if len(strategyName) < int(strategy) {
		return ""
	}
	return strategyName[strategy]
}

func (strategy marzbanLimitResetStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(strategy.String())
}

func (strategy marzbanLimitResetStrategy) MarshalText() ([]byte, error) {
	return []byte(strategy.String()), nil
}

func (userStatus marzbanUserStatus) String() string {
	if len(userStatusName) < int(userStatus) {
		return ""
	}
	return userStatusName[userStatus]
}

func (userStatus marzbanUserStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(userStatus.String())
}

func (userStatus marzbanUserStatus) MarshalText() ([]byte, error) {
	return []byte(userStatus.String()), nil
}

func (protocolType marzbanProtocolType) String() string {
	if len(protocolTypeName) < int(protocolType) {
		return ""
	}
	return protocolTypeName[protocolType]
}

func (protocolType marzbanProtocolType) MarshalJSON() ([]byte, error) {
	return json.Marshal(protocolType.String())
}

func (protocolType marzbanProtocolType) MarshalText() ([]byte, error) {
	return []byte(protocolType.String()), nil
}

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

type marzbanProxySettings struct {
	Id   uuid.UUID `json:"id"`
	Flow string    `json:"flow,omitempty"`
}

type marzbanUserConf struct {
	DataLimit              uint                                         `json:"data_limit"`
	DataLimitResetStrategy marzbanLimitResetStrategy                    `json:"data_limit_reset_strategy"`
	Expire                 uint                                         `json:"expire"`
	Inbounds               map[marzbanProtocolType][]string             `json:"inbounds"`
	Proxies                map[marzbanProtocolType]marzbanProxySettings `json:"proxies"`
	NextPlan               map[string]string                            `json:"next_plan"`
	Note                   string                                       `json:"note"`
	OnHoldExpireDuration   uint                                         `json:"on_hold_expire_duration"`
	OnHoldTimeout          time.Time                                    `json:"on_hold_timeout"`
	Status                 marzbanUserStatus                            `json:"status"`
	Username               string                                       `json:"username"`
}
