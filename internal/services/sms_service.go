package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type SMSService struct {
	APIKey string
}

func NewSMSService(apiKey string) *SMSService {
	return &SMSService{APIKey: apiKey}
}

// SendVerificationCode генерирует 4-значный код и отправляет его через Mobizon
func (s *SMSService) SendVerificationCode(phone string) (int, error) {
	if s.APIKey == "" {
		return 0, fmt.Errorf("mobizon api key not configured")
	}

	// Генерация 4-значного кода
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9000) + 1000

	message := fmt.Sprintf("Ваш код подтверждения: %d. Компания Ozge Contract.", code)

	endpoint := "https://api.mobizon.kz/service/message/sendsmsmessage"
	values := url.Values{}
	values.Set("apiKey", s.APIKey)
	values.Set("recipient", phone)
	values.Set("text", message)

	resp, err := http.PostForm(endpoint, values)
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
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if result.Code != 0 {
			return 0, fmt.Errorf("mobizon error: %s", result.Message)
		}
	}

	return code, nil
}
