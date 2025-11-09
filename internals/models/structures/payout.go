package structures

type PayoutRequest struct {
	TransactionID         string `json:"transaction_id"`
	OperatorTransactionID string `json:"operator_transaction_id"`
	OrderID               string `json:"order_id"`
	UserID                string `json:"user_id"`
	MobileNumber          string `json:"mobile_number"`
	AccountNumber         string `json:"acount_number"`
	IFSCCode              string `json:"ifsc_code"`
	BankName              string `json:"bank_name"`
	BeneficiaryName        string `json:"benificary_name"`
	Amount                string `json:"amount"`
	TransferType          string `json:"transfer_type"`
	TransactionStatus     string `json:"transaction_status"`
	Remarks               string `json:"remarks"`
}

type PayoutInitilizationRequest struct {
	UserID            string `json:"user_id"`
	MobileNumber      string `json:"mobile_number"`
	AccountNumber     string `json:"account_number"`
	IFSCCode          string `json:"ifsc_code"`
	BankName          string `json:"bank_name"`
	BeneficiaryName    string `json:"benificary_name"`
	Amount            string `json:"amount"`
	TransferType      string `json:"transfer_type"`
	TransactionStatus string `json:"transaction_status"`
	Remarks           string `json:"remarks"`
}

type PayoutFinalSuccessRequest struct {
	TransactionID         string `json:"transaction_id"`
	OperatorTransactionID string `json:"operator_transaction_id"`
	OrderID               string `json:"order_id"`
	TransactionStatus     string `json:"transaction_status"`
}
