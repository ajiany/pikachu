package feishu

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ResponseData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type TenantAccessToken struct {
	Token  string `json:"tenant_access_token"`
	Expire int64  `json:"expire"`
}

type TextMessage struct {
	Text string `json:"text"`
}

func (t TextMessage) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type Pagination struct {
	HasMore   bool   `json:"has_more"`
	PageToken string `json:"page_token"`
}

type Chat struct {
	Id   string `json:"chat_id"`
	Name string `json:"name"`
}

type User struct {
	UnionId      string       `json:"union_id"`
	UserId       string       `json:"user_id"`
	OpenId       string       `json:"open_id"`
	Name         string       `json:"name"`
	EnName       string       `json:"en_name"`
	Nickname     string       `json:"nickname"`
	Email        string       `json:"email"`
	Gender       int64        `json:"gender"`
	City         string       `json:"city"`
	Country      string       `json:"country"`
	WorkStation  string       `json:"work_station"`
	JoinTime     int64        `json:"join_time"`
	EmployeeNo   string       `json:"employee_no"`
	EmployeeType int64        `json:"employee_type"`
	JobTitle     string       `json:"job_title"`
	CustomAttrs  []CustomAttr `json:"custom_attrs"`
	Status       *UserStatus  `json:"status"`
}

type UserStatus struct {
	IsFrozen    bool `json:"is_frozen"`
	IsResigned  bool `json:"is_resigned"`
	IsActivated bool `json:"is_activated"`
	IsExited    bool `json:"is_exited"`
	IsUnjoin    bool `json:"is_unjoin"`
}

type CustomAttr struct {
	Id    string      `json:"id"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type UserGroup struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserIdItem struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

// ****************** Approval ******************

// 审批实例
type ApprovalInstance struct {
	ApprovalCode string         `json:"approval_code"`
	ApprovalName string         `json:"approval_name"`
	StartTime    int64          `json:"start_time"`
	EndTime      int64          `json:"end_time"`
	UUID         string         `json:"uuid"`
	UserId       string         `json:"user_id"`
	OpenId       string         `json:"open_id"`
	SerialNumber string         `json:"serial_number"`
	DepartmentId string         `json:"department_id"`
	Status       string         `json:"status"`
	Form         string         `json:"form"`
	formData     []ApprovalForm `json:"-"`
}

func (a ApprovalInstance) FormData() []ApprovalForm {
	if len(a.formData) > 0 {
		return a.formData
	}

	if err := json.Unmarshal([]byte(a.Form), &a.formData); err != nil {
		logrus.WithError(err).Error("unmarshal fail")
	}
	return a.formData
}

func (a ApprovalInstance) GetForm(id string, name string) *ApprovalForm {
	for _, item := range a.FormData() {
		if item.Id == id || strings.Contains(item.Name, name) {
			return &item
		}
	}

	return nil
}

// 审批控件
type ApprovalForm struct {
	Id       string      `json:"id"`
	CustomId string      `json:"custom_id,omitempty"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	Ext      interface{} `json:"ext"`
}

func (f ApprovalForm) StrVal() string {
	return fmt.Sprint(f.Value)
}

func (f ApprovalForm) StrVals() []string {
	var vals []string
	v, ok := f.Value.([]interface{})
	if ok {
		for _, item := range v {
			vals = append(vals, fmt.Sprint(item))
		}
	}
	return vals
}

func (f ApprovalForm) TimeVal() time.Time {
	t, _ := time.ParseInLocation(time.RFC3339, fmt.Sprint(f.Value), time.Local)
	return t
}

func (f ApprovalForm) TimeIntervalVal() ApprovalFormTimeInterval {
	var i ApprovalFormTimeInterval
	Convert(f.Value, &i)
	return i
}

type ApprovalFormTimeInterval struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Interval string `json:"interval"`
}

func (a ApprovalFormTimeInterval) StartTime() time.Time {
	t, _ := time.ParseInLocation(time.RFC3339, fmt.Sprint(a.Start), time.Local)
	return t
}

func (a ApprovalFormTimeInterval) EndTime() time.Time {
	t, _ := time.ParseInLocation(time.RFC3339, fmt.Sprint(a.End), time.Local)
	return t
}

// ****************** Calendar ******************
type Calendar struct {
	CalendarId string `json:"calendar_id"`
	Summary    string `json:"summary"`
	Type       string `json:"type"`
}

type CalendarEventTime struct {
	Date      string `json:"date,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
}

type CalendarEventVchat struct {
	VcType      string `json:"vc_type"`
	Description string `json:"description"`
	MeetingUrl  string `json:"meeting_url"`
}

// 日程
type CalendarEvent struct {
	EventId          string              `json:"event_id"`
	Summary          string              `json:"summary"`
	Description      string              `json:"description"`
	NeedNotification bool                `json:"need_notification"`
	StartTime        CalendarEventTime   `json:"start_time"`
	EndTime          CalendarEventTime   `json:"end_time"`
	Vchat            *CalendarEventVchat `json:"vchat,omitempty"`
	Recurrence       string              `json:"recurrence,omitempty"`
	Status           string              `json:"status"`
}

// 日程参与者
type CalendarEventAttendee struct {
	AttendeeId      string `json:"attendee_id"`
	Type            string `json:"type"`
	IsOptional      bool   `json:"is_optional,omitempty"`
	UserId          string `json:"user_id,omitempty"`
	ChatId          string `json:"chat_id,omitempty"`
	RoomId          string `json:"room_id,omitempty"`
	ThirdPartyEmail string `json:"third_party_email,omitempty"`
}

// ****************** Meeting Room ******************
type OrganizerInfo struct {
	Name string `json:"name"`
}

// 闲忙
type FreeBusyItem struct {
	Uid           string         `json:"uid"`
	StartTime     string         `json:"start_time"`
	EndTime       string         `json:"end_time"`
	OriginalTime  int64          `json:"original_time"`
	OrganizerInfo *OrganizerInfo `json:"organizer_info"`
}

// 建筑物
type Building struct {
	BuildingId  string   `json:"building_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Floors      []string `json:"floors"`
	CountryId   string   `json:"country_id"`
	DistrictId  string   `json:"district_id"`
}

// 会议室
type MeetingRoom struct {
	RoomId       string `json:"room_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayId    string `json:"display_id"`
	Capacity     int64  `json:"capacity"`
	IsDisabled   bool   `json:"is_disabled"`
	BuildingId   string `json:"building_id"`
	BuildingName string `json:"building_name"`
	FloorName    string `json:"floor_name"`
}

func (r MeetingRoom) DisplayName() string {
	return fmt.Sprintf("%s-%s(%d) %s", r.FloorName, r.Name, r.Capacity, r.BuildingName)
}

// ****************** Event ******************
type EncryptEventData struct {
	Encrypt string `json:"encrypt"`
}

type EventData struct {
	RawData []byte `json:"-"`
	// v2
	Schema string `json:"schema"` // 事件格式的版本。无此字段的即为1.0
	Header struct {
		EventId    string    `json:"event_id"`    // 事件的唯一标识
		Token      string    `json:"token"`       // 即Verification Token
		CreateTime string    `json:"create_time"` //  事件发送的时间
		EventType  EventType `json:"event_type"`  // 事件类型
		TenantKey  string    `json:"tenant_key"`  // 企业标识
		AppId      string    `json:"app_id"`      // 应用ID

	} `json:"header"`

	// v1
	Timestamp string    `json:"ts"`    // 事件发送的时间，一般近似于事件发生的时间。
	UUID      string    `json:"uuid"`  // 事件的唯一标识
	Token     string    `json:"token"` // 即Verification Token
	Type      EventType `json:"type"`  // event_callback-事件推送，url_verification-url地址验证
	Challenge string    `json:"challenge"`

	Event map[string]interface{} `json:"event"`
}

func (e EventData) IsV2() bool {
	return e.Schema != ""
}

func (e EventData) EventType() EventType {
	if e.IsV2() && e.Header.EventType != "" {
		return e.Header.EventType
	}

	return EventType(e.Type)
}

func (e EventData) VerificationToken() string {
	if e.IsV2() {
		return e.Header.Token
	}

	return e.Token
}

type EventUserId struct {
	UnionId string `json:"union_id"`
	UserId  string `json:"user_id"`
	OpenId  string `json:"open_id"`
}

type EventUser struct {
	Name      string      `json:"name"`
	TenantKey string      `json:"tenant_key"`
	UserId    EventUserId `json:"user_id"`
}

type EventChatUserAddedData struct {
	ChatId            string      `json:"chat_id"`
	OperatorId        EventUserId `json:"operator_id"`
	External          bool        `json:"external"`
	OperatorTenantKey string      `json:"operator_tenant_key"`
	Users             []EventUser `json:"users"`
}

type EventStaffAddedData struct {
	Object User `json:"object"`
}

type EventStaffUpdatedData struct {
	Object    User `json:"object"`
	OldObject User `json:"old_object"`
}

type EventMeetingRoomStatusChangedData struct {
	RoomId   string `json:"room_id"`
	RoomName string `json:"room_name"`
}
