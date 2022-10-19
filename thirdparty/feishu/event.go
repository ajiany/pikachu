package feishu

type EventType string

const (
	EventChatUserAdded   EventType = "im.chat.member.user.added_v1"
	EventUrlVerification EventType = "url_verification"
	EventStaffAdded      EventType = "contact.user.created_v3"
	EventStaffUpdated    EventType = "contact.user.updated_v3"

	EventMeetingRoomStatusChanged EventType = "meeting_room.meeting_room.status_changed_v1"
)
