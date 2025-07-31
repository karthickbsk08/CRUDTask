package constants

const (
	SuccessCode           = "S"
	ErrorCode             = "S"
	DeletionSuccessQuote  = "Task successfully deleted"
	MethodsWithBody       = "POST,PUT"
	RedisServerAddress    = "localhost:30000"
	TaskHandler           = "task:send_reminder"
	UniqueTaskIndentifier = "send:email"
)

var (
	Offset_Auto_increment = 0 // Used to track the offset for auto-incrementing IDs
)
