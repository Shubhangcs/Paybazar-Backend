package structures

type AdminWallet struct {
	WalletID string `json:"wallet_id"`
	AdminID  string `json:"admin_id"`
	Balance  string `json:"balance"`
}

type MasterDistributorWallet struct {
	WalletID string `json:"wallet_id"`
	UserID   string `json:"user_id"`
	Balance  string `json:"balance"`
}

type DistributorWallet struct {
	WalletID      string `json:"wallet_id"`
	DistributorID string `json:"distributor_id"`
	Balance       string `json:"balance"`
}

type UserWallet struct {
	WalletID string `json:"wallet_id"`
	UserID   string `json:"user_id"`
	Balance  string `json:"balance"`
}

// Admin Wallet Models

type AdminWalletTopupRequest struct {
	AdminId string `json:"admin_id"`
	Amount  string `json:"amount"`
	Remarks string `json:"remarks"`
}

type MasterDistributorWalletTransactions struct {
	TransactionId       string `json:"transaction_id"`
	MasterDistributorId string `json:"master_distributor_id"`
	Amount              string `json:"amount"`
	TransactionType     string `json:"transaction_type"`
	TransactionService  string `json:"transaction_service"`
	ReferenceId         string `json:"reference_id"`
	Remarks             string `json:"remarks"`
	CreatedAt           string `json:"created_at"`
}

type AdminWalletTransactions struct {
	TransactionId      string `json:"transaction_id"`
	AdminId            string `json:"admin_id"`
	Amount             string `json:"amount"`
	TransactionType    string `json:"transaction_type"`
	TransactionService string `json:"transaction_service"`
	ReferenceId        string `json:"reference_id"`
	Remarks            string `json:"remarks"`
	CreatedAt          string `json:"created_at"`
}

type DistributorWalletTransactions struct {
	TransactionId      string `json:"transaction_id"`
	DistributorId      string `json:"distributor_id"`
	Amount             string `json:"amount"`
	TransactionType    string `json:"transaction_type"`
	TransactionService string `json:"transaction_service"`
	ReferenceId        string `json:"reference_id"`
	Remarks            string `json:"remarks"`
	CreatedAt          string `json:"created_at"`
}

type UserWalletTransactions struct {
	TransactionId      string `json:"transaction_id"`
	UserId             string `json:"user_id"`
	Amount             string `json:"amount"`
	TransactionType    string `json:"transaction_type"`
	TransactionService string `json:"transaction_service"`
	ReferenceId        string `json:"reference_id"`
	Remarks            string `json:"remarks"`
	CreatedAt          string `json:"created_at"`
}

type WalletResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
}
