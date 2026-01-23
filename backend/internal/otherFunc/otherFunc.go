package otherfunc

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	m "github.com/vova1001/Website-Ylia-fitness/internal/model"
)

func EmailCheck(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GeneratorToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("err gener rand byte")
	}
	return hex.EncodeToString(b), nil
}

func SendResetEmail(toEmail, resetLink string) error {
	// Загрузка переменных окружения
	emailBotDefault := os.Getenv("EMAIL_BOT")
	emailPassDefault := os.Getenv("EMAIL_BOT_PASS")
	emailBotYandex := os.Getenv("EMAIL_BOT_YANDEX")
	emailPassYandex := os.Getenv("EMAIL_PASS_YANDEX")
	smtpHostDefault := os.Getenv("SMTP_HOST")
	smtpPortDefault := os.Getenv("SMTP_PORT")

	// КРИТИЧЕСКАЯ ОТЛАДКА: что реально загрузилось
	log.Println("========== EMAIL DEBUG START ==========")
	log.Printf("[DEBUG] Recipient email: %s", toEmail)
	log.Printf("[DEBUG] Loaded env variables:")
	log.Printf("[DEBUG]   EMAIL_BOT: '%s' (len pass: %d)",
		emailBotDefault, len(emailPassDefault))
	log.Printf("[DEBUG]   EMAIL_BOT_YANDEX: '%s'", emailBotYandex)
	log.Printf("[DEBUG]   EMAIL_PASS_YANDEX exists: %v (length: %d)",
		emailPassYandex != "", len(emailPassYandex))
	if emailPassYandex != "" {
		// Показываем только первые 2 символа пароля для безопасности
		log.Printf("[DEBUG]   EMAIL_PASS_YANDEX starts with: '%s...'",
			emailPassYandex[:min(2, len(emailPassYandex))])
	}
	log.Printf("[DEBUG]   SMTP_HOST: '%s'", smtpHostDefault)
	log.Printf("[DEBUG]   SMTP_PORT: '%s'", smtpPortDefault)

	var from, pass, smtpHost, smtpPort string

	log.Println("[EMAIL] Step 1: checking recipient domain")

	// Выбираем настройки в зависимости от домена получателя
	switch {
	case strings.HasSuffix(toEmail, "@yandex.ru"):
		// Для Яндекс - используем специальные настройки
		smtpHost = "smtp.yandex.ru"
		smtpPort = "587" // Порт 587 для Яндекс (с STARTTLS)
		from = emailBotYandex
		pass = emailPassYandex

		log.Println("[EMAIL] Using YANDEX credentials for recipient:", toEmail)
		log.Printf("[DEBUG YANDEX] Using: host=%s, port=%s", smtpHost, smtpPort)
		log.Printf("[DEBUG YANDEX] From address: '%s'", from)
		log.Printf("[DEBUG YANDEX] Password provided: %v (len: %d)",
			pass != "", len(pass))

		// Проверка обязательных полей
		if from == "" {
			log.Println("[ERROR] Yandex from email is EMPTY!")
			return fmt.Errorf("Yandex sender email not configured")
		}
		if pass == "" {
			log.Println("[ERROR] Yandex password is EMPTY!")
			return fmt.Errorf("Yandex password not configured")
		}

	case strings.HasSuffix(toEmail, "@mail.ru"):
		// Для Mail.ru
		smtpHost = "smtp.mail.ru"
		smtpPort = "587"
		from = emailBotDefault
		pass = emailPassDefault
		log.Println("[EMAIL] Using MAIL.RU credentials")

	case strings.HasSuffix(toEmail, "@rambler.ru"):
		// Для Rambler
		smtpHost = "smtp.rambler.ru"
		smtpPort = "587"
		from = emailBotDefault
		pass = emailPassDefault
		log.Println("[EMAIL] Using RAMBLER credentials")

	default:
		smtpHost = smtpHostDefault
		smtpPort = smtpPortDefault
		from = emailBotDefault
		pass = emailPassDefault
		log.Println("[EMAIL] Using DEFAULT credentials")
	}

	log.Println("[EMAIL] Step 2: smtpHost =", smtpHost, "from =", from)
	log.Printf("[DEBUG FINAL] Using: host=%s, port=%s, from=%s, pass_len=%d",
		smtpHost, smtpPort, from, len(pass))

	htmlBody := fmt.Sprintf(`
    <html>
    <body style="font-family:Arial,sans-serif; text-align:center;">
        <h2>Сброс пароля</h2>
        <p>Нажмите кнопку ниже, чтобы установить новый пароль:</p>
        <a href="%s"
           style="display:inline-block;
                  padding:10px 20px;
                  background-color:#4CAF50;
                  color:white;
                  text-decoration:none;
                  border-radius:5px;">
           Изменить пароль
        </a>
        <p>Если вы не запрашивали сброс, проигнорируйте это письмо.</p>
    </body>
    </html>`, resetLink)

	msg := "From: " + from + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: Сброс пароля\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n\n" +
		htmlBody

	log.Println("[EMAIL] Step 3: calling SendMail...")
	log.Printf("[DEBUG] SMTP server: %s:%s", smtpHost, smtpPort)

	// Временная отладка: логируем всё перед отправкой
	log.Println("========== EMAIL DEBUG END ==========")

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		smtp.PlainAuth("", from, pass, smtpHost),
		from,
		[]string{toEmail},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("[EMAIL ERROR] SendMail failed: %v", err)
	} else {
		log.Println("[EMAIL SUCCESS] Email sent successfully")
	}

	return err
}

// Вспомогательная функция
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewYookassaClient(shopID, apiKey string) *m.YookassaClient {
	return &m.YookassaClient{
		ShopID:  shopID,
		ApiKey:  apiKey,
		BaseURL: "https://api.yookassa.ru/v3",
		Client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func CreatePayment(yc *m.YookassaClient, amount float64, description string, metadata map[string]string) (*m.YookassaPaymentResponse, error) {
	req := &m.YookassaPaymentRequest{
		Amount: struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		}{
			Value:    fmt.Sprintf("%.2f", amount),
			Currency: "RUB",
		},
		Capture:     true,
		Description: description,
		Metadata:    metadata,
		Confirmation: struct {
			Type      string `json:"type"`
			ReturnURL string `json:"return_url"`
		}{
			Type:      "redirect",
			ReturnURL: "https://website-ylia-fitness-frontend.onrender.com/",
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", yc.BaseURL+"/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	httpReq.SetBasicAuth(yc.ShopID, yc.ApiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotence-Key", uuid.New().String())
	httpReq.Header.Set("User-Agent", "GoClient/1.0")

	resp, err := yc.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("yookassa returned error %d: %s", resp.StatusCode, string(body))
	}

	var paymentResp m.YookassaPaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	if paymentResp.Confirmation.ConfirmationURL == "" {
		return nil, fmt.Errorf("confirmation_url пустой, raw response: %s", string(body))
	}

	return &paymentResp, nil
}

func DifferenceDead(d time.Duration) m.ResponseDuration {
	if d < 0 {
		return m.ResponseDuration{Days: 0, Hours: 0, Text: "Course is Dead"}
	}
	AllHours := int(d.Hours())
	return m.ResponseDuration{Days: AllHours / 24, Hours: AllHours % 24, Text: "There is time"}
}
