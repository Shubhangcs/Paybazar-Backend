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

type AdminWalletTopupResponse struct {
	AdminId       string `json:"admin_id"`
	TransactionId string `json:"transaction_id"`
	Balance       string `json:"balance"`
}

// Master Distributor Wallet Models

type MasterDistributorWalletTopupRequest struct {
	AdminId             string `json:"admin_id"`
	MasterDistributorId string `json:"master_distributor_id"`
	Amount              string `json:"amount"`
	Remarks             string `json:"remarks"`
}

type MasterDistributorWalletTopupResponse struct {
	MasterDistributorId string `json:"master_distributor_id"`
	TransactionId       string `json:"transaction_id"`
	Balance             string `json:"balance"`
}

// Distributor Wallet Models

type DistributorWalletTopupRequest struct {
	AdminId       string `json:"admin_id"`
	DistributorId string `json:"distributor_id"`
	Amount        string `json:"amount"`
	Remarks       string `json:"remarks"`
}

type DistributorWalletTopupResponse struct {
	DistributorId string `json:"distributor_id"`
	TransactionId string `json:"transaction_id"`
	Balance       string `json:"balance"`
}

// User Wallet Models

type UserWalletTopupRequest struct {
	AdminId string `json:"admin_id"`
	UserId  string `json:"user_id"`
	Amount  string `json:"amount"`
	Remarks string `json:"remarks"`
}

type UserWalletTopupResponse struct {
	UserId        string `json:"user_id"`
	TransactionId string `json:"transaction_id"`
	Balance       string `json:"balance"`
}
