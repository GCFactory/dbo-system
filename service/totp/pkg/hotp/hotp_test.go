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
	"github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"encoding/base32"
	"testing"
)

type tc struct {
	Counter uint64
	TOTP    string
	Mode    config.Algorithm
	Secret  string
}

var (
	secSha1 = base32.StdEncoding.EncodeToString([]byte("12345678901234567890"))

	rfcMatrixTCs = []tc{
		{0, "755224", config.AlgorithmSHA1, secSha1},
		{1, "287082", config.AlgorithmSHA1, secSha1},
		{2, "359152", config.AlgorithmSHA1, secSha1},
		{3, "969429", config.AlgorithmSHA1, secSha1},
		{4, "338314", config.AlgorithmSHA1, secSha1},
		{5, "254676", config.AlgorithmSHA1, secSha1},
		{6, "287922", config.AlgorithmSHA1, secSha1},
		{7, "162583", config.AlgorithmSHA1, secSha1},
		{8, "399871", config.AlgorithmSHA1, secSha1},
		{9, "520489", config.AlgorithmSHA1, secSha1},
	}
)

// Test values from http://tools.ietf.org/html/rfc4226#appendix-D
func TestValidateRFCMatrix(t *testing.T) {

	for _, tx := range rfcMatrixTCs {
		valid, err := ValidateCustom(tx.TOTP, tx.Counter, tx.Secret,
			ValidateOpts{
				Digits:    config.DigitsSix,
				Algorithm: tx.Mode,
			})
		require.NoError(t, err,
			"unexpected error totp=%s mode=%v counter=%v", tx.TOTP, tx.Mode, tx.Counter)
		require.True(t, valid,
			"unexpected totp failure totp=%s mode=%v counter=%v", tx.TOTP, tx.Mode, tx.Counter)
	}
}

func TestGenerateRFCMatrix(t *testing.T) {
	for _, tx := range rfcMatrixTCs {
		passcode, err := GenerateCodeCustom(tx.Secret, tx.Counter,
			ValidateOpts{
				Digits:    config.DigitsSix,
				Algorithm: tx.Mode,
			})
		assert.Nil(t, err)
		assert.Equal(t, tx.TOTP, passcode)
	}
}

func TestGenerateCodeCustom(t *testing.T) {
	secSha1 := base32.StdEncoding.EncodeToString([]byte("12345678901234567890"))

	code, err := GenerateCodeCustom("foo", 1, ValidateOpts{})
	print(code)
	require.Equal(t, config.ErrValidateSecretInvalidBase32, err, "Decoding of secret as base32 failed.")
	require.Equal(t, "", code, "Code should be empty string when we have an error.")

	code, err = GenerateCodeCustom(secSha1, 1, ValidateOpts{})
	require.Equal(t, 6, len(code), "Code should be 6 digits when we have not an error.")
	require.NoError(t, err, "Expected no error.")
}

func TestValidateInvalid(t *testing.T) {
	secSha1 := base32.StdEncoding.EncodeToString([]byte("12345678901234567890"))

	valid, err := ValidateCustom("foo", 11, secSha1,
		ValidateOpts{
			Digits:    config.DigitsSix,
			Algorithm: config.AlgorithmSHA1,
		})
	require.Equal(t, config.ErrValidateInputInvalidLength, err, "Expected Invalid length error.")
	require.Equal(t, false, valid, "Valid should be false when we have an error.")

	valid, err = ValidateCustom("foo", 11, secSha1,
		ValidateOpts{
			Digits:    config.DigitsEight,
			Algorithm: config.AlgorithmSHA1,
		})
	require.Equal(t, config.ErrValidateInputInvalidLength, err, "Expected Invalid length error.")
	require.Equal(t, false, valid, "Valid should be false when we have an error.")

	valid, err = ValidateCustom("000000", 11, secSha1,
		ValidateOpts{
			Digits:    config.DigitsSix,
			Algorithm: config.AlgorithmSHA1,
		})
	require.NoError(t, err, "Expected no error.")
	require.Equal(t, false, valid, "Valid should be false.")

	valid = Validate("000000", 11, secSha1)
	require.Equal(t, false, valid, "Valid should be false.")
}

// This tests for issue #10 - secrets without padding
func TestValidatePadding(t *testing.T) {
	valid, err := ValidateCustom("831097", 0, "JBSWY3DPEHPK3PX",
		ValidateOpts{
			Digits:    config.DigitsSix,
			Algorithm: config.AlgorithmSHA1,
		})
	require.NoError(t, err, "Expected no error.")
	require.Equal(t, true, valid, "Valid should be true.")
}

func TestValidateLowerCaseSecret(t *testing.T) {
	valid, err := ValidateCustom("831097", 0, "jbswy3dpehpk3px",
		ValidateOpts{
			Digits:    config.DigitsSix,
			Algorithm: config.AlgorithmSHA1,
		})
	require.NoError(t, err, "Expected no error.")
	require.Equal(t, true, valid, "Valid should be true.")
}
