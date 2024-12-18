//go:generate easyjson -all .
package dto

//easyjson:json
type Signup struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
