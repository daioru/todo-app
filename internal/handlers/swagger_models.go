package handlers

type UserData struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateTaskData struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

type UpdateTaskData struct {
	Title       string `json:"title" validate:"optional"`
	Description string `json:"description" validate:"optional"`
	Status      string `json:"status" validate:"optional"`
}
