// smsc_whatsapp.go
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
)

type SMSCWAService struct {
	Login    string // SMSC_KZ_LOGIN
	Password string // SMSC_KZ_PASSWORD (пароль API, не от кабинета, если включен отдельный)
	SenderWA string // Номер WA-бота в формате 7XXXXXXXXXX без +
	BaseURL  string // опц.: "https://smsc.kz/sys/send.php"
	HTTP     *http.Client
}

// Формат JSON-ответа smsc.kz при fmt=3
type smscResponse struct {
	ID        *int64  `json:"id,omitempty"`
	Cnt       *int    `json:"cnt,omitempty"`
	Cost      *string `json:"cost,omitempty"`
	Error     *string `json:"error,omitempty"`
	ErrorCode *int    `json:"error_code,omitempty"`
}

func NewWhatsAppSMSC(login, password, senderWA string) *WhatsAppService {
	return &WhatsAppService{
		SMSC: &SMSCWAService{
			Login:    login,
			Password: password,
			SenderWA: senderWA, // формат 7XXXXXXXXXX без +
			BaseURL:  "https://smsc.kz/sys/send.php",
			HTTP:     nil, // по умолчанию клиент со своим таймаутом внутри
		},
	}
}

// Нормализуем номер получателя к "7XXXXXXXXXX"
func normalizeKZPhone(p string) (string, error) {
	// оставить только цифры
	var b strings.Builder
	for _, r := range p {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	digits := b.String()
	if digits == "" {
		return "", errors.New("empty phone")
	}

	// варианты: 8XXXXXXXXXX -> 7XXXXXXXXXX, +7XXXXXXXXXX -> 7XXXXXXXXXX, 77XXXXXXXXX ок
	if strings.HasPrefix(digits, "8") && len(digits) == 11 {
		digits = "7" + digits[1:]
	}
	if strings.HasPrefix(digits, "7") && len(digits) == 11 {
		return digits, nil
	}
	if strings.HasPrefix(digits, "007") && len(digits) == 13 {
		return digits[2:], nil
	}
	return "", fmt.Errorf("unexpected phone format: %s", digits)
}

// Отправка обычного WA-текста через smsc.kz
func (s *SMSCWAService) sendWhatsAppText(to, message string) (*smscResponse, error) {
	if s.Login == "" || s.Password == "" || s.SenderWA == "" {
		return nil, errors.New("smsc.kz credentials not configured")
	}
	base := s.BaseURL
	if base == "" {
		base = "https://smsc.kz/sys/send.php"
	}

	toNorm, err := normalizeKZPhone(to)
	if err != nil {
		return nil, err
	}

	// ВАЖНО: для WhatsApp sender должен быть вида "wa:<номер_бота>"
	params := url.Values{
		"login":  {s.Login},
		"psw":    {s.Password},
		"phones": {toNorm},
		"mes":    {message},
		"sender": {"wa:" + s.SenderWA},
		"fmt":    {"3"}, // JSON
		// "charset": {"utf-8"}, // по желанию
		// "test": {"1"},        // тестовый режим без фактической отправки
	}

	reqURL := base + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	client := s.HTTP
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr smscResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("decode smsc response: %w", err)
	}
	if sr.Error != nil {
		return &sr, fmt.Errorf("smsc error %v: %v", valueOr(sr.ErrorCode, -1), *sr.Error)
	}
	return &sr, nil
}

func valueOr[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// -------------------------------------------
// Адаптация под ваш текущий интерфейс сервиса
// -------------------------------------------

type WhatsAppService struct {
	// Старые поля можно оставить, но использовать новые
	// Token, PhoneNumberID string // больше не нужны
	SMSC *SMSCWAService
}

func (s *WhatsAppService) SendVerificationCode(phone string) (int, error) {
	if s.SMSC == nil {
		return 0, fmt.Errorf("smsc.kz service not configured")
	}
	code, err := generate4DigitCode()
	if err != nil {
		return 0, fmt.Errorf("failed to generate code: %w", err)
	}
	message := fmt.Sprintf("Ваш код подтверждения: %d. Компания Ozge Contract.", code)

	if _, err := s.SMSC.sendWhatsAppText(phone, message); err != nil {
		return 0, err
	}
	return code, nil
}
