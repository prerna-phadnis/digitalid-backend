package handlers

type RegisterRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Passport string `json:"passport"`
	Email    string `json:"email"`
	DOB      string `json:"dob"`
	Phone    int    `json:"phone"`
}
