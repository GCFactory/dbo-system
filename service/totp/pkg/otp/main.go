package otp

import (
	"crypto/rand"
	platformConfig "github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	totpPkg "github.com/GCFactory/dbo-system/service/totp/pkg"
	"github.com/GCFactory/dbo-system/service/totp/pkg/hotp"
	hotpConfig "github.com/GCFactory/dbo-system/service/totp/pkg/hotp/config"
	"github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	totpConfig "github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TotpStruct struct {
	cfg *platformConfig.Config
	log logger.Logger
}

// ValidateCustom validates a TOTP given a user specified time and custom options.
// Most users should use Validate() to provide an interpolatable TOTP experience.
func (totp TotpStruct) ValidateCustom(passcode string, secret string, t time.Time, opts config.ValidateOpts) (bool, error) {
	if opts.Period == 0 {
		opts.Period = 30
	}

	counters := []uint64{}
	counter := int64(math.Floor(float64(t.Unix()) / float64(opts.Period)))

	counters = append(counters, uint64(counter))
	for i := 1; i <= int(opts.Skew); i++ {
		counters = append(counters, uint64(counter+int64(i)))
		counters = append(counters, uint64(counter-int64(i)))
	}

	for _, counter := range counters {
		rv, err := hotp.ValidateCustom(passcode, counter, secret, hotpConfig.ValidateOpts{
			Digits:    opts.Digits,
			Algorithm: opts.Algorithm,
		})

		if err != nil {
			return false, err
		}

		if rv == true {
			return true, nil
		}
	}

	return false, nil
}

// Validate a TOTP using the current time.
// A shortcut for ValidateCustom, Validate uses a configuration
// that is compatible with Google-Authenticator and most clients.
func (totp TotpStruct) Validate(passcode string, secret string) bool {
	rv, _ := totp.ValidateCustom(
		passcode,
		secret,
		time.Now().UTC(),
		config.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    config.DigitsSix,
			Algorithm: config.AlgorithmSHA1,
		},
	)
	return rv
}

// GenerateCodeCustom takes a timepoint and produces a passcode using a
// secret and the provided opts. (Under the hood, this is making an adapted
// call to hotp.GenerateCodeCustom)
func (totp TotpStruct) GenerateCodeCustom(secret string, t time.Time, opts config.ValidateOpts) (passcode string, err error) {
	if opts.Period == 0 {
		opts.Period = 30
	}
	counter := uint64(math.Floor(float64(t.Unix()) / float64(opts.Period)))
	passcode, err = hotp.GenerateCodeCustom(secret, counter, hotpConfig.ValidateOpts{
		Digits:    opts.Digits,
		Algorithm: opts.Algorithm,
	})
	if err != nil {
		return "", err
	}
	return passcode, nil
}

// GenerateCode creates a TOTP token using the current time.
// A shortcut for GenerateCodeCustom, GenerateCode uses a configuration
// that is compatible with Google-Authenticator and most clients.
func (totp TotpStruct) GenerateCode(secret string, t time.Time) (string, error) {
	return totp.GenerateCodeCustom(secret, t, config.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    config.DigitsSix,
		Algorithm: config.AlgorithmSHA1,
	})
}

func (totp TotpStruct) Generate(opts config.GenerateOpts) (*string, *string, error) {
	// url encode the Issuer/AccountName
	if opts.Issuer == "" {
		return nil, nil, totpConfig.ErrGenerateMissingIssuer
	}

	if opts.AccountName == "" {
		return nil, nil, totpConfig.ErrGenerateMissingAccountName
	}

	if opts.Period == 0 {
		opts.Period = config.DefaultPeriod
	}

	if opts.SecretSize == 0 {
		opts.SecretSize = config.DefaultSecretLength
	}

	if opts.Digits == 0 {
		opts.Digits = config.DigitsSix
	}

	if opts.Rand == nil {
		opts.Rand = rand.Reader
	}

	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example

	v := url.Values{}
	if len(opts.Secret) != 0 {
		v.Set("secret", config.B32NoPadding.EncodeToString(opts.Secret))
	} else {
		secret := make([]byte, opts.SecretSize)
		_, err := opts.Rand.Read(secret)
		if err != nil {
			return nil, nil, err
		}
		v.Set("secret", config.B32NoPadding.EncodeToString(secret))
	}

	v.Set("issuer", opts.Issuer)
	v.Set("period", strconv.FormatUint(uint64(opts.Period), 10))
	v.Set("algorithm", opts.Algorithm.String())
	v.Set("digits", opts.Digits.String())

	u := url.URL{
		Scheme:   "otpauth",
		Host:     "totp",
		Path:     "/" + opts.Issuer + ":" + opts.AccountName,
		RawQuery: totp.EncodeQuery(v),
	}
	str1 := u.String()
	str2 := v.Get("secret")

	return &str2, &str1, nil
}

func (totp TotpStruct) EncodeQuery(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.PathEscape(k) // changed from url.QueryEscape
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.PathEscape(v)) // changed from url.QueryEscape
		}
	}
	return buf.String()
}

func NewTOTPStruct(cfg *platformConfig.Config, logger logger.Logger) totpPkg.Totp {
	return &TotpStruct{cfg: cfg, log: logger}
}
