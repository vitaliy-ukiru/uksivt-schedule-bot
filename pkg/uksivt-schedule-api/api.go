package scheduleapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

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

var ErrInvalidWeekStart = errors.New("week start is not monday")

func (c Client) Lessons(ctx context.Context, group Group, weekStart time.Time) (WeekOfLessons, error) {
	if weekStart.Weekday() != time.Monday {
		return WeekOfLessons{}, ErrInvalidWeekStart
	}

	local := weekStart.Format("2006-01-02")
	path := fmt.Sprintf("%s/%s/from_date/%s", methodGetGroups, group.String(), local)
	response := make(map[string][]Lesson)
	if err := c.get(ctx, path, &response); err != nil {
		return WeekOfLessons{}, err
	}

	return setToWeek(response)
}
