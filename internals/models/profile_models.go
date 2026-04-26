package models

type Profile struct {
	ID       int64
	UserID   int64
	FullName *string
	Phone    *string
	Email    *string
}

type CreateProfileRequest struct {
	FullName string `validate:"required,min=2"`
	Phone    string `validate:"required,len=10,numeric"`
}

type ProfileResponse struct {
	User       User    `json:"User"`
	Profile    Profile `json:"Profile"`
	ManualUser bool    `json:"ManualUser"`
}
