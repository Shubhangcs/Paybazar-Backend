package structures

type Ticket struct {
	AdminID string `json:"admin_id"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	Message string `json:"message"`
}
