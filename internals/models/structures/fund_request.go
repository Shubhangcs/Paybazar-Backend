package structures

type FundRequest struct {
	AdminId       string `json:"admin_id" validate:"required,uuid4"`
	RequestId     string `json:"request_id"`
	RequesterName string `json:"requester_name" validate:"required"`
	RequesterId   string `json:"requester_id" validate:"required,uuid4"`
	RequesterType string `json:"requester_type" validate:"required"`
	Amount        string `json:"amount" validate:"required"`
	BankName      string `json:"bank_name" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	IFSCCode      string `json:"ifsc_code" validate:"required"`
	BankBranch    string `json:"bank_branch" validate:"required"`
	UTRNumber     string `json:"utr_number" validate:"required"`
	Remarks       string `json:"remarks" validate:"required"`
	RequestStatus string `json:"request_status"`
}

type AcceptFundRequest struct {
	AdminId   string `json:"admin_id" validate:"required,uuid4"`
	RequestId string `json:"request_id" validate:"required,uuid4"`
}

type FundRequestResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
}
