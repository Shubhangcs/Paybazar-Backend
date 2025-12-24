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
	BeneficiaryName       string `json:"benificary_name"`
	Amount                string `json:"amount"`
	TransferType          string `json:"transfer_type"`
	TransactionStatus     string `json:"transaction_status"`
	Remarks               string `json:"remarks"`
	Commission            string `json:"commission"`
}

type PayoutInitilizationRequest struct {
	UserID          string `json:"user_id"`
	MobileNumber    string `json:"mobile_number"`
	AccountNumber   string `json:"account_number"`
	IFSCCode        string `json:"ifsc_code"`
	BankName        string `json:"bank_name"`
	BeneficiaryName string `json:"beneficiary_name"`
	Amount          string `json:"amount"`
	TransferType    string `json:"transfer_type"`
	Remarks         string `json:"remarks"`
	Commission      string `json:"commission"`
	MPIN            string `json:"mpin"`
}

type PayoutApiRequest struct {
	PartnerRequestID string `json:"partner_request_id"`
	MobileNumber     string `json:"mobile_no"`
	AccountNumber    string `json:"account_no"`
	IFSCCode         string `json:"ifsc"`
	BankName         string `json:"bank_name"`
	BeneficiaryName  string `json:"beneficiary_name"`
	Amount           string `json:"amount"`
	TransferType     string `json:"transfer_type"`
}

type PayoutFinal struct {
	Status                int    `json:"status"`
	OpertaorTransactionID string `json:"optransid"`
	OrderID               string `json:"orderid"`
	PartnerRequestID      string `json:"partnerreqid"`
}

type PayoutApiSuccessResponse struct {
	Error                 int    `json:"error"`
	Message               string `json:"msg"`
	Status                int    `json:"status"`
	OrderID               string `json:"orderid"`
	OperatorTransactionID string `json:"optransid"`
	PartnerRequestID      string `json:"partnerreqid"`
}

type GetPayoutLogs struct {
	PayoutTransactionID    string `json:"payout_transaction_id"`
	TransactionID          string `json:"transaction_id"`
	PhoneNumber            string `json:"phone_number"`
	BankName               string `json:"bank_name"`
	BeneficiaryName        string `json:"beneficiary_name"`
	Amount                 string `json:"amount"`
	Commission             string `json:"commission"`
	TransferType           string `json:"transfer_type"`
	TransactionStatus      string `json:"transaction_status"`
	TransactionDateAndTime string `json:"transaction_date_and_time"`
	AccountNumber          string `json:"account_number"`
	UserID                 string `json:"user_id"`
}

type PayoutApiFailureResponse struct {
	Error               int    `json:"error"`
	Message             string `json:"message"`
	PayoutTransactionID string `json:"payout_transaction_id"`
}

type PayoutVerifyAccountResponse struct {
	StatusCode int               `json:"status_code"`
	Status     bool              `json:"status"`
	Message    string            `json:"message"`
	Data       PayoutUserDetails `json:"data"`
}

type PayoutUserDetails struct {
	UserName   string `json:"c_name"`
	BankName   string `json:"bank_name"`
	BranchName string `json:"branch_name"`
}

type PayoutRefund struct {
	TransactionID string `json:"transaction_id"`
}
