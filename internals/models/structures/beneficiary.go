package structures

type BeneficiaryModel struct {
	BeneficiaryID       string `json:"beneficiary_id"`
	MobileNumber        string `json:"mobile_number"`
	BankName            string `json:"bank_name"`
	IFSCCode            string `json:"ifsc_code"`
	AccountNumber       string `json:"account_number"`
	BeneficiaryName     string `json:"beneficiary_name"`
	BeneficiaryPhone    string `json:"beneficiary_phone"`
	BeneficiaryVerified bool   `json:"beneficiary_verified"`
}
