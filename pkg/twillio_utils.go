package pkg

import (
	"fmt"
	"os"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwillioUtils struct {
}

func (*TwillioUtils) SendOTP(phoneNumber string, OTP string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})

	// Define message parameters
	params := &openapi.CreateMessageParams{}
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER")) // Twilio registered number
	params.SetTo("+91"+phoneNumber)                        // Recipient number (include country code)
	params.SetBody(fmt.Sprintf("This is a Verification Message from Paybazaar with OTP: %s", OTP))

	// Send SMS
	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	return nil
}
