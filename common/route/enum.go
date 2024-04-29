package route

type Type string

const (
	LogIn     Type = "login"
	LogOut    Type = "logout"
	Chat      Type = "chat"
	AddGroup  Type = "addGroup"
	QuitGroup Type = "quitGroup"
	UserList  Type = "userlist"
	Reject    Type = "reject"
	Recover   Type = "recover"
	Broadcast Type = "broadcast"
	Heartbeat Type = "heartbeat"
)
