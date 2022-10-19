package feishu

import (
	"context"
	"fmt"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

func (cli *Client) GetPrimaryCalendar(ctx context.Context) (*Calendar, error) {
	r, err := cli.Api.POST(ctx, "/calendar/v4/calendars/primary", &httpx.Option{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Calendars []struct {
				Calendar Calendar `json:"calendar"`
				UserId   string   `json:"user_id"`
			} `json:"calendars"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return &resp.Data.Calendars[0].Calendar, nil
}

// CreateCalendarEvent 创建日程
func (cli *Client) CreateCalendarEvent(ctx context.Context, calendarId string, event CalendarEvent) (*CalendarEvent, error) {
	url := fmt.Sprintf("/calendar/v4/calendars/%s/events", calendarId)
	r, err := cli.Api.POST(ctx, url, &httpx.Option{Params: event})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Event CalendarEvent `json:"event"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return &resp.Data.Event, nil
}

// UpdateCalendarEvent 更新日程
func (cli *Client) UpdateCalendarEvent(ctx context.Context, calendarId string, event CalendarEvent) (*CalendarEvent, error) {
	url := fmt.Sprintf("/calendar/v4/calendars/%s/events/%s", calendarId, event.EventId)
	r, err := cli.Api.PATCH(ctx, url, &httpx.Option{Params: event})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Event CalendarEvent `json:"event"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return &resp.Data.Event, nil
}

// DeleteCalendarEvent 删除日程
func (cli *Client) DeleteCalendarEvent(ctx context.Context, calendarId string, eventId string) error {
	url := fmt.Sprintf("/calendar/v4/calendars/%s/events/%s", calendarId, eventId)
	_, err := cli.Api.DELETE(ctx, url, &httpx.Option{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// AddCalendarEventAttendees 添加日程参与者
func (cli *Client) AddCalendarEventAttendees(ctx context.Context, calendarId, eventId string, attendees []CalendarEventAttendee) error {
	url := fmt.Sprintf("/calendar/v4/calendars/%s/events/%s/attendees", calendarId, eventId)
	opts := httpx.Option{
		Query: map[string]interface{}{
			"user_id_type": TypeUserId,
		},
		Params: map[string]interface{}{
			"attendees": attendees,
		},
	}
	_, err := cli.Api.POST(ctx, url, &opts)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// AllCalendarEventAttendees 获取所有日程参与人
func (cli *Client) AllCalendarEventAttendees(ctx context.Context, calendarId, eventId string) ([]CalendarEventAttendee, error) {
	var pageToken string
	records := make([]CalendarEventAttendee, 0)

	for {
		items, pag, err := cli.GetCalendarEventAttendees(ctx, calendarId, eventId, pageToken)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		records = append(records, items...)

		if len(items) == 0 || !pag.HasMore {
			break
		}
	}

	return records, nil
}

// 获取日程参与人列表
func (cli *Client) GetCalendarEventAttendees(ctx context.Context, calendarId, eventId string, pageToken string) ([]CalendarEventAttendee, *Pagination, error) {
	query := map[string]interface{}{
		"user_id_type": TypeUserId,
		"page_token":   pageToken,
		"page_size":    100,
	}
	url := fmt.Sprintf("/calendar/v4/calendars/%s/events/%s/attendees", calendarId, eventId)
	r, err := cli.Api.GET(ctx, url, &httpx.Option{Query: query})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Pagination
			Items []CalendarEventAttendee `json:"items"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return resp.Data.Items, &resp.Data.Pagination, nil
}

// DeleteCalendarEventAttendees 移除日程参与者
func (cli *Client) DeleteCalendarEventAttendees(ctx context.Context, calendarId, eventId string, attendeeIds []string) error {
	url := fmt.Sprintf("/calendar/v4/calendars/%s/events/%s/attendees/batch_delete", calendarId, eventId)
	opts := httpx.Option{
		Params: map[string]interface{}{
			"attendee_ids": attendeeIds,
		},
	}
	_, err := cli.Api.POST(ctx, url, &opts)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
