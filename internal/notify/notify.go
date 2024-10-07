package notify

type Notification interface {
	Send(msg Notice) error
}

type Notice struct {
	EventName  string `json:"eventName"`
	Content    string `json:"content"`
	NoticeType string `json:"noticeType"` // 消息类型 ['success', 'info', 'warn', 'error', 'secondary', 'contrast']
	Provider   string `json:"provider"`   // 消息提供者
}
