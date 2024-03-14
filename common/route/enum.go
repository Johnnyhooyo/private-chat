package route

type Type string

const (
	LogIn    Type = "login"
	LogOut   Type = "logout"
	Chat     Type = "chat"
	UserList Type = "userList"
	Reject   Type = "reject"
	Recover  Type = "recover"
)
