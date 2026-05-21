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

type FieldsGetRequest struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type FieldInfo struct {
	FieldID     int     `json:"FieldID"`
	Number      string  `json:"Number"`
	Polygon     string  `json:"Polygon"`
	Status      int     `json:"Status"`
	CultureID   int     `json:"CultureID"`
	CompanyID   int     `json:"CompanyID"`
	CompanyName string  `json:"CompanyName"`
	Square      float64 `json:"Square"`
	StationsNum int     `json:"StationsNum"`
}

type FieldsGetResponse struct {
	Res []FieldInfo `json:"res"`
}

type FieldMeteoParamsRequest struct {
	FieldID int `json:"fieldid"`
}

type FieldMeteoParamsResponse struct {
	Res []struct {
		ParamID   int    `json:"ParamID"`
		ParamName string `json:"ParamName"`
		Unit      string `json:"Unit"`
		Stations  []struct {
			ID           int    `json:"ID"`
			Name         string `json:"Name"`
			StationParam int    `json:"StationParam"`
		} `json:"Stations"`
	} `json:"res"`
}

type FieldForecastRequest struct {
	FieldID   int    `json:"fieldid"`
	ApiSource string `json:"apisource"`
}

type FieldForecastResponse struct {
	Res json.RawMessage `json:"res"`
}

type StationMeteoHistoryRequest struct {
	StationID    int    `json:"stationid"`
	ParamID      int    `json:"paramid"`
	StationParam int    `json:"stationparam"`
	FieldID      int    `json:"fieldid,omitempty"`
	DateStart    string `json:"datestart"`
	DateFinish   string `json:"datefinish"`
	AvgPeriod    string `json:"avgperiod"`
	AvgType      string `json:"avgtype"`
}

type StationMeteoHistoryResponse struct {
	Res []struct {
		Date  string  `json:"Date"`
		Value float64 `json:"Value"`
	} `json:"res"`
}

type StationGeoDataRequest struct {
	StationID  int    `json:"stationid"`
	DateStart  string `json:"datestart"`
	DateFinish string `json:"datefinish"`
}

type StationGeoDataResponse struct {
	Res []struct {
		PointTime string  `json:"PointTime"`
		Lat       float64 `json:"Lat"`
		Lon       float64 `json:"Lon"`
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

func (c *Client) FieldsGet(ctx context.Context, page, limit int, search string) (FieldsGetResponse, error) {
	if err := c.authenticate(ctx); err != nil {
		return FieldsGetResponse{}, err
	}

	var resp FieldsGetResponse
	request := FieldsGetRequest{Page: page, Limit: limit, Search: search}
	if err := c.doPost(ctx, "/api/fieldsget", request, &resp); err != nil {
		return FieldsGetResponse{}, err
	}
	return resp, nil
}

func (c *Client) FieldMeteoParamsGet(ctx context.Context, fieldID int) (FieldMeteoParamsResponse, error) {
	if err := c.authenticate(ctx); err != nil {
		return FieldMeteoParamsResponse{}, err
	}

	var resp FieldMeteoParamsResponse
	request := FieldMeteoParamsRequest{FieldID: fieldID}
	if err := c.doPost(ctx, "/api/fieldmeteoparamsget", request, &resp); err != nil {
		return FieldMeteoParamsResponse{}, err
	}
	return resp, nil
}

func (c *Client) FieldForecastGet(ctx context.Context, fieldID int) (FieldForecastResponse, error) {
	if err := c.authenticate(ctx); err != nil {
		return FieldForecastResponse{}, err
	}

	var resp FieldForecastResponse
	request := FieldForecastRequest{FieldID: fieldID, ApiSource: "worldweatheronline"}
	if err := c.doPost(ctx, "/api/fieldforecastget", request, &resp); err != nil {
		return FieldForecastResponse{}, err
	}
	return resp, nil
}

func (c *Client) StationMeteoHistoryGet(ctx context.Context, req StationMeteoHistoryRequest) (StationMeteoHistoryResponse, error) {
	if err := c.authenticate(ctx); err != nil {
		return StationMeteoHistoryResponse{}, err
	}

	var resp StationMeteoHistoryResponse
	if err := c.doPost(ctx, "/api/stationmeteohistoryget", req, &resp); err != nil {
		return StationMeteoHistoryResponse{}, err
	}
	return resp, nil
}

func (c *Client) StationGeoDataGet(ctx context.Context, req StationGeoDataRequest) (StationGeoDataResponse, error) {
	if err := c.authenticate(ctx); err != nil {
		return StationGeoDataResponse{}, err
	}

	var resp StationGeoDataResponse
	if err := c.doPost(ctx, "/api/stationgeodataget", req, &resp); err != nil {
		return StationGeoDataResponse{}, err
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
		req.Header.Set("Authorization", c.token)
	}
	req.Header.Set("Accept-Language", "ua")

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
