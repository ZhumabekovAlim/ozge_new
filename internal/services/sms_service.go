package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SMSService provides methods to send SMS messages via Mobizon API.
type SMSService struct {
	APIKey string
}

// NewSMSService creates a new SMSService with the provided API key.
func NewSMSService(apiKey string) *SMSService {
	return &SMSService{APIKey: apiKey}
}

// SendSMS sends an SMS message using Mobizon service.
func (s *SMSService) SendSMS(phone, text string) error {
	if s.APIKey == "" {
		return fmt.Errorf("mobizon api key not configured")
	}
	endpoint := "https://api.mobizon.kz/service/message/sendsmsmessage"
	values := url.Values{}
	values.Set("apiKey", s.APIKey)
	values.Set("recipient", phone)
	values.Set("text", text)

	resp, err := http.PostForm(endpoint, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mobizon status: %s", resp.Status)
	}
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if result.Code != 0 {
			return fmt.Errorf("mobizon error: %s", result.Message)
		}
	}
	return nil
}
