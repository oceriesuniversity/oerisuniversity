package services

import (
	"fmt"
	
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets environment variable as integer or returns default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// sendSuccessEmail sends application success email
func SendSuccessEmail(toEmail, userName string) error {
	// Get SMTP configuration from environment variables
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnvAsInt("SMTP_PORT", 587)
	smtpUsername := getEnv("SMTP_USERNAME", "")
	smtpPassword := getEnv("SMTP_PASSWORD", "")
	fromEmail := getEnv("FROM_EMAIL", "ocerisumps@gmail.com")
	fromName := getEnv("FROM_NAME", "Oceries University Colledge Of Metaphisical Sciences")

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromEmail, fromName))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Application Successful! 🎉")
	
	// HTML email content
	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Application Successful</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); 
					color: white; padding: 20px; text-align: center; border-radius: 10px 10px 0 0; }
			.content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
			.footer { text-align: center; margin-top: 20px; color: #666; font-size: 14px; }
			.button { display: inline-block; padding: 12px 24px; background: #667eea; 
					color: white; text-decoration: none; border-radius: 5px; margin: 20px 0; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Congratulations, %s! 🎊</h1>
			</div>
			<div class="content">
				<h2>Your Application Has Been Approved!</h2>
				<p>We're excited to inform you that your application has been successfully processed and approved.</p>
				<p>You can now access all the features and benefits of our platform.</p>
				
				

				<p><strong>What's next?</strong></p>
				<ul>
					<li>Complete your profile setup</li>
					<li>Explore our features</li>
					<li>Start using our services immediately</li>
				</ul>

				<p>If you have any questions or need assistance, don't hesitate to contact our support team.</p>
			</div>
			<div class="footer">
				<p>Best regards,<br>The %s Team</p>
				<p><small>This is an automated message, please do not reply to this email.</small></p>
			</div>
		</div>
	</body>
	</html>
	`, userName, fromName)

	m.SetBody("text/html", htmlBody)

	// Create dialer and send email
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email to %s: %v", toEmail, err)
		fmt.Println("\n Email service error!!!!")
	}

	fmt.Printf("✅ Email successfully sent to %s", toEmail)
	return nil
}






