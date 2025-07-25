package panelclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	contract "github.com/luckyComet55/marzban-api-gtw/contract"
)

type MarzbanPanelClient interface {
	GetUsers() ([]*contract.UserInfo, error)
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

type marzbanPanelClientImpl struct {
	httpClient    *http.Client
	logger        *log.Logger
	PanelBaseUrl  string
	PanelAuthPair marzbanPanelAuthPair
	panelAuthJwt  string
}

func NewMarzbanPanelClient(panelBaseUrl, panelUsername, panelPassword string) MarzbanPanelClient {
	cli := &marzbanPanelClientImpl{
		httpClient:    &http.Client{},
		PanelBaseUrl:  panelBaseUrl,
		PanelAuthPair: marzbanPanelAuthPair{panelUsername, panelPassword},
	}
	err := cli.getJwtAccessToken()
	if err != nil {
		log.Panic(err)
	}
	return cli
}

func (cli *marzbanPanelClientImpl) getJwtAccessToken() error {
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
		errorMessage := fmt.Sprintf("auth request resturned with %d status code", requestResult.StatusCode)
		cli.logger.Println(errorMessage)
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

func (cli *marzbanPanelClientImpl) fetchUsers() (*http.Response, error) {
	usersApiUrl, err := url.ParseRequestURI(cli.PanelBaseUrl)

	if err != nil {
		cli.logger.Println(err)
		return nil, err
	}

	usersApiUrl.Path = "/api/users"
	stringUrl := usersApiUrl.String()

	request, err := http.NewRequest("GET", stringUrl, nil)

	if err != nil {
		cli.logger.Println(err)
		return nil, err
	}

	bearerTokenString := fmt.Sprintf("Bearer %s", cli.panelAuthJwt)
	request.Header.Add("Authorization", bearerTokenString)

	requestResult, err := cli.httpClient.Do(request)

	if err != nil {
		cli.logger.Println(err)
		return nil, err
	}

	return requestResult, nil
}

func (cli *marzbanPanelClientImpl) GetUsers() ([]*contract.UserInfo, error) {
	if cli.panelAuthJwt == "" {
		if err := cli.getJwtAccessToken(); err != nil {
			cli.logger.Printf("error while retrieving access token: %s\n", err)
			return nil, err
		}
	}

	requestResult, err := cli.fetchUsers()
	if err != nil {
		cli.logger.Printf("error while fetching users: %s\n", err)
		return nil, err
	}

	if requestResult.StatusCode == 401 {
		cli.logger.Println("access token expired, retrieving new one")
		if err := cli.getJwtAccessToken(); err != nil {
			cli.logger.Printf("error while retrieving access token: %s\n", err)
			return nil, err
		}
		return cli.GetUsers()
	} else if requestResult.StatusCode != 200 {
		errorMessage := fmt.Sprintf("unhandled status code from Marzban: [%d] %s", requestResult.StatusCode, requestResult.Status)
		cli.logger.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	defer requestResult.Body.Close()

	responseBodyStr, err := io.ReadAll(requestResult.Body)
	if err != nil {
		cli.logger.Println(err)
		return nil, err
	}

	var usersReponseUnmarshalled *marzbanUsersResponse

	if err := json.Unmarshal(responseBodyStr, &usersReponseUnmarshalled); err != nil {
		cli.logger.Println(err)
		return nil, err
	}

	return usersReponseUnmarshalled.Users, nil
}
