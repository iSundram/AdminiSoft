
package auth

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"

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
