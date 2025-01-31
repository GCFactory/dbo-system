//go:generate mockgen -source totp.go -destination ../internal/totp/mock/totp_mock.go -package mock
package pkg

import (
	"github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"net/url"
	"time"
)

type Totp interface {
	ValidateCustom(passcode string, secret string, t time.Time, opts config.ValidateOpts) (bool, error)
	Validate(passcode string, secret string) bool
	GenerateCodeCustom(secret string, t time.Time, opts config.ValidateOpts) (passcode string, err error)
	GenerateCode(secret string, t time.Time) (string, error)
	Generate(opts config.GenerateOpts) (*string, *string, error)
	EncodeQuery(v url.Values) string
}
