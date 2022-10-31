package user

import (
	"fmt"
	osuser "os/user"
)

func CurrentUser() *osuser.User {
	currentUser, err := osuser.Current()

	if err != nil {
		panic(err)
	}

	return currentUser
}

func UserUidAndGuidString() string {
	currentUser := CurrentUser()

	return fmt.Sprintf("%s:%s", currentUser.Uid, currentUser.Gid)
}
