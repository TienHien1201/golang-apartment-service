package cloudinary

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"thomas.vn/apartment_service/internal/config"
)

func NewCloudinary(cfg config.CloudinaryConfig) (*cloudinary.Cloudinary, error) {
	if cfg.CloudName == "" || cfg.APIKey == "" || cfg.APISecret == "" {
		return nil, fmt.Errorf("cloudinary config is missing")
	}

	return cloudinary.NewFromParams(
		cfg.CloudName,
		cfg.APIKey,
		cfg.APISecret,
	)
}
