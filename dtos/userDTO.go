package dtos

// UserDTO digunakan untuk mengontrol data yang dikembalikan
type UserDTO struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role"`
}
