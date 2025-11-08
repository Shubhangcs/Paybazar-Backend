package structures

type TransactionRequest struct {
	TransactorId   string `json:"transactor_id"`
	TransactorType string `json:"transactor_type"`
}

type TransactionResponse struct {
	TransactionId      string `json:"transaction_id"`
	TransactorId       string `json:"transactor_id"`
	Amount             string `json:"amount"`
	TransactionType    string `json:"transaction_type"`
	TransactionService string `json:"transaction_service"`
	ReferenceId        string `json:"reference_id"`
	Remarks            string `json:"remarks"`
	CreatedAt          string `json:"created_at"`
	TransactorType     string `json:"transactor_type"`
}


