package hetznerrobot

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HetznerRobotClient struct {
	username string
	password string
	url      string
}

func NewHetznerRobotClient(username string, password string, url string) HetznerRobotClient {
	return HetznerRobotClient{
		username: username,
		password: password,
		url:      url,
	}
}

func (c *HetznerRobotClient) makeAPICall(ctx context.Context, method string, uri string, body io.Reader) ([]byte, error) {
	r, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	r.SetBasicAuth(c.username, c.password)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Hetzner response status %d\n%s", response.StatusCode, bytes)
	if response.StatusCode > 400 {
		return nil, fmt.Errorf("Hetzner API response HTTP %d: %s", response.StatusCode, bytes)
	}

	return bytes, nil
}
