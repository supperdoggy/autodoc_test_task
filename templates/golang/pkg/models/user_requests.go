package models

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	User  *User  `json:"user"`
	Error string `json:"error"`
}

type GetUserRequest struct {
	ID    *string `json:"id"`
	Email *string `json:"email"`
}

type DeleteUserRequest struct {
	ID string `json:"id"`
}

type ListUserResponse struct {
	Users []User `json:"users"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
