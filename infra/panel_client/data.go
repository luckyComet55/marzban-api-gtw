package panelclient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MarzbanLimitResetStrategy int
type MarzbanUserStatus int
type MarzbanProtocolType int
type MarzbanDateTime time.Time

func (dateTime *MarzbanDateTime) UnmarshalJSON(data []byte) error {
	var dateTimeStr string
	if err := json.Unmarshal(data, &dateTimeStr); err != nil {
		return err
	}
	if dateTimeStr == "" {
		dateTime = nil
		return nil
	}
	var err error
	var dateTimeTime time.Time
	dateTimeTime, err = time.Parse("2006-01-02T15:04:05.000000", dateTimeStr)
	if err != nil {
		return err
	}
	*dateTime = MarzbanDateTime(dateTimeTime)
	return nil
}

func (dateTime MarzbanDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dateTime))
}

const (
	NoResetStrategy MarzbanLimitResetStrategy = iota
	DailyResetStrategy
	WeekyResetStrategy
	MonthlyResetStrategy
	YearlyResetStrategy
)

const (
	ActiveStatus MarzbanUserStatus = iota
	OnHoldStatus
)

const (
	VmessProtocolType MarzbanProtocolType = iota
	VlessProtocolType
	TrojanProtocolType
	ShadowsocksProtocolType
)

var strategyName = map[MarzbanLimitResetStrategy]string{
	NoResetStrategy:      "no_reset",
	DailyResetStrategy:   "day",
	WeekyResetStrategy:   "week",
	MonthlyResetStrategy: "month",
	YearlyResetStrategy:  "year",
}

var userStatusName = map[MarzbanUserStatus]string{
	ActiveStatus: "active",
	OnHoldStatus: "on_hold",
}

var protocolTypeName = map[MarzbanProtocolType]string{
	VmessProtocolType:       "vmess",
	VlessProtocolType:       "vless",
	TrojanProtocolType:      "trojan",
	ShadowsocksProtocolType: "shadowsocks",
}

func (strategy MarzbanLimitResetStrategy) String() string {
	value, ok := strategyName[strategy]
	if !ok {
		panic(fmt.Sprintf("bad String() call for MarzbanLimitResetStrategy enum with value %v", int(strategy)))
	}
	return value
}

func ToMarzbanLimitResetStrategy(value string) (MarzbanLimitResetStrategy, error) {
	switch value {
	case "no_reset":
		return NoResetStrategy, nil
	case "day":
		return DailyResetStrategy, nil
	case "week":
		return WeekyResetStrategy, nil
	case "month":
		return MonthlyResetStrategy, nil
	case "year":
		return YearlyResetStrategy, nil
	default:
		return -1, fmt.Errorf("unable to convert value '%s' to MarzbanLimitResetStrategy enum", value)
	}
}

func (strategy MarzbanLimitResetStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(strategy.String())
}

func (strategy MarzbanLimitResetStrategy) MarshalText() ([]byte, error) {
	return []byte(strategy.String()), nil
}

func (strategy *MarzbanLimitResetStrategy) UnmarshalJSON(data []byte) error {
	var strategyStr string
	if err := json.Unmarshal(data, &strategyStr); err != nil {
		return err
	}
	var err error
	*strategy, err = ToMarzbanLimitResetStrategy(strategyStr)
	return err
}

func (strategy *MarzbanLimitResetStrategy) UnmarshalText(text []byte) error {
	strategyStr := string(text)
	var err error
	*strategy, err = ToMarzbanLimitResetStrategy(strategyStr)
	return err
}

func (userStatus MarzbanUserStatus) String() string {
	value, ok := userStatusName[userStatus]
	if !ok {
		panic(fmt.Sprintf("bad String() call for MarzbanUserStatus enum with value %v", int(userStatus)))
	}
	return value
}

func ToMarzbanUserStatus(value string) (MarzbanUserStatus, error) {
	switch value {
	case "active":
		return ActiveStatus, nil
	case "on_hold":
		return OnHoldStatus, nil
	default:
		return -1, fmt.Errorf("unable to convert value '%s' to MarzbanUserStatus enum", value)
	}
}

func (userStatus MarzbanUserStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(userStatus.String())
}

func (userStatus MarzbanUserStatus) MarshalText() ([]byte, error) {
	return []byte(userStatus.String()), nil
}

func (userStatus *MarzbanUserStatus) UnmarshalJSON(data []byte) error {
	var userStatusStr string
	if err := json.Unmarshal(data, &userStatusStr); err != nil {
		return err
	}
	var err error
	*userStatus, err = ToMarzbanUserStatus(userStatusStr)
	return err
}

func (userStatus *MarzbanUserStatus) UnmarshalText(text []byte) error {
	userStatusStr := string(text)
	var err error
	*userStatus, err = ToMarzbanUserStatus(userStatusStr)
	return err
}

func (protocolType MarzbanProtocolType) String() string {
	value, ok := protocolTypeName[protocolType]
	if !ok {
		panic(fmt.Sprintf("bad String() call for MarzbanProtocolType enum with value %v", int(protocolType)))
	}
	return value
}

func ToMarzbanProtocolType(value string) (MarzbanProtocolType, error) {
	switch value {
	case "vmess":
		return VmessProtocolType, nil
	case "vless":
		return VlessProtocolType, nil
	case "trojan":
		return TrojanProtocolType, nil
	case "shadowsocks":
		return ShadowsocksProtocolType, nil
	default:
		return -1, fmt.Errorf("unable to convert value '%s' to MarzbanProtocolType enum", value)
	}
}

func (protocolType MarzbanProtocolType) MarshalJSON() ([]byte, error) {
	return json.Marshal(protocolType.String())
}

func (protocolType MarzbanProtocolType) MarshalText() ([]byte, error) {
	return []byte(protocolType.String()), nil
}

func (protocolType *MarzbanProtocolType) UnmarshalJSON(data []byte) error {
	var protocolTypeStr string
	if err := json.Unmarshal(data, &protocolTypeStr); err != nil {
		return err
	}
	var err error
	*protocolType, err = ToMarzbanProtocolType(protocolTypeStr)
	return err
}

func (protocolType *MarzbanProtocolType) UnmarshalText(text []byte) error {
	protocolTypeStr := string(text)
	var err error
	*protocolType, err = ToMarzbanProtocolType(protocolTypeStr)
	return err
}

type MarzbanProxySettings struct {
	Id   uuid.UUID `json:"id"`
	Flow string    `json:"flow,omitempty"`
}

type MarzbanUserConf struct {
	DataLimit              uint                                         `json:"data_limit"`
	DataLimitResetStrategy MarzbanLimitResetStrategy                    `json:"data_limit_reset_strategy"`
	Expire                 uint                                         `json:"expire"`
	Inbounds               map[MarzbanProtocolType][]string             `json:"inbounds"`
	Proxies                map[MarzbanProtocolType]MarzbanProxySettings `json:"proxies"`
	NextPlan               map[string]string                            `json:"next_plan"`
	Note                   string                                       `json:"note"`
	OnHoldExpireDuration   uint                                         `json:"on_hold_expire_duration"`
	OnHoldTimeout          MarzbanDateTime                              `json:"on_hold_timeout"`
	Status                 MarzbanUserStatus                            `json:"status"`
	Username               string                                       `json:"username"`
}

type MarzbanUserInfo struct {
	MarzbanUserConf
	ConfigLinks []string        `json:"links"`
	UsedTraffic int             `json:"used_traffic"`
	CreatedAt   MarzbanDateTime `json:"created_at"`
}
