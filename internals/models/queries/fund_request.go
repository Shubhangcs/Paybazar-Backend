package queries

import (
	"context"

	"github.com/Srujankm12/paybazar-api/internals/models/structures"
)

func (q *Query) GetFundRequests(requesterId string) (*[]structures.FundRequest, error) {
	var fundRequest structures.FundRequest
	var fundRequests []structures.FundRequest

	query := `
		SELECT admin_id, requester_id, requester_type, amount, bank_name, account_number, ifsc_code, bank_branch, utr_number, remarks, request_status FROM fund_requests WHERE requester_id=$1
	`
	res, err := q.Pool.Query(context.Background(), query, requesterId)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		if err := res.Scan(
			&fundRequest.AdminId,
			&fundRequest.RequesterId,
			&fundRequest.RequesterType,
			&fundRequest.Amount,
			&fundRequest.BankName,
			&fundRequest.IFSCCode,
			&fundRequest.BankBranch,
			&fundRequest.UTRNumber,
			&fundRequest.Remarks,
			&fundRequest.RequestStatus,
		); err != nil {
			return nil, err
		}
		fundRequests = append(fundRequests, fundRequest)
	}

	if res.Err() != nil {
		return nil, res.Err()
	}
	return &fundRequests, nil
}

func (q *Query) RejectFundRequest() error {
	
}
