package feishu

import (
	"context"
	"time"

	"github.com/ajiany/pikachu/tools/httpx"
	"github.com/pkg/errors"
)

// CreateCalendarEvent 创建日程
func (cli *Client) GetMeetingRoomFreeBusy(ctx context.Context, roomIds []string, timeMin, timeMax time.Time) (map[string][]FreeBusyItem, error) {
	query := map[string]interface{}{
		"room_ids": roomIds,
		"time_min": timeMin.Format(time.RFC3339),
		"time_max": timeMax.Format(time.RFC3339),
	}
	r, err := cli.Api.GET(ctx, "/meeting_room/freebusy/batch_get", &httpx.Option{Query: query})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			FreeBusy map[string][]FreeBusyItem `json:"free_busy"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, errors.WithStack(err)
	}

	return resp.Data.FreeBusy, nil
}

// 获取建筑物列表
func (cli *Client) GetBuildings(ctx context.Context, pageToken string) ([]Building, *Pagination, error) {
	query := map[string]interface{}{
		"page_token": pageToken,
		"page_size":  100,
	}
	r, err := cli.Api.GET(ctx, "/meeting_room/building/list", &httpx.Option{Query: query})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Pagination
			Buildings []Building `json:"buildings"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return resp.Data.Buildings, &resp.Data.Pagination, nil
}

// 获取全部建筑物
func (cli *Client) AllBuildings(ctx context.Context) ([]Building, error) {
	var pageToken string
	buildings := make([]Building, 0)

	for {
		items, pag, err := cli.GetBuildings(ctx, pageToken)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		buildings = append(buildings, items...)

		if len(items) == 0 || !pag.HasMore {
			break
		}
	}

	return buildings, nil
}

// 获取建筑会议室列表
func (cli *Client) GetMeetingRooms(ctx context.Context, buildingId string, pageToken string) ([]MeetingRoom, *Pagination, error) {
	query := map[string]interface{}{
		"building_id": buildingId,
		"page_token":  pageToken,
		"page_size":   100,
	}
	r, err := cli.Api.GET(ctx, "/meeting_room/room/list", &httpx.Option{Query: query})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	var resp struct {
		Data struct {
			Pagination
			Rooms []MeetingRoom `json:"rooms"`
		} `json:"data"`
	}
	if err := r.ParsedStruct(&resp); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return resp.Data.Rooms, &resp.Data.Pagination, nil
}

// 获取建筑全部会议室
func (cli *Client) AllMeetingRooms(ctx context.Context, buildingId string) ([]MeetingRoom, error) {
	var pageToken string
	rooms := make([]MeetingRoom, 0)

	for {
		items, pag, err := cli.GetMeetingRooms(ctx, buildingId, pageToken)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		rooms = append(rooms, items...)

		if len(items) == 0 || !pag.HasMore {
			break
		}
	}

	return rooms, nil
}
