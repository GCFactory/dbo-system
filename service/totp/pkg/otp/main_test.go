package otp

import (
	"encoding/base32"
	"github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type tc struct {
	TS     int64
	TOTP   string
	Mode   config.Algorithm
	Secret string
}

var (
	secSha1   = base32.StdEncoding.EncodeToString([]byte("12345678901234567890"))
	secSha256 = base32.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012"))
	secSha512 = base32.StdEncoding.EncodeToString([]byte("1234567890123456789012345678901234567890123456789012345678901234"))

	rfcMatrixTCs = []tc{
		{59, "94287082", config.AlgorithmSHA1, secSha1},
		{59, "46119246", config.AlgorithmSHA256, secSha256},
		{59, "90693936", config.AlgorithmSHA512, secSha512},
		{1111111109, "07081804", config.AlgorithmSHA1, secSha1},
		{1111111109, "68084774", config.AlgorithmSHA256, secSha256},
		{1111111109, "25091201", config.AlgorithmSHA512, secSha512},
		{1111111111, "14050471", config.AlgorithmSHA1, secSha1},
		{1111111111, "67062674", config.AlgorithmSHA256, secSha256},
		{1111111111, "99943326", config.AlgorithmSHA512, secSha512},
		{1234567890, "89005924", config.AlgorithmSHA1, secSha1},
		{1234567890, "91819424", config.AlgorithmSHA256, secSha256},
		{1234567890, "93441116", config.AlgorithmSHA512, secSha512},
		{2000000000, "69279037", config.AlgorithmSHA1, secSha1},
		{2000000000, "90698825", config.AlgorithmSHA256, secSha256},
		{2000000000, "38618901", config.AlgorithmSHA512, secSha512},
		{20000000000, "65353130", config.AlgorithmSHA1, secSha1},
		{20000000000, "77737706", config.AlgorithmSHA256, secSha256},
		{20000000000, "47863826", config.AlgorithmSHA512, secSha512},
	}
)

// Test vectors from http://tools.ietf.org/html/rfc6238#appendix-B
// NOTE -- the test vectors are documented as having the SAME
// secret -- this is WRONG -- they have a variable secret
// depending upon the hmac algorithm:
//
//	http://www.rfc-editor.org/errata_search.php?rfc=6238
//
// this only took a few hours of head/desk interaction to figure out.
func TestValidateRFCMatrix(t *testing.T) {
	for _, tx := range rfcMatrixTCs {
		valid, err := ValidateCustom(tx.TOTP, tx.Secret, time.Unix(tx.TS, 0).UTC(),
			config.ValidateOpts{
				Digits:    config.DigitsEight,
				Algorithm: tx.Mode,
			})
		require.NoError(t, err,
			"unexpected error totp=%s mode=%v ts=%v", tx.TOTP, tx.Mode, tx.TS)
		require.True(t, valid,
			"unexpected totp failure totp=%s mode=%v ts=%v", tx.TOTP, tx.Mode, tx.TS)
	}
}

func TestGenerateRFCTCs(t *testing.T) {
	for _, tx := range rfcMatrixTCs {
		passcode, err := GenerateCodeCustom(tx.Secret, time.Unix(tx.TS, 0).UTC(),
			config.ValidateOpts{
				Digits:    config.DigitsEight,
				Algorithm: tx.Mode,
			})
		assert.Nil(t, err)
		assert.Equal(t, tx.TOTP, passcode)
	}
}

func TestValidateSkew(t *testing.T) {
	secSha1 := base32.StdEncoding.EncodeToString([]byte("12345678901234567890"))

	tests := []tc{
		{29, "94287082", config.AlgorithmSHA1, secSha1},
		{59, "94287082", config.AlgorithmSHA1, secSha1},
		{61, "94287082", config.AlgorithmSHA1, secSha1},
	}

	for _, tx := range tests {
		valid, err := ValidateCustom(tx.TOTP, tx.Secret, time.Unix(tx.TS, 0).UTC(),
			config.ValidateOpts{
				Digits:    config.DigitsEight,
				Algorithm: tx.Mode,
				Skew:      1,
			})
		require.NoError(t, err,
			"unexpected error totp=%s mode=%v ts=%v", tx.TOTP, tx.Mode, tx.TS)
		require.True(t, valid,
			"unexpected totp failure totp=%s mode=%v ts=%v", tx.TOTP, tx.Mode, tx.TS)
	}
}

func TestGenerate(t *testing.T) {
	secret, url, err := Generate(config.GenerateOpts{
		Issuer:      "SnakeOil",
		AccountName: "alice@example.com",
	})

	require.NoError(t, err, "generate basic TOTP")
	require.Equal(t, 32, len(*secret), "Secret is 32 bytes long as base32.")

	secret, url, err = Generate(config.GenerateOpts{
		Issuer:      "Snake Oil",
		AccountName: "alice@example.com",
	})
	require.NoError(t, err, "issuer with a space in the name")
	require.Contains(t, *url, "issuer=Snake%20Oil")

	secret, url, err = Generate(config.GenerateOpts{
		Issuer:      "SnakeOil",
		AccountName: "alice@example.com",
		SecretSize:  20,
	})
	require.NoError(t, err, "generate larger TOTP")
	require.Equal(t, 32, len(*secret), "Secret is 32 bytes long as base32.")

	secret, url, err = Generate(config.GenerateOpts{
		Issuer:      "SnakeOil",
		AccountName: "alice@example.com",
		SecretSize:  13, // anything that is not divisible by 5, really
	})
	require.NoError(t, err, "Secret size is valid when length not divisible by 5.")
	require.NotContains(t, *secret, "=", "Secret has no escaped characters.")

	secret, url, err = Generate(config.GenerateOpts{
		Issuer:      "SnakeOil",
		AccountName: "alice@example.com",
		Secret:      []byte("helloworld"),
	})
	require.NoError(t, err, "Secret generation failed")
	sec, err := B32NoPadding.DecodeString(*secret)
	require.NoError(t, err, "Secret wa not valid base32")
	require.Equal(t, sec, []byte("helloworld"), "Specified Secret was not kept")
}

//
//func TestGoogleLowerCaseSecret(t *testing.T) {
//	w, err := config.NewKeyFromURL(`otpauth://totp/Google%3Afoo%40example.com?secret=qlt6vmy6svfx4bt4rpmisaiyol6hihca&issuer=Google`)
//	require.NoError(t, err)
//	sec := w.Secret()
//	require.Equal(t, "qlt6vmy6svfx4bt4rpmisaiyol6hihca", sec)
//
//	n := time.Now().UTC()
//	code, err := GenerateCode(w.Secret(), n)
//	require.NoError(t, err)
//
//	valid := Validate(code, w.Secret())
//	require.True(t, valid)
//}
