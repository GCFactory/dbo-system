/**
 *  Copyright 2014 Paul Querna
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package hotp

import (
	"crypto/hmac"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"math"
	"strings"
)

const debug = false

// Validate a HOTP passcode given a counter and secret.
// This is a shortcut for ValidateCustom, with parameters that
// are compataible with Google-Authenticator.
func Validate(passcode string, counter uint64, secret string) bool {
	rv, _ := ValidateCustom(
		passcode,
		counter,
		secret,
		ValidateOpts{
			Digits:    config.DigitsSix,
			Algorithm: config.AlgorithmSHA1,
		},
	)
	return rv
}

// ValidateOpts provides options for ValidateCustom().
type ValidateOpts struct {
	// Digits as part of the input. Defaults to 6.
	Digits config.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm config.Algorithm
}

// GenerateCode creates a HOTP passcode given a counter and secret.
// This is a shortcut for GenerateCodeCustom, with parameters that
// are compataible with Google-Authenticator.
func GenerateCode(secret string, counter uint64) (string, error) {
	return GenerateCodeCustom(secret, counter, ValidateOpts{
		Digits:    config.DigitsSix,
		Algorithm: config.AlgorithmSHA1,
	})
}

// GenerateCodeCustom uses a counter and secret value and options struct to
// create a passcode.
func GenerateCodeCustom(secret string, counter uint64, opts ValidateOpts) (passcode string, err error) {
	//Set default value
	if opts.Digits == 0 {
		opts.Digits = config.DigitsSix
	}
	// As noted in issue #10 and #17 this adds support for TOTP secrets that are
	// missing their padding.
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	// As noted in issue #24 Google has started producing base32 in lower case,
	// but the StdEncoding (and the RFC), expect a dictionary of only upper case letters.
	secret = strings.ToUpper(secret)

	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", config.ErrValidateSecretInvalidBase32
	}

	buf := make([]byte, 8)
	mac := hmac.New(opts.Algorithm.Hash, secretBytes)
	binary.BigEndian.PutUint64(buf, counter)
	if debug {
		fmt.Printf("counter=%v\n", counter)
		fmt.Printf("buf=%v\n", buf)
	}

	mac.Write(buf)
	sum := mac.Sum(nil)

	// "Dynamic truncation" in RFC 4226
	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := opts.Digits.Length()
	mod := int32(value % int64(math.Pow10(l)))

	if debug {
		fmt.Printf("offset=%v\n", offset)
		fmt.Printf("value=%v\n", value)
		fmt.Printf("mod'ed=%v\n", mod)
	}

	return opts.Digits.Format(mod), nil
}

// ValidateCustom validates an HOTP with customizable options. Most users should
// use Validate().
func ValidateCustom(passcode string, counter uint64, secret string, opts ValidateOpts) (bool, error) {
	passcode = strings.TrimSpace(passcode)

	if len(passcode) != opts.Digits.Length() {
		return false, config.ErrValidateInputInvalidLength
	}

	otpstr, err := GenerateCodeCustom(secret, counter, opts)
	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare([]byte(otpstr), []byte(passcode)) == 1 {
		return true, nil
	}

	return false, nil
}
