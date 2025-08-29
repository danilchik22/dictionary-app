package auth

type User struct {
	Id       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"size:255;unique" json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Sex      bool   `json:"sex"`
}

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshBody struct {
	RefreshToken string `json:"refresh_token"`
}

type JSONResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type NewUserResponse struct {
	Message string `json:"message"`
	UserId  int    `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
