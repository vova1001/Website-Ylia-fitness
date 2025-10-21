package otherfunc

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
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
	from := "yliafitness_helper@mail.ru"
	pass := "aWfFfGRklLhggzbyfwfu"
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

	switch {
	case strings.HasSuffix(toEmail, "@mail.ru"):
		smtpHost = "smtp.mail.ru"
	case strings.HasSuffix(toEmail, "@yandex.ru"):
		smtpHost = "smtp.yandex.ru"
	case strings.HasSuffix(toEmail, "@rambler.ru"):
		smtpHost = "smtp.rambler.ru"
	}

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

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		smtp.PlainAuth("", from, pass, smtpHost),
		from,
		[]string{toEmail},
		[]byte(msg),
	)
}
