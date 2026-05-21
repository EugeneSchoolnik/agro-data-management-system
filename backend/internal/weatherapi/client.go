package weatherapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	login      string
	password   string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, login, password string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		login:      login,
		password:   password,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authResponse struct {
	Res struct {
		Token string `json:"Token"`
	} `json:"res"`
}

type StationLastDataRequest struct {
	StationID int `json:"stationid"`
}

type StationLastDataResponse struct {
	Res []struct {
		ParamID      int     `json:"ParamID"`
		ParamName    string  `json:"ParamName"`
		Unit         string  `json:"Unit"`
		StationParam int     `json:"StationParam"`
		Date         string  `json:"Date"`
		Value        float64 `json:"Value"`
	} `json:"res"`
}

type FieldLastDataRequest struct {
	FieldID int `json:"fieldid"`
}

type FieldLastDataResponse struct {
	Res []struct {
		ParamID      int    `json:"ParamID"`
		ParamName    string `json:"ParamName"`
		Unit         string `json:"Unit"`
		StationValue []struct {
			StationID    int     `json:"StationID"`
			StationName  string  `json:"StationName"`
			Date         string  `json:"Date"`
			Value        float64 `json:"Value"`
			StationParam int     `json:"StationParam"`
		} `json:"StationValue"`
	} `json:"res"`
}

func (c *Client) authenticate(ctx context.Context) error {
	if c.token != "" {
		return nil
	}

	request := authRequest{Login: c.login, Password: c.password}
	var resp authResponse
	if err := c.doPost(ctx, "/api/authorization", request, &resp); err != nil {
		return err
	}
	if resp.Res.Token == "" {
		return fmt.Errorf("authorization failed: empty token")
	}
	c.token = resp.Res.Token
	return nil
}

func (c *Client) StationLastDataGet(ctx context.Context, stationID int) (StationLastDataResponse, error) {
	if err := c.authenticate(ctx); err != nil {
		return StationLastDataResponse{}, err
	}

	var resp StationLastDataResponse
	request := StationLastDataRequest{StationID: stationID}
	if err := c.doPost(ctx, "/api/stationlastdataget", request, &resp); err != nil {
		return StationLastDataResponse{}, err
	}
	return resp, nil
}

func (c *Client) FieldLastDataGet(ctx context.Context, fieldID int) (FieldLastDataResponse, error) {
	if err := c.authenticate(ctx); err != nil {
		return FieldLastDataResponse{}, err
	}

	var resp FieldLastDataResponse
	request := FieldLastDataRequest{FieldID: fieldID}
	if err := c.doPost(ctx, "/api/fieldlastdataget", request, &resp); err != nil {
		return FieldLastDataResponse{}, err
	}
	return resp, nil
}

func (c *Client) doPost(ctx context.Context, path string, request interface{}, response interface{}) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("weather api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		if path == "/api/authorization" {
			content, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("weather api authorization failed: %s", strings.TrimSpace(string(content)))
		}
		c.token = ""
		if err := c.authenticate(ctx); err != nil {
			return err
		}
		return c.doPost(ctx, path, request, response)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		content, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("weather api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(content)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read weather api response: %w", err)
	}

	if err := json.Unmarshal(data, response); err != nil {
		return fmt.Errorf("failed to decode weather api response: %w", err)
	}

	return nil
}
