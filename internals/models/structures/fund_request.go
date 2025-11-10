package structures

type CreateFundRequestModel struct {
	AdminID           string `json:"admin_id" validate:"required,uuid4"`
	RequesterID       string `json:"requester_id" validate:"required,uuid4"`
	RequesterUniqueID string `json:"requester_unique_id" validate:"required"`
	RequesterName     string `json:"requester_name" validate:"required"`
	RequesterType     string `json:"requster_type" validate:"required"`
	Amount            string `json:"amount" validate:"required"`
	BankName          string `json:"bank_name" validate:"required"`
	AccountNumber     string `json:"account_number" validate:"required"`
	IFSCCode          string `json:"ifsc_code" validate:"required"`
	BankBranch        string `json:"bank_branch" validate:"required"`
	UTRNumber         string `json:"utr_number" validate:"required"`
	Remarks           string `json:"remarks" validate:"required"`
}

type GetFundRequestModel struct {
	RequestId         string `json:"request_id" validate:"required,uuid4"`
	RequestUniqueId   string `json:"request_unique_id" validate:"required"`
	RequesterId       string `json:"requester_id" validate:"required,uuid4"`
	RequesterUniqueId string `json:"requester_unique_id" validate:"required"`
	RequesterName     string `json:"requester_name" validate:"required"`
	RequesterType     string `json:"requester_type" validate:"required"`
	Amount            string `json:"amount" validate:"required"`
	BankName          string `json:"bank_name" validate:"required"`
	AccountNumber     string `json:"account_number" validate:"required"`
	IFSCCode          string `json:"ifsc_code" validate:"required"`
	BankBranch        string `json:"bank_branch" validate:"required"`
	UTRNumber         string `json:"utr_number" validate:"required"`
	Remarks           string `json:"remarks" validate:"required"`
	RequestStatus     string `json:"request_status" validate:"required"`
}

type AcceptFundRequestModel struct {
	AdminID   string `json:"admin_id" validate:"required,uuid4"`
	RequestID string `json:"request_id" validate:"required,uuid4"`
}

type FundRequestResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
}
