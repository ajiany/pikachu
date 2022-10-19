package feishu

// 消息类型
const (
	MsgTypeText = "text" // 文本消息
	MsgTypePost = "post" // 富文本消息
)

// 消息接收者ID类型
const (
	RecvTypeOpenId  = "open_id"
	RecvTypeUserId  = "user_id"
	RecvTypeUnionId = "union_id"
	RecvTypeEmail   = "email"
	RecvTypeChatId  = "chat_id"
)

// 用户 ID 类型
const (
	TypeOpenId  = "open_id"
	TypeUserId  = "user_id"
	TypeUnionId = "union_id"
)

// 日程参与人类型
const (
	AttendeeTypeUser       = "user"        // 用户
	AttendeeTypeChat       = "chat"        // 群组
	AttendeeTypeResource   = "resource"    // 会议室
	AttendeeTypeThirdParty = "third_party" // 邮箱
)

// 审批状态
const (
	ApprovalPending  = "PENDING"  // 审批中
	ApprovalApproved = "APPROVED" // 通过
	ApprovalRejected = "REJECTED" // 拒绝
	ApprovalCanceled = "CANCELED" // 撤回
	ApprovalDeleted  = "DELETED"  // 删除
)
