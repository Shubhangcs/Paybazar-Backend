package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/labstack/echo/v4"
)

type payoutRepo struct {
	query *queries.Query
}

func NewPayoutRepository(query *queries.Query) *payoutRepo {
	return &payoutRepo{
		query: query,
	}
}

func (pr *payoutRepo) PayoutRequest(e echo.Context) (string, error) {
	var req structures.PayoutInitilizationRequest
	if err := e.Bind(&req); err != nil {
		return "", fmt.Errorf("invalid request format: %w", err)
	}
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("invalid request data: %w", err)
	}

	amt, err := strconv.ParseFloat(req.Amount , 64)
	if err != nil {
		return "" , fmt.Errorf("failed to parse amount: %w", err)
	}

	if amt < 1000 {
		return "" , fmt.Errorf("failed to execuite minimum transaction is 1000")
	}

	// Check User Balance
	hasBalance, err := pr.query.CheckUserBalance(req.UserID, req.Amount)
	if err != nil {
		return "", fmt.Errorf("failed to check user wallet balance: %w", err)
	}
	if !hasBalance {
		return "", fmt.Errorf("insufficient balance")
	}

	// Check Payout Limit
	hasExceded, err := pr.query.CheckPayoutLimit(req.UserID, req.Amount)
	if err != nil {
		return "", fmt.Errorf("failed to check payout limit: %w", err)
	}
	if !hasExceded {
		return "", fmt.Errorf("payout limit exceded")
	}

	// Initilize Payout Request
	apiReqBody, err := pr.query.InitilizePayoutRequest(&structures.PayoutInitilizationRequest{
		UserID:          req.UserID,
		MobileNumber:    req.MobileNumber,
		AccountNumber:   req.AccountNumber,
		IFSCCode:        req.IFSCCode,
		BankName:        req.BankName,
		BeneficiaryName: req.BeneficiaryName,
		Amount:          req.Amount,
		TransferType:    req.TransferType,
		Remarks:         req.Remarks,
		Commission:      req.Commission,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initilize payout request: %w", err)
	}

	// Api Request
	if err := e.Validate(req); err != nil {
		return "", fmt.Errorf("failed to validate api request format: %w", err)
	}
	token := os.Getenv("RKIT_API_TOKEN")
	if token == "" {
		return "", fmt.Errorf("missing RKIT_API_TOKEN")
	}
	var url string = "https://v2bapi.rechargkit.biz/rkitpayout/payoutTransfer"
	reqBody, err := json.Marshal(apiReqBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode api request json: %w", err)
	}

	apiRequest, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create payout request: %w", err)
	}
	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 20 * time.Second, // extra guard; ctx still rules
	}
	resp, err := client.Do(apiRequest)
	if err != nil {
		return "", fmt.Errorf("failed to send api request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var base struct {
		Error int `json:"error"`
	}

	if err := json.Unmarshal(respBytes, &base); err != nil {
		return "", fmt.Errorf("failed to unmarshal base response: %w", err)
	}

	if base.Error != 0 {
		var apiFailureResponse structures.PayoutApiFailureResponse
		if err := json.Unmarshal(respBytes, &apiFailureResponse); err != nil {
			return "", fmt.Errorf("failed to unmarshal failure response: %w", err)
		}
		apiFailureResponse.PayoutTransactionID = apiReqBody.PartnerRequestID
		if err := pr.query.PayoutFailure(&apiFailureResponse); err != nil {
			return "", fmt.Errorf("failed to update api failure: %w", err)
		}
		return "", fmt.Errorf("failed to complete api request: %v", string(respBytes))
	}

	var apiSuccessResponse structures.PayoutApiSuccessResponse
	if err := json.Unmarshal(respBytes, &apiSuccessResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal success response: %w", err)
	}

	if err := pr.query.PayoutSuccess(&apiSuccessResponse); err != nil {
		return "", fmt.Errorf("failed to update api success response: %w", err)
	}
	return "transaction successfull", nil
}
