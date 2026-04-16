package APIs

import (
	"time"

	"github.com/google/uuid"
)

type API struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
type Create struct {
	Title string `json:"title"`
}
