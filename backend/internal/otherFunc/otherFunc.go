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
	from := os.Getenv("EMAIL_BOT")
	pass := os.Getenv("EMAIL_BOT_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	log.Println("[EMAIL] Step 1: starting function")

	switch {
	case strings.HasSuffix(toEmail, "@mail.ru"):
		smtpHost = "smtp.mail.ru"
	case strings.HasSuffix(toEmail, "@yandex.ru"):
		smtpHost = "smtp.yandex.ru"
	case strings.HasSuffix(toEmail, "@rambler.ru"):
		smtpHost = "smtp.rambler.ru"
	}

	log.Println("[EMAIL] Step 2: smtpHost =", smtpHost)

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

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		smtp.PlainAuth("", from, pass, smtpHost),
		from,
		[]string{toEmail},
		[]byte(msg),
	)
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

<<<<<<< Updated upstream
func DifferenceDead(d time.Duration) m.ResponseDuration {
	if d < 0 {
		return m.ResponseDuration{Days: 0, Hours: 0, Text: "Course is Dead"}
	}
	AllHours := int(d.Hours())
	return m.ResponseDuration{Days: AllHours / 24, Hours: AllHours % 24, Text: "There is time"}
=======
func FormatDuration(time time.Duration) {
	AllHours := time.Hours()
	Days := AllHours / 24

>>>>>>> Stashed changes
}
