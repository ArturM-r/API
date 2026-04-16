package user

import (
	"time"

	"github.com/google/uuid"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	User  *Response `json:"user"`
	Token string    `json:"token"`
}

type Response struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
