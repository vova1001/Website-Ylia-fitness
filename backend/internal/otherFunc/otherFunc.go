package otherfunc

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
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
			Timeout: 30 * time.Second,
		},
	}
}

func CreatePayment(yc *m.YookassaClient, amount float64, description string) (*m.YookassaPaymentResponse, error) {
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
		Confirmation: struct {
			Type      string `json:"type"`
			ReturnURL string `json:"return_url"`
		}{
			Type:      "redirect",
			ReturnURL: "https://website-ylia-fitness-frontend.onrender.com/",
		},
	}

	return SendRequest(yc, req)
}

func SendRequest(yc *m.YookassaClient, req *m.YookassaPaymentRequest) (*m.YookassaPaymentResponse, error) {

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(yc.ShopID + ":" + yc.ApiKey))

	httpReq, err := http.NewRequest("POST", yc.BaseURL+"/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotence-Key", fmt.Sprintf("%d", time.Now().UnixNano()))

	httpReq.Header.Set("User-Agent", "GoClient/1.0")

	log.Println("Sending request to Yookassa...")

	resp, err := yc.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	log.Println("Yookassa response code:", resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	log.Println("Yookassa raw response:", string(body))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("yookassa returned error %d: %s", resp.StatusCode, body)
	}

	var paymentResp m.YookassaPaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, fmt.Errorf("json decode error: %w", err)
	}

	return &paymentResp, nil
}
