package user

import (
	"fmt"
	osuser "os/user"
)

func UserUidAndGuidString() string {
	currentUser, err := osuser.Current()

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s:%s", currentUser.Uid, currentUser.Gid)
}
