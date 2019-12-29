package user

import "time"

// User represents standard user entity of issue#1.
// bookmarkedPosts map contains the postId mapped to the time it was bookmarked.
type User struct {
	Username        string            `json:"username"`
	Email           string            `json:"email,omitempty"`
	FirstName       string            `json:"firstName,omitempty"`
	MiddleName      string            `json:"middleName,omitempty"`
	LastName        string            `json:"lastName,omitempty"`
	CreationTime    time.Time         `json:"creationTime,omitempty"`
	Bio             string            `json:"bio,omitempty"`
	BookmarkedPosts map[int]time.Time `json:"bookmarkedPosts,omitempty"`
	Password        string            `json:"password,omitempty"`
}
