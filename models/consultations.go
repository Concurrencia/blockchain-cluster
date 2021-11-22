package models

type Consultation struct {
	ID               int    `json:"id"`
	LoanAmount       string `json:"loanAmount"`
	CreditHistory    string `json:"creditHistory"`
	PropertyAreaNum  string `json:"propertyAreaNum"`
	CantMultas       string `json:"cantMultas"`
	NivelGravedadNum string `json:"nivelGravedadNum"`
	Result           string `json:"result"`
	UserID           string `json:"userId"`
}
