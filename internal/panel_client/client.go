package panelclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	pcl "github.com/luckyComet55/marzban-api-gtw/infra/panel_client"
	contract "github.com/luckyComet55/marzban-proto-contract/gen/go/contract"
)

type marzbanPanelClient struct {
	httpClient    *http.Client
	logger        *slog.Logger
	PanelBaseUrl  string
	PanelAuthPair marzbanPanelAuthPair
	panelAuthJwt  string
}

func NewMarzbanPanelClient(c MarzbanPanelClientConfig, logger *slog.Logger) pcl.MarzbanPanelClient {
	cli := &marzbanPanelClient{
		httpClient:    &http.Client{},
		PanelBaseUrl:  c.MarzbanBaseUrl,
		PanelAuthPair: marzbanPanelAuthPair{c.Username, c.Password},
		logger:        logger,
	}
	err := cli.getJwtAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func (cli *marzbanPanelClient) getJwtAccessToken() error {
	urlEncodedPost := url.Values{}
	urlEncodedPost.Set("grant_type", "password")
	urlEncodedPost.Set("username", cli.PanelAuthPair.Username)
	urlEncodedPost.Set("password", cli.PanelAuthPair.Password)

	authApiUrl, err := url.ParseRequestURI(cli.PanelBaseUrl)
	if err != nil {
		return err
	}
	authApiUrl.Path = "/api/admin/token"
	stringUrl := authApiUrl.String()

	request, err := http.NewRequest("POST", stringUrl, strings.NewReader(urlEncodedPost.Encode()))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	requestResult, err := cli.httpClient.Do(request)
	if err != nil {
		return err
	}
	if requestResult.StatusCode != 200 {
		errorMessage := fmt.Sprintf("unable to authorize, auth request returned with %d status code", requestResult.StatusCode)
		cli.logger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	defer requestResult.Body.Close()

	responseBodyStr, err := io.ReadAll(requestResult.Body)
	if err != nil {
		return err
	}

	var jwtUnmarshalledData marzbanJwtData

	if err := json.Unmarshal(responseBodyStr, &jwtUnmarshalledData); err != nil {
		return err
	}

	cli.panelAuthJwt = jwtUnmarshalledData.AccessToken

	return nil
}

func (cli *marzbanPanelClient) requestWrapper(request *http.Request) (*http.Response, error) {
	if request.Header.Get("Authorization") == "" {
		cli.logger.Warn(fmt.Sprintf("empty authorization for request to %s", request.URL.Path))
	}

	result, err := cli.httpClient.Do(request)
	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}
	cli.logger.Info(fmt.Sprintf("request to %s returned with result %s [%d]", request.URL.Path, result.Status, result.StatusCode))

	return result, nil
}

func (cli *marzbanPanelClient) requestAuthWrapper(request *http.Request) (*http.Response, error) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.panelAuthJwt))
	result, err := cli.requestWrapper(request)
	if err != nil {
		return nil, err
	}

	if result.StatusCode >= 200 && result.StatusCode < 399 {
		return result, nil
	}

	if result.StatusCode != 401 {
		return nil, fmt.Errorf("HTTP error %s [%d]", result.Status, result.StatusCode)
	}

	if err := cli.getJwtAccessToken(); err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.panelAuthJwt))
	result2, err := cli.requestWrapper(request)
	if err != nil {
		return nil, err
	}

	if result2.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %s [%d]", result2.Status, result2.StatusCode)
	}

	return result2, nil
}

func (cli *marzbanPanelClient) fetchUsers() (*marzbanUsersResponse, error) {
	usersApiUrl, err := url.ParseRequestURI(cli.PanelBaseUrl)

	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	usersApiUrl.Path = "/api/users"
	stringUrl := usersApiUrl.String()

	request, err := http.NewRequest("GET", stringUrl, nil)
	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	requestResult, err := cli.requestAuthWrapper(request)
	if err != nil {
		return nil, err
	}

	defer requestResult.Body.Close()

	responseBodyStr, err := io.ReadAll(requestResult.Body)
	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	var usersReponseUnmarshalled *marzbanUsersResponse

	if err := json.Unmarshal(responseBodyStr, &usersReponseUnmarshalled); err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	return usersReponseUnmarshalled, nil
}

func (cli *marzbanPanelClient) GetUsers() ([]*contract.UserInfo, error) {
	users, err := cli.fetchUsers()
	if err != nil {
		return nil, err
	}

	return users.Users, nil
}

func (cli *marzbanPanelClient) createUser(userInfo pcl.MarzbanUserConf) (*contract.UserInfo, error) {
	usersApiUrl, err := url.ParseRequestURI(cli.PanelBaseUrl)

	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	usersApiUrl.Path = "/api/user"
	stringUrl := usersApiUrl.String()

	cli.logger.Debug(fmt.Sprintf("%v", userInfo))

	postData, err := json.Marshal(userInfo)
	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	cli.logger.Debug(string(postData))

	request, err := http.NewRequest("POST", stringUrl, bytes.NewReader(postData))

	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	result, err := cli.requestAuthWrapper(request)
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()

	responseBodyStr, err := io.ReadAll(result.Body)
	if err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	cli.logger.Debug(string(responseBodyStr))

	var userCreateResult *contract.UserInfo
	if err := json.Unmarshal(responseBodyStr, &userCreateResult); err != nil {
		cli.logger.Error(err.Error())
		return nil, err
	}

	return userCreateResult, nil
}

func (cli *marzbanPanelClient) CreateUser(user *contract.CreateUserInfo) (*contract.UserInfo, error) {
	vlessProxySetting := pcl.MarzbanProxySettings{
		Id:   uuid.New(),
		Flow: "",
	}
	nowTime := time.Now()
	userData := pcl.MarzbanUserConf{
		Username: user.Username,
		Proxies: map[pcl.MarzbanProtocolType]pcl.MarzbanProxySettings{
			pcl.VlessProtocolType: vlessProxySetting,
		},
		Inbounds: map[pcl.MarzbanProtocolType][]string{
			pcl.VlessProtocolType: []string{user.ProxyProtocol},
		},
		DataLimit:              0,
		DataLimitResetStrategy: pcl.NoResetStrategy,
		Expire:                 0,
		NextPlan:               nil,
		Note:                   "",
		Status:                 pcl.ActiveStatus,
		OnHoldExpireDuration:   0,
		OnHoldTimeout:          nowTime,
	}

	return cli.createUser(userData)
}
