
package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strconv"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret(userID uint) (string, string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", "", err
	}

	secretBase32 := base32.StdEncoding.EncodeToString(secret)
	
	// Generate QR code URL
	qrCodeURL := fmt.Sprintf("otpauth://totp/AdminiSoftware:%d?secret=%s&issuer=AdminiSoftware",
		userID, secretBase32)

	return secretBase32, qrCodeURL, nil
}

func ValidateTOTP(token, secret string) bool {
	return totp.Validate(token, secret)
}

func GenerateBackupCodes() ([]string, error) {
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code := make([]byte, 4)
		_, err := rand.Read(code)
		if err != nil {
			return nil, err
		}
		
		codes[i] = fmt.Sprintf("%08d", int(code[0])<<24|int(code[1])<<16|int(code[2])<<8|int(code[3]))
	}
	return codes, nil
}
package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"github.com/pquerna/otp/totp"
)

type TwoFactorAuth struct{}

func NewTwoFactorAuth() *TwoFactorAuth {
	return &TwoFactorAuth{}
}

func (tfa *TwoFactorAuth) GenerateSecret(email string) (string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate secret: %w", err)
	}
	
	return base32.StdEncoding.EncodeToString(secret), nil
}

func (tfa *TwoFactorAuth) GenerateQRCode(email, secret string) (string, error) {
	return totp.QRCodeURL(email, "AdminiSoftware", secret)
}

func (tfa *TwoFactorAuth) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

func (tfa *TwoFactorAuth) GenerateBackupCodes() []string {
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code := make([]byte, 8)
		rand.Read(code)
		codes[i] = fmt.Sprintf("%x", code)
	}
	return codes
}
package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TwoFactorManager struct{}

func NewTwoFactorManager() *TwoFactorManager {
	return &TwoFactorManager{}
}

func (tf *TwoFactorManager) GenerateSecret(username string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "AdminiSoftware",
		AccountName: username,
		SecretSize:  32,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func (tf *TwoFactorManager) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

func (tf *TwoFactorManager) GenerateBackupCodes() ([]string, error) {
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code, err := tf.generateRandomCode()
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return codes, nil
}

func (tf *TwoFactorManager) generateRandomCode() (string, error) {
	bytes := make([]byte, 5)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(bytes)[:8], nil
}

func (tf *TwoFactorManager) GenerateQRCode(secret, username string) (string, error) {
	key, err := otp.NewKeyFromURL(fmt.Sprintf("otpauth://totp/AdminiSoftware:%s?secret=%s&issuer=AdminiSoftware", username, secret))
	if err != nil {
		return "", err
	}
	return key.URL(), nil
}
