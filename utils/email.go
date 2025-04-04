package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendPasswordResetEmail sends a password recovery email
func SendPasswordResetEmail(to string) error {
	smtpHost := "smtp.example.com"
	smtpPort := "587"
	senderEmail := "no-reply@example.com"
	senderPassword := os.Getenv("SMTP_PASSWORD") // Ensure this is stored securely

	subject := "Password Reset Request"
	body := fmt.Sprintf("Hello,\n\nClick the following link to reset your password: https://example.com/reset-password?email=%s\n\nBest regards,\nYour Team", to)

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", senderEmail, to, subject, body)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{to}, []byte(msg))

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
