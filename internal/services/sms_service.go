package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SMSService struct {
	APIKey string
}

// NewSMSService returns a new SMSService.
func NewSMSService(apiKey string) *SMSService {
	return &SMSService{APIKey: apiKey}
}

// SendVerificationCode генерирует 4-значный код и отправляет его через Mobizon.
func (s *SMSService) SendVerificationCode(phone string) (int, error) {
	if s.APIKey == "" {
		return 0, fmt.Errorf("mobizon api key not configured")
	}

	code, err := generate4DigitCode()
	if err != nil {
		return 0, fmt.Errorf("failed to generate code: %w", err)
	}

	message := fmt.Sprintf("Ваш код подтверждения: %d. Компания Ozge Contract.", code)

	// Корректный endpoint (KZ)
	endpoint := "https://api.mobizon.kz/service/Message/SendSmsMessage"

	values := url.Values{}
	values.Set("apiKey", s.APIKey)
	values.Set("recipient", phone)
	values.Set("text", message)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(endpoint, values)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("mobizon status: %s", resp.Status)
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	// Если Mobizon вернул JSON — прочтём и проверим код
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if result.Code != 0 {
			return 0, fmt.Errorf("mobizon error: %s", result.Message)
		}
	}

	return code, nil
}
