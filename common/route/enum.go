package route

type Type string

const (
	LogIn     Type = "login"
	LogOut    Type = "logout"
	Chat      Type = "chat"
	UserList  Type = "userlist"
	Reject    Type = "reject"
	Recover   Type = "recover"
	Broadcast Type = "broadcast"
	Heartbeat Type = "heartbeat"
)
