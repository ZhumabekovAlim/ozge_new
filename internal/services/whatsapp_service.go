package services

import (
	"bytes"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"
)

type WhatsAppService struct {
	Token         string
	PhoneNumberID string
}

// NewWhatsAppService returns a new WhatsAppService.
func NewWhatsAppService(token, phoneNumberID string) *WhatsAppService {
	return &WhatsAppService{Token: token, PhoneNumberID: phoneNumberID}
}

// generate4DigitCode returns a cryptographically-strong 4-digit code (1000..9999).
func generate4DigitCode() (int, error) {
	n, err := crand.Int(crand.Reader, big.NewInt(9000))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 1000, nil
}

// SendVerificationCode generates a 4-digit code and sends it via WhatsApp Cloud API as a text message.
func (s *WhatsAppService) SendVerificationCode(phone string) (int, error) {
	if s.Token == "" || s.PhoneNumberID == "" {
		return 0, fmt.Errorf("whatsapp credentials not configured")
	}

	code, err := generate4DigitCode()
	if err != nil {
		return 0, fmt.Errorf("failed to generate code: %w", err)
	}

	message := fmt.Sprintf("Ваш код подтверждения: %d. Компания Ozge Contract.", code)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                phone,
		"type":              "text",
		"text": map[string]string{
			"body": message,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	endpoint := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", s.PhoneNumberID)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return 0, fmt.Errorf("whatsapp status: %s", resp.Status)
	}

	var result struct {
		Error interface{} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if result.Error != nil {
			return 0, fmt.Errorf("whatsapp error: %v", result.Error)
		}
	}

	return code, nil
}
