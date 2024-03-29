package scheduleapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client for uksivt-schedule api.
type Client struct {
	client *http.Client
	apiUrl string
}

func NewClient(client *http.Client, apiUrl string) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{client: client, apiUrl: apiUrl}
}

func (c Client) call(req *http.Request, dest any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(dest)
}

func (c Client) get(ctx context.Context, path string, dest any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiUrl+path, nil)
	if err != nil {
		return err
	}
	return c.call(req, dest)
}

const methodGetGroups = "/college_group"

func (c Client) Groups(ctx context.Context) ([]string, error) {
	var result []string

	if err := c.get(ctx, methodGetGroups, &result); err != nil {
		return nil, err
	}

	return result, nil
}

const methodGetTeachers = "/teacher"

func (c Client) Teachers(ctx context.Context) ([]string, error) {
	var result []string

	if err := c.get(ctx, methodGetTeachers, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c Client) Lessons(ctx context.Context, group string, weekStart time.Time) (map[time.Weekday][]Lesson, error) {
	if !MatchGroup(group) {
		return nil, ErrInvalidGroup
	}

	local := weekStart.Format("2006-01-02")
	path := fmt.Sprintf("%s/%s/from_date/%s", methodGetGroups, group, local)

	response := make(map[time.Weekday][]Lesson)
	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}
