package pkgtotp

import (
	"encoding/base64"
	"fmt"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type GenerateResult struct {
	Secret string `json:"secret"`
	QRCode string `json:"qrCode"`
}

func Generate(email string) (*GenerateResult, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Apartment Business",
		AccountName: email,
	})
	if err != nil {
		return nil, err
	}
	otpURL := key.URL()

	png, err := qrcode.Encode(otpURL, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	base64QR := base64.StdEncoding.EncodeToString(png)

	return &GenerateResult{
		Secret: key.Secret(),
		QRCode: fmt.Sprintf("data:image/png;base64,%s", base64QR),
	}, nil
}

func Verify(token string, secret string) bool {
	return totp.Validate(token, secret)
}
