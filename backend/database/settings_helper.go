package database

import (
	"crypto/tls"
	"fmt"
	"hy-board-backend/models"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)

// GetSetting gets a setting value from the database.
func GetSetting(key string) string {
	if DB == nil {
		return ""
	}
	var setting models.SystemSetting
	if err := DB.Where("`key` = ?", key).First(&setting).Error; err == nil {
		return setting.Value
	}
	return ""
}

// GetSettingBool gets a setting as a boolean.
func GetSettingBool(key string) bool {
	val := GetSetting(key)
	return strings.ToLower(val) == "true" || val == "1"
}

// GetSettingInt gets a setting as an integer.
func GetSettingInt(key string) int {
	val := GetSetting(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

// UpdateSetting updates or creates a setting.
func UpdateSetting(key string, value string) error {
	if DB == nil {
		return nil
	}
	var setting models.SystemSetting
	err := DB.Where("`key` = ?", key).First(&setting).Error
	if err == nil {
		setting.Value = value
		return DB.Save(&setting).Error
	}
	// Create new
	setting = models.SystemSetting{
		Key:   key,
		Value: value,
	}
	return DB.Create(&setting).Error
}

// SendEmail sends an email using the SMTP settings configured in the system.
func SendEmail(to, subject, body string) error {
	host := GetSetting("smtp_host")
	portStr := GetSetting("smtp_port")
	username := GetSetting("smtp_username")
	password := GetSetting("smtp_password")
	from := GetSetting("smtp_from")
	encryption := GetSetting("smtp_encryption") // SSL, STARTTLS, none

	if host == "" || portStr == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	// Format message
	fromHeader := mail.Address{Name: GetSetting("site_name"), Address: from}
	toHeader := mail.Address{Address: to}

	header := make(map[string]string)
	header["From"] = fromHeader.String()
	header["To"] = toHeader.String()
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", username, password, host)

	if encryption == "SSL" || port == 465 {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("failed to dial TLS: %v", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %v", err)
		}
		defer client.Close()

		if username != "" && password != "" {
			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("failed to authenticate SMTP: %v", err)
			}
		}

		if err = client.Mail(from); err != nil {
			return fmt.Errorf("failed to set sender: %v", err)
		}

		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to add recipient: %v", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %v", err)
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write message body: %v", err)
		}

		err = w.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %v", err)
		}

		return client.Quit()
	} else {
		// STARTTLS or none
		client, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %v", err)
		}
		defer client.Close()

		if encryption == "STARTTLS" || port == 587 {
			tlsconfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host,
			}
			if err := client.StartTLS(tlsconfig); err != nil {
				return fmt.Errorf("failed to start TLS: %v", err)
			}
		}

		if username != "" && password != "" {
			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("failed to authenticate SMTP: %v", err)
			}
		}

		if err = client.Mail(from); err != nil {
			return fmt.Errorf("failed to set sender: %v", err)
		}

		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to add recipient: %v", err)
		}

		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %v", err)
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write message body: %v", err)
		}

		err = w.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %v", err)
		}

		return client.Quit()
	}
}
