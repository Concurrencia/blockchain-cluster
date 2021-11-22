package models

type Consultation struct {
	ID     int    `json:"id"`
	Dato1  string `json:"dato1"`
	Dato2  string `json:"dato2"`
	Result string `json:"result"`
	UserId string `json:"userId"`
}
