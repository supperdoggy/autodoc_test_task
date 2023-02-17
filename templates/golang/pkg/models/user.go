package models

type User struct {
	ID    string `json:"id,omitempty" bson:"_id"`
	Email string `json:"email,omitempty" bson:"email"`

	Name      string `json:"name,omitempty" bson:"name"`
	IsDeleted bool   `json:"is_deleted,omitempty" bson:"is_deleted"`
}
