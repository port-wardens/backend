package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jayendramadaram/port-wardens/model"
)

var envVar string

type LoginResponse struct {
	Token string `json:"token"`
}

type Client interface {
	// Doesnt need JWT

	// Auth Routes
	SignUp(model.CreateUser) error
	Login(model.LoginUser) (LoginResponse, error)

	// JWT Routes
}

type client struct {
	url string
	jwt string
}

func NewPortWardenClient(url string) Client {
	return &client{url: url}
}

func (c *client) SignUp(req model.CreateUser) error {
	if _, err := c.submitRequest(req, "/signup", "POST"); err != nil {
		return err
	}
	return nil
}

func (c *client) Login(req model.LoginUser) (LoginResponse, error) {
	resp, err := c.submitRequest(req, "/login", "POST")
	if err != nil {
		return LoginResponse{}, err
	}

	var login LoginResponse
	if err := json.NewDecoder(resp).Decode(&login); err != nil {
		return LoginResponse{}, err
	}
	c.jwt = login.Token
	return login, nil
}

func (c *client) submitRequest(req interface{}, endpoint string, method string) (io.ReadCloser, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return nil, err
	}

	var resp *http.Request
	var err error

	switch method {
	case "GET":
		resp, err = http.NewRequest(http.MethodGet, c.url+endpoint, nil)
	case "POST":
		resp, err = http.NewRequest(http.MethodPost, c.url+endpoint, buf)
	case "PUT":
		resp, err = http.NewRequest(http.MethodPut, c.url+endpoint, buf)
	default:
		return nil, fmt.Errorf("invalid method %v", method)
	}
	if err != nil {
		return nil, err
	}

	resp.Header.Set("Content-Type", "application/json")
	if c.jwt != "" {
		// fmt.Println("jwt: ", c.jwt)
		resp.Header.Set("Authorization", c.jwt)
	}

	client := &http.Client{}
	// fmt.Println("respunse: ", resp, "err: ", err)
	response, err := client.Do(resp)
	if err != nil {
		return nil, err
	}
	// defer response.Body.Close()
	fmt.Println("response status: ", response.StatusCode)
	if response.StatusCode < 200 || response.StatusCode > 300 {
		errObj := struct {
			Error string `json:"error"`
		}{}

		if err := json.NewDecoder(response.Body).Decode(&errObj); err != nil {
			errMsg, err := io.ReadAll(response.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read the error message %v", err)
			}
			return nil, fmt.Errorf("failed to decode the error %v", string(errMsg))
		}
		return nil, fmt.Errorf("request failed with status %v", errObj.Error)
	}
	return response.Body, err
}
