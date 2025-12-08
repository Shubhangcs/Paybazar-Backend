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

type WalletTransaction struct {
	TransactionID      string  `json:"transaction_id"`
	TransactorID       string  `json:"transactor_id"`
	ReceiverID         string  `json:"receiver_id"`
	TransactorName     string  `json:"transactor_name"`
	ReceiverName       string  `json:"receiver_name"`
	TransactorType     string  `json:"transactor_type"`
	ReceiverType       string  `json:"receiver_type"`
	TransactionType    string  `json:"transaction_type"`
	TransactionService string  `json:"transaction_service"`
	Amount             float64 `json:"amount"`
	TransactionStatus  string  `json:"transaction_status"`
	Remarks            string  `json:"remarks"`
}

type WalletResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
}

type RefundRequest struct {
	AdminID     string `json:"admin_id"`
	PhoneNumber string `json:"phone_number"`
	Amount      string `json:"amount"`
}

type MasterDistributorRefundRetailerRequest struct {
	MasterDistributorID string `json:"master_distributor_id"`
	PhoneNumber         string `json:"phone_number"`
	Amount              string `json:"amount"`
}

type DistributorRefundRetailerRequest struct {
	DistributorID string `json:"distributor_id"`
	PhoneNumber   string `json:"phone_number"`
	Amount        string `json:"amount"`
}
