package models

type User struct {
	Email         string         `json:"email"`
	Password      string         `json:"password"`
	Consultations []Consultation `json:"consultations"`
}

type UserResponse struct {
	ID            string         `json:"id"`
	Email         string         `json:"email"`
	Password      string         `json:"password"`
	Consultations []Consultation `json:"consultations"`
}
