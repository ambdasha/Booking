package testutil
//помогает удобно слать HTTP JSON-запросы
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func DoJSON(method, url, token string, body any) (int, []byte, error) {
	var r io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		r = bytes.NewReader(b)
	}

	//создание HTTP-запроса
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, data, nil
}

//для удобного получения токена из ответа логина
func MustTokenFromLoginResponse(status int, body []byte) (string, error) {
	if status != 200 {
		return "", fmt.Errorf("login failed: status=%d body=%s", status, string(body))
	}
	var v struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return "", err
	}
	if v.AccessToken == "" {
		return "", fmt.Errorf("empty access_token in response: %s", string(body))
	}
	return v.AccessToken, nil
}