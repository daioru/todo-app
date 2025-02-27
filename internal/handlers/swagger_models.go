package handlers

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
