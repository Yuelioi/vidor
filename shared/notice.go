package shared

// APP通知
type Notice interface {
	ProgressUpdate(part Part)
}

type NoticeData struct {
	EventName string
	Message   interface{}
}

type Callback func(data NoticeData)

type eventName struct {
	Start string
	Stop  string
	Error string
}

var EventName = eventName{}
