//go:build !unit
// +build !unit

package weatherapi

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getLiveClient(t *testing.T) *Client {
	t.Helper()
	loadDotEnv(t)
	if os.Getenv("WEATHER_API_INTEGRATION") != "1" {
		t.Skip("skip live weather API tests unless WEATHER_API_INTEGRATION=1")
	}

	login := os.Getenv("WEATHER_API_LOGIN")
	password := os.Getenv("WEATHER_API_PASSWORD")
	if login == "" || password == "" {
		t.Fatal("WEATHER_API_LOGIN and WEATHER_API_PASSWORD must be set for live integration tests")
	}

	baseURL := os.Getenv("WEATHER_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.meteotrek.ua"
	}

	return NewClient(baseURL, login, password)
}

func loadDotEnv(t *testing.T) {
	t.Helper()
	if os.Getenv("WEATHER_API_INTEGRATION") != "" && os.Getenv("WEATHER_API_LOGIN") != "" && os.Getenv("WEATHER_API_PASSWORD") != "" {
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Logf("loadDotEnv: cannot get current working dir: %v", err)
		return
	}

	for i := 0; i < 6; i++ {
		path := filepath.Join(cwd, ".env")
		if _, err := os.Stat(path); err == nil {
			if err := parseDotEnv(path); err != nil {
				t.Logf("loadDotEnv: parse error: %v", err)
			}
			return
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}
}

func parseDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")
		if key == "" {
			continue
		}
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
	return scanner.Err()
}

func logResponsePartial(t *testing.T, prefix string, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		t.Logf("%s: failed to marshal response for log: %v", prefix, err)
		return
	}
	message := string(data)
	if len(message) > 400 {
		message = message[:400] + "..."
	}
	t.Logf("%s: %s", prefix, message)
}

func TestClient_LiveFieldsGet(t *testing.T) {
	client := getLiveClient(t)

	resp, err := client.FieldsGet(context.Background(), 0, 5, "")
	assert.NoError(t, err)
	assert.NotNil(t, resp.Res, "expected live fields response to be present")
	assert.GreaterOrEqual(t, len(resp.Res), 1, "expected at least one field in live fields response")
	for i, field := range resp.Res {
		if i >= 3 {
			break
		}
		t.Logf("field %d id=%d name=%s stations=%d", i+1, field.FieldID, field.Number, field.StationsNum)
	}
	logResponsePartial(t, "fieldsget response", resp.Res)
}

func TestClient_LiveFieldMeteoParamsGet(t *testing.T) {
	client := getLiveClient(t)

	fieldIDStr := os.Getenv("WEATHER_API_FIELD_ID")
	if fieldIDStr == "" {
		t.Skip("skip live field meteo params test unless WEATHER_API_FIELD_ID is set")
	}

	fieldID, err := strconv.Atoi(fieldIDStr)
	assert.NoError(t, err)

	resp, err := client.FieldMeteoParamsGet(context.Background(), fieldID)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Res, "expected live field meteo params response to be present")
	assert.GreaterOrEqual(t, len(resp.Res), 1, "expected at least one meteo param in live field response")
	for i, param := range resp.Res {
		if i >= 3 {
			break
		}
		t.Logf("field param %d id=%d name=%s stations=%d", i+1, param.ParamID, param.ParamName, len(param.Stations))
	}
	logResponsePartial(t, "fieldmeteoparamsget response", resp.Res)
}

func TestClient_LiveFieldForecastGet(t *testing.T) {
	client := getLiveClient(t)

	fieldIDStr := os.Getenv("WEATHER_API_FIELD_ID")
	if fieldIDStr == "" {
		t.Skip("skip live field forecast test unless WEATHER_API_FIELD_ID is set")
	}

	fieldID, err := strconv.Atoi(fieldIDStr)
	assert.NoError(t, err)

	resp, err := client.FieldForecastGet(context.Background(), fieldID)
	assert.NoError(t, err)
	logResponsePartial(t, "fieldforecastget response", resp.Res)
}

func TestClient_LiveStationMeteoHistoryGet(t *testing.T) {
	client := getLiveClient(t)

	stationIDStr := os.Getenv("WEATHER_API_STATION_ID")
	fieldIDStr := os.Getenv("WEATHER_API_FIELD_ID")
	if stationIDStr == "" || fieldIDStr == "" {
		t.Skip("skip live station meteo history test unless both WEATHER_API_STATION_ID and WEATHER_API_FIELD_ID are set")
	}

	stationID, err := strconv.Atoi(stationIDStr)
	assert.NoError(t, err)
	fieldID, err := strconv.Atoi(fieldIDStr)
	assert.NoError(t, err)

	paramsResp, err := client.FieldMeteoParamsGet(context.Background(), fieldID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(paramsResp.Res), 1, "expected at least one field meteo param to build history request")
	paramID := paramsResp.Res[0].ParamID
	if len(paramsResp.Res[0].Stations) == 0 {
		t.Skip("skip live station meteo history test because first field meteo param has no station mapping")
	}
	stationParam := paramsResp.Res[0].Stations[0].StationParam

	req := StationMeteoHistoryRequest{
		StationID:    stationID,
		ParamID:      paramID,
		StationParam: stationParam,
		FieldID:      fieldID,
		DateStart:    "2026-05-01T00:00:00Z",
		DateFinish:   "2026-05-20T23:59:59Z",
		AvgPeriod:    "day",
		AvgType:      "avg",
	}

	resp, err := client.StationMeteoHistoryGet(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Res, "expected live station meteo history response to be present")
	for i, item := range resp.Res {
		if i >= 3 {
			break
		}
		t.Logf("history %d date=%s value=%v", i+1, item.Date, item.Value)
	}
	logResponsePartial(t, "stationmeteohistoryget response", resp.Res)
}

func TestClient_LiveStationGeoDataGet(t *testing.T) {
	client := getLiveClient(t)

	stationIDStr := os.Getenv("WEATHER_API_STATION_ID")
	if stationIDStr == "" {
		t.Skip("skip live station geo data test unless WEATHER_API_STATION_ID is set")
	}

	stationID, err := strconv.Atoi(stationIDStr)
	assert.NoError(t, err)

	req := StationGeoDataRequest{
		StationID:  stationID,
		DateStart:  "2026-05-01T00:00:00Z",
		DateFinish: "2026-05-01T02:00:00Z",
	}

	resp, err := client.StationGeoDataGet(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Res, "expected live station geo data response to be present")
	for i, item := range resp.Res {
		if i >= 3 {
			break
		}
		t.Logf("geodata %d time=%s lat=%v lon=%v", i+1, item.PointTime, item.Lat, item.Lon)
	}
	logResponsePartial(t, "stationgeodataget response", resp.Res)
}

func TestClient_LiveStationLastDataGet(t *testing.T) {
	client := getLiveClient(t)

	stationIDStr := os.Getenv("WEATHER_API_STATION_ID")
	if stationIDStr == "" {
		t.Skip("skip live station data test unless WEATHER_API_STATION_ID is set")
	}

	stationID, err := strconv.Atoi(stationIDStr)
	assert.NoError(t, err)

	resp, err := client.StationLastDataGet(context.Background(), stationID)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Res, "expected live station data response to contain results")
	for _, item := range resp.Res {
		t.Logf("station param %d %s = %v %s", item.ParamID, item.ParamName, item.Value, item.Unit)
	}
}

func TestClient_LiveFieldLastDataGet(t *testing.T) {
	client := getLiveClient(t)

	fieldIDStr := os.Getenv("WEATHER_API_FIELD_ID")
	if fieldIDStr == "" {
		t.Skip("skip live field data test unless WEATHER_API_FIELD_ID is set")
	}

	fieldID, err := strconv.Atoi(fieldIDStr)
	assert.NoError(t, err)

	resp, err := client.FieldLastDataGet(context.Background(), fieldID)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Res, "expected live field data response to contain results")
	for _, item := range resp.Res {
		assert.NotEmpty(t, item.StationValue, "expected station values in live field response")
		t.Logf("field param %d %s -> %d station values", item.ParamID, item.ParamName, len(item.StationValue))
	}
}
