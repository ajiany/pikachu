package feishu

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ResponseData struct {
	Code int    `test_helper:"code"`
	Msg  string `test_helper:"msg"`
}

type TenantAccessToken struct {
	Token  string `test_helper:"tenant_access_token"`
	Expire int64  `test_helper:"expire"`
}

type TextMessage struct {
	Text string `test_helper:"text"`
}

func (t TextMessage) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

type Pagination struct {
	HasMore   bool   `test_helper:"has_more"`
	PageToken string `test_helper:"page_token"`
}

type Chat struct {
	Id   string `test_helper:"chat_id"`
	Name string `test_helper:"name"`
}

type User struct {
	UnionId      string       `test_helper:"union_id"`
	UserId       string       `test_helper:"user_id"`
	OpenId       string       `test_helper:"open_id"`
	Name         string       `test_helper:"name"`
	EnName       string       `test_helper:"en_name"`
	Nickname     string       `test_helper:"nickname"`
	Email        string       `test_helper:"email"`
	Gender       int64        `test_helper:"gender"`
	City         string       `test_helper:"city"`
	Country      string       `test_helper:"country"`
	WorkStation  string       `test_helper:"work_station"`
	JoinTime     int64        `test_helper:"join_time"`
	EmployeeNo   string       `test_helper:"employee_no"`
	EmployeeType int64        `test_helper:"employee_type"`
	JobTitle     string       `test_helper:"job_title"`
	CustomAttrs  []CustomAttr `test_helper:"custom_attrs"`
	Status       *UserStatus  `test_helper:"status"`
}

type UserStatus struct {
	IsFrozen    bool `test_helper:"is_frozen"`
	IsResigned  bool `test_helper:"is_resigned"`
	IsActivated bool `test_helper:"is_activated"`
	IsExited    bool `test_helper:"is_exited"`
	IsUnjoin    bool `test_helper:"is_unjoin"`
}

type CustomAttr struct {
	Id    string      `test_helper:"id"`
	Type  string      `test_helper:"type"`
	Value interface{} `test_helper:"value"`
}

type UserGroup struct {
	Id   string `test_helper:"id"`
	Name string `test_helper:"name"`
}

type UserIdItem struct {
	UserId string `test_helper:"user_id"`
	Email  string `test_helper:"email"`
	Phone  string `test_helper:"phone"`
}

// ****************** Approval ******************

// 审批实例
type ApprovalInstance struct {
	ApprovalCode string         `test_helper:"approval_code"`
	ApprovalName string         `test_helper:"approval_name"`
	StartTime    int64          `test_helper:"start_time"`
	EndTime      int64          `test_helper:"end_time"`
	UUID         string         `test_helper:"uuid"`
	UserId       string         `test_helper:"user_id"`
	OpenId       string         `test_helper:"open_id"`
	SerialNumber string         `test_helper:"serial_number"`
	DepartmentId string         `test_helper:"department_id"`
	Status       string         `test_helper:"status"`
	Form         string         `test_helper:"form"`
	formData     []ApprovalForm `test_helper:"-"`
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
	Id       string      `test_helper:"id"`
	CustomId string      `test_helper:"custom_id,omitempty"`
	Name     string      `test_helper:"name"`
	Type     string      `test_helper:"type"`
	Value    interface{} `test_helper:"value"`
	Ext      interface{} `test_helper:"ext"`
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
	Start    string `test_helper:"start"`
	End      string `test_helper:"end"`
	Interval string `test_helper:"interval"`
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
	CalendarId string `test_helper:"calendar_id"`
	Summary    string `test_helper:"summary"`
	Type       string `test_helper:"type"`
}

type CalendarEventTime struct {
	Date      string `test_helper:"date,omitempty"`
	Timestamp string `test_helper:"timestamp,omitempty"`
	Timezone  string `test_helper:"timezone,omitempty"`
}

type CalendarEventVchat struct {
	VcType      string `test_helper:"vc_type"`
	Description string `test_helper:"description"`
	MeetingUrl  string `test_helper:"meeting_url"`
}

// 日程
type CalendarEvent struct {
	EventId          string              `test_helper:"event_id"`
	Summary          string              `test_helper:"summary"`
	Description      string              `test_helper:"description"`
	NeedNotification bool                `test_helper:"need_notification"`
	StartTime        CalendarEventTime   `test_helper:"start_time"`
	EndTime          CalendarEventTime   `test_helper:"end_time"`
	Vchat            *CalendarEventVchat `test_helper:"vchat,omitempty"`
	Recurrence       string              `test_helper:"recurrence,omitempty"`
	Status           string              `test_helper:"status"`
}

// 日程参与者
type CalendarEventAttendee struct {
	AttendeeId      string `test_helper:"attendee_id"`
	Type            string `test_helper:"type"`
	IsOptional      bool   `test_helper:"is_optional,omitempty"`
	UserId          string `test_helper:"user_id,omitempty"`
	ChatId          string `test_helper:"chat_id,omitempty"`
	RoomId          string `test_helper:"room_id,omitempty"`
	ThirdPartyEmail string `test_helper:"third_party_email,omitempty"`
}

// ****************** Meeting Room ******************
type OrganizerInfo struct {
	Name string `test_helper:"name"`
}

// 闲忙
type FreeBusyItem struct {
	Uid           string         `test_helper:"uid"`
	StartTime     string         `test_helper:"start_time"`
	EndTime       string         `test_helper:"end_time"`
	OriginalTime  int64          `test_helper:"original_time"`
	OrganizerInfo *OrganizerInfo `test_helper:"organizer_info"`
}

// 建筑物
type Building struct {
	BuildingId  string   `test_helper:"building_id"`
	Name        string   `test_helper:"name"`
	Description string   `test_helper:"description"`
	Floors      []string `test_helper:"floors"`
	CountryId   string   `test_helper:"country_id"`
	DistrictId  string   `test_helper:"district_id"`
}

// 会议室
type MeetingRoom struct {
	RoomId       string `test_helper:"room_id"`
	Name         string `test_helper:"name"`
	Description  string `test_helper:"description"`
	DisplayId    string `test_helper:"display_id"`
	Capacity     int64  `test_helper:"capacity"`
	IsDisabled   bool   `test_helper:"is_disabled"`
	BuildingId   string `test_helper:"building_id"`
	BuildingName string `test_helper:"building_name"`
	FloorName    string `test_helper:"floor_name"`
}

func (r MeetingRoom) DisplayName() string {
	return fmt.Sprintf("%s-%s(%d) %s", r.FloorName, r.Name, r.Capacity, r.BuildingName)
}

// ****************** Event ******************
type EncryptEventData struct {
	Encrypt string `test_helper:"encrypt"`
}

type EventData struct {
	RawData []byte `test_helper:"-"`
	// v2
	Schema string `test_helper:"schema"` // 事件格式的版本。无此字段的即为1.0
	Header struct {
		EventId    string    `test_helper:"event_id"`    // 事件的唯一标识
		Token      string    `test_helper:"token"`       // 即Verification Token
		CreateTime string    `test_helper:"create_time"` //  事件发送的时间
		EventType  EventType `test_helper:"event_type"`  // 事件类型
		TenantKey  string    `test_helper:"tenant_key"`  // 企业标识
		AppId      string    `test_helper:"app_id"`      // 应用ID

	} `test_helper:"header"`

	// v1
	Timestamp string    `test_helper:"ts"`    // 事件发送的时间，一般近似于事件发生的时间。
	UUID      string    `test_helper:"uuid"`  // 事件的唯一标识
	Token     string    `test_helper:"token"` // 即Verification Token
	Type      EventType `test_helper:"type"`  // event_callback-事件推送，url_verification-url地址验证
	Challenge string    `test_helper:"challenge"`

	Event map[string]interface{} `test_helper:"event"`
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
	UnionId string `test_helper:"union_id"`
	UserId  string `test_helper:"user_id"`
	OpenId  string `test_helper:"open_id"`
}

type EventUser struct {
	Name      string      `test_helper:"name"`
	TenantKey string      `test_helper:"tenant_key"`
	UserId    EventUserId `test_helper:"user_id"`
}

type EventChatUserAddedData struct {
	ChatId            string      `test_helper:"chat_id"`
	OperatorId        EventUserId `test_helper:"operator_id"`
	External          bool        `test_helper:"external"`
	OperatorTenantKey string      `test_helper:"operator_tenant_key"`
	Users             []EventUser `test_helper:"users"`
}

type EventStaffAddedData struct {
	Object User `test_helper:"object"`
}

type EventStaffUpdatedData struct {
	Object    User `test_helper:"object"`
	OldObject User `test_helper:"old_object"`
}

type EventMeetingRoomStatusChangedData struct {
	RoomId   string `test_helper:"room_id"`
	RoomName string `test_helper:"room_name"`
}
