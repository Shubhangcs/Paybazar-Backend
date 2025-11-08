package structures

type FundRequest struct {
	AdminId       string `json:"admin_id"`
	RequesterId   string `json:"requester_id"`
	RequesterType string `json:"requester_type"`
	Amount        string `json:"amount"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	IFSCCode      string `json:"ifsc_code"`
	BankBranch    string `json:"bank_branch"`
	UTRNumber     string `json:"utr_number"`
	Remarks       string `json:"remarks"`
	RequestStatus string `json:"request_status"`
}

type FundRequestResponse struct {
	RequestId     string `json:"request_id"`
	RequestStatus string `json:"request_status"`
}
