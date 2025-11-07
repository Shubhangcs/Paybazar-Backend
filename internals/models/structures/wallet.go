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

type AdminAddAmountRequest struct {
	AdminID            string `json:"admin_id"`
	Amount             string `json:"amount"`
	TransactionType    string `json:"transaction_type"`
	TransactionService string `json:"transaction_service"`
	RefrenceID         string `json:"reference_id"`
	Remarks            string `json:"remarks"`
}

type AdminAddAmountToUserWalletRequest struct {
	AdminID string `json:"admin_id"`
	UserID  string `json:"user_id"`
	Admount string `json:"amount"`
	AdminTransactionType string `json:""`
}
