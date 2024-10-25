package dto

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	AvatarId string `json:"avatar_id"`
	Status   int    `json:"status"`
}
