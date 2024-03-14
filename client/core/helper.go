package core

import "fmt"

func PrintHelper() {
	helpMsg := `help: print help message.
exit: exit client, shutdown application.
to_${userName}:your message: send message to your friend
userlist: get user list who register in server'
reject_${userName}: do not receive this user's message anymore.
recover_%{userName}: recover the connect with this friend.'`
	fmt.Println(helpMsg)
}
