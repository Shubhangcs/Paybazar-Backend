package repositories

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Srujankm12/paybazar-api/internals/models/queries"
	"github.com/Srujankm12/paybazar-api/internals/models/structures"
	"github.com/Srujankm12/paybazar-api/pkg"
	"github.com/labstack/echo/v4"
)

type payoutRepo struct {
	query    *queries.Query
	jwtUtils *pkg.JwtUtils
}

func NewPayoutRepository(query *queries.Query, jwtUtils *pkg.JwtUtils) *payoutRepo {
	return &payoutRepo{
		query:    query,
		jwtUtils: jwtUtils,
	}
}

func (pr *payoutRepo) bindAndValidate(e echo.Context, v interface{}) error {
	if err := e.Bind(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request format")
	}
	if err := e.Validate(v); err != nil {
		return echo.NewHTTPError(400, "Invalid request data")
	}
	return nil
}

func (pr *payoutRepo) PayoutRequest(e echo.Context) (string, error) {
	var req structures.PayoutInitilizationRequest
	if err := pr.bindAndValidate(e, &req); err != nil {
		return "", err
	}

	amt, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		log.Println("parse amount error:", err)
		return "", echo.NewHTTPError(400, "Invalid amount")
	}

	if amt < 1000 || amt > 25000 {
		return "", echo.NewHTTPError(400, "Minimum transaction amount is 1000")
	}

	// Check User Balance
	hasBalance, err := pr.query.CheckUserBalance(req.UserID, req.Amount, req.Commission)
	if err != nil {
		log.Println("DB check user balance error:", err)
		return "", echo.NewHTTPError(500, "Failed to verify wallet balance")
	}
	if !hasBalance {
		return "", echo.NewHTTPError(400, "Insufficient balance")
	}

	// Check Payout Limit
	// hasNotExceeded, err := pr.query.CheckPayoutLimit(req.UserID, req.Amount)
	// if err != nil {
	// 	log.Println("DB check payout limit error:", err)
	// 	return "", echo.NewHTTPError(500, "Failed to verify payout limit")
	// }
	// if !hasNotExceeded {
	// 	return "", echo.NewHTTPError(400, "Payout limit exceeded")
	// }

	// Check MPIN
	hasMpin, err := pr.query.CheckMpin(req.UserID, req.MPIN)
	if err != nil {
		log.Println("DB check mpin error:", err)
		return "", echo.NewHTTPError(500, "Wrong MPIN")
	}
	if !hasMpin {
		return "", echo.NewHTTPError(400, "Wrong MPIN")
	}

	// Initialize Payout Request (prepare DB entry / partner request id etc.)
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
		log.Println("DB initialize payout request error:", err)
		return "", echo.NewHTTPError(500, "Failed to initialize payout")
	}

	// Prepare external API request
	token := os.Getenv("RKIT_API_TOKEN")
	if token == "" {
		log.Println("missing RKIT_API_TOKEN")
		return "", echo.NewHTTPError(500, "Payout provider configuration error")
	}

	url := "https://v2bapi.rechargkit.biz/rkitpayout/payoutTransfer"
	reqBody, err := json.Marshal(apiReqBody)
	if err != nil {
		log.Println("marshal api request error:", err)
		return "", echo.NewHTTPError(500, "Failed to prepare payout request")
	}

	apiRequest, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Println("create api request error:", err)
		return "", echo.NewHTTPError(500, "Failed to create payout request")
	}
	apiRequest.Header.Set("Content-Type", "application/json")
	apiRequest.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Do(apiRequest)
	if err != nil {
		log.Println("send api request error:", err)
		return "", echo.NewHTTPError(502, "Failed to contact payout provider")
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read api response error:", err)
		return "", echo.NewHTTPError(502, "Failed to read payout provider response")
	}

	var base struct {
		Error int `json:"error"`
	}

	if err := json.Unmarshal(respBytes, &base); err != nil {
		log.Println("unmarshal base response error:", err, "response:", string(respBytes))
		return "", echo.NewHTTPError(502, "Unexpected response from payout provider")
	}

	if base.Error != 0 {
		var apiFailureResponse structures.PayoutApiFailureResponse
		if err := json.Unmarshal(respBytes, &apiFailureResponse); err != nil {
			// still log the raw response for debugging
			log.Println("unmarshal failure response error:", err, "response:", string(respBytes))
			return "", echo.NewHTTPError(502, "Payout provider returned an error")
		}
		// attach our partner request id to failure record and persist
		apiFailureResponse.PayoutTransactionID = apiReqBody.PartnerRequestID
		if err := pr.query.PayoutFailure(&apiFailureResponse); err != nil {
			log.Println("DB record payout failure error:", err)
			// don't expose DB internals — but report provider error to client
			return "", echo.NewHTTPError(502, "Payout failed")
		}
		return "", echo.NewHTTPError(502, "Payout failed")
	}

	var apiSuccessResponse structures.PayoutApiSuccessResponse
	if err := json.Unmarshal(respBytes, &apiSuccessResponse); err != nil {
		log.Println("unmarshal success response error:", err, "response:", string(respBytes))
		return "", echo.NewHTTPError(502, "Unexpected response from payout provider")
	}

	if err := pr.query.PayoutSuccess(&apiSuccessResponse); err != nil {
		log.Println("DB update payout success error:", err)
		// We successfully reached provider but failed to persist; still inform user provider succeeded.
		return "", echo.NewHTTPError(500, "Payout succeeded but saving status failed")
	}

	return "Transaction successful", nil
}

func (pr *payoutRepo) GetPayoutTransactions(e echo.Context) (*[]structures.GetPayoutLogs, error) {
	var userId = e.Param("user_id")
	res, err := pr.query.GetPayoutTransactions(userId)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to fetch payout transactions")
	}
	return res, nil
}

// helper: generate a unique numeric string of given length using crypto/rand
func generateUniqueNumericString(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}
	max := big.NewInt(10)
	var sb strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteString(n.String())
	}
	return sb.String(), nil
}

func (pr *payoutRepo) VerifyAccountNumber(e echo.Context) (*structures.PayoutVerifyAccountResponse, error) {
	refID := e.Param("phone")
	accNum := e.Param("account_number")
	ifsc := e.Param("ifsc")

	// generate unique random numeric string (nonce)
	nonce, _ := generateUniqueNumericString(12)

	// generate JWT token
	token, err := pr.jwtUtils.GenerateTokenForExternalAPI(nonce)
	if err != nil {
		return nil, err
	}

	// payload with dynamic values
	payload := map[string]string{
		"refid":          refID,
		"account_number": accNum,
		"ifsc_code": ifsc,
	}

	bodyBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST",
		"https://uat.paysprint.in/sprintverify-uat/api/v1/verification/penny_drop_v2",
		strings.NewReader(string(bodyBytes)),
	)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Token", token)
	req.Header.Add("authorisedkey", "TVRJek5EVTJOelUwTnpKRFQxSlFNREF3TURFPQ==")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	// decode API response directly and return
	var resp structures.PayoutVerifyAccountResponse
	_ = json.Unmarshal(body, &resp)

	log.Println(string(body))

	return &resp, nil
}
