package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		return "", fmt.Errorf("invalid request body: %w", err)
	}

	// create your transaction id, etc.
	transactionID, err := pr.query.PayoutRequestInitilizationRequest(&req)
	if err != nil {
		return "", fmt.Errorf("failed to initilize payout request: %w", err)
	}

	// ---- RKIT payoutTransfer call ----
	url := "https://v2bapi.rchargekit.biz/rkitpayout/payoutTransfer"

	// Token from env (or inject via config)
	token := os.Getenv("RKIT_API_TOKEN")
	if token == "" {
		return "", fmt.Errorf("missing RKIT_API_TOKEN")
	}

	// Build request payload as RKIT expects
	// Map your incoming req fields to RKIT's names
	payload := map[string]interface{}{
		"mobile_no":          req.MobileNumber,    // string, 10-digit
		"account_no":         req.AccountNumber,   // string
		"ifsc":               req.IFSCCode,        // string
		"bank_name":          req.BankName,        // string (<= 20 chars)
		"beneficiary_name":   req.BeneficiaryName, // string
		"amount":             req.Amount,          // float or string; RKIT doc says float
		"transfer_type":      req.TransferType,    // "5" for IMPS, "6" for NEFT
		"partner_request_id": transactionID,       // use your generated id
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("encode payout payload: %w", err)
	}

	// Context with timeout (good practice for external calls)
	ctx, cancel := context.WithTimeout(e.Request().Context(), 15*time.Second)
	defer cancel()

	payoutReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create payout request: %w", err)
	}
	payoutReq.Header.Set("Content-Type", "application/json")
	payoutReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 20 * time.Second, // extra guard; ctx still rules
	}

	resp, err := client.Do(payoutReq)
	if err != nil {
		return "", fmt.Errorf("send payout request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read payout response: %w", err)
	}

	// You can optionally parse it to act on status (1=success IMPS, 2=pending NEFT)
	// per RKIT docs. :contentReference[oaicite:1]{index=1}
	// type rkitResp struct {
	// 	Error        int    `json:"error"`
	// 	Msg          string `json:"msg"`
	// 	Status       int    `json:"status"`     // 1=success (IMPS), 2=pending (NEFT)
	// 	OrderID      string `json:"orderid"`
	// 	OpTransID    string `json:"optransid"`
	// 	PartnerReqID string `json:"partnerreqid"`
	// }

	// Return raw response JSON (your signature requires string)
	return string(respBytes), nil
}
