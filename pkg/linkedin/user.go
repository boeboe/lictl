package linkedin

import "time"

// User represents the structure of a LinkedIn user.
type User struct {
	ConnectionCount int    `json:"connectionCount"`
	FollowerCount   int    `json:"followerCount"`
	UserTitle       string `json:"userTitle"`
	Location        string `json:"location"`
	Name            string `json:"name"`
	UserLink        string `json:"userLink"`
}

func SearchUsersOnline(keywords []string, interval time.Duration, debug bool) ([]User, error) {
	return make([]User, 0), nil
}