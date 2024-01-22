package usecase

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	totpPkg "github.com/GCFactory/dbo-system/service/totp/pkg/otp"
	totpPkgConfig "github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type totpUC struct {
	cfg      *config.Config
	totpRepo totp.Repository
	logger   logger.Logger
}

func (t totpUC) Enroll(ctx context.Context, totpConfig models.TOTPConfig) (*models.TOTPEnroll, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "totpUC.Enroll")
	defer span.Finish()

	totpConfig.Id = uuid.New()
	totpConfig.AccountName = "admin"
	totpConfig.IsActive = true
	totpConfig.Issuer = "dbo.gcfactory.space"

	secret, url, err := totpPkg.Generate(totpPkgConfig.GenerateOpts{
		Issuer:      totpConfig.Issuer,
		AccountName: totpConfig.AccountName,
		SecretSize:  20,
		Algorithm:   totpPkgConfig.AlgorithmSHA256,
	})

	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "totpUC.Enroll.Generate"))
	}
	totpConfig.Secret = *secret
	totpConfig.URL = *url

	err = t.totpRepo.CreateConfig(ctx, totpConfig)
	if err != nil {
		return nil, httpErrors.NewRestError(http.StatusBadRequest, err.Error(), nil)
	}

	totpEnroll := models.TOTPEnroll{
		Base32:    totpConfig.Secret,
		OTPathURL: *url,
	}

	return &totpEnroll, nil
}

func (t totpUC) Verify(ctx context.Context, url string) (*models.TOTPVerify, error) {
	if len(url) == 0 {
		return &models.TOTPVerify{Status: "Empty TOTP URL"}, errors.New("Empty TOTP URL")
	}
	algorithm := strings.Contains(url, "algorithm")
	if !algorithm {
		return &models.TOTPVerify{Status: "No algorithm field"}, errors.New("URL hadn't algorithm field")
	}
	digits := strings.Contains(url, "digits")
	if !digits {
		return &models.TOTPVerify{Status: "No digits field"}, errors.New("URL hadn't digits field")
	}
	issuer := strings.Contains(url, "issuer")
	if !issuer {
		return &models.TOTPVerify{Status: "No issuer field"}, errors.New("URL hadn't issuer field")
	}
	period := strings.Contains(url, "period")
	if !period {
		return &models.TOTPVerify{Status: "No period field"}, errors.New("URL hadn't period field")
	}
	secret := strings.Contains(url, "secret")
	if !secret {
		return &models.TOTPVerify{Status: "No secret field"}, errors.New("URL hadn't secret field")
	}
	return &models.TOTPVerify{Status: "OK"}, nil
}

func (t totpUC) Validate(ctx context.Context, id uuid.UUID, code string) (*models.TOTPValidate, error) {
	if len(code) == 0 {
		return &models.TOTPValidate{Status: "Empty code"}, errors.New("Empty code")
	}
	url, err := t.totpRepo.GetURL(ctx, id)
	if err != nil {
		return &models.TOTPValidate{Status: "Non-existent user id"}, nil
	}
	validateOpts := totpPkgConfig.ValidateOpts{}

	//Extract secret
	reg_expr, err := regexp.Compile("secret=([\\d\\w]{32})")
	if err != nil {
		return &models.TOTPValidate{Status: "Secret match compile error"}, errors.New("Secret match compile error")
	}
	secret := reg_expr.FindString(url)
	secret = secret[7:]

	//Extract period
	reg_expr, err = regexp.Compile("period=(\\d{1,2})")
	if err != nil {
		return &models.TOTPValidate{Status: "Period match compile error"}, errors.New("Period match compile error")
	}
	period_str := reg_expr.FindString(url)
	var period uint
	if period_str == "" {
		period = totpPkgConfig.DefaultPeriod
	} else {
		period_str = period_str[7:]
		tmp, err := strconv.ParseUint(period_str, 10, 32)
		if err != nil {
			return &models.TOTPValidate{Status: "Parse period string to uint error"}, errors.New("Parse period string to uint error")
		}
		period = uint(tmp)
	}
	validateOpts.Period = period

	//Extract digits
	reg_expr, err = regexp.Compile("digits=(6|8)")
	if err != nil {
		return &models.TOTPValidate{Status: "Digits match compile error"}, errors.New("Digits match compile error")
	}
	digits_str := reg_expr.FindString(url)
	var digits totpPkgConfig.Digits
	if digits_str == "" {
		digits = totpPkgConfig.DefaultDigits
	} else {
		digits_str = digits_str[7:]
		tmp, err := strconv.ParseInt(digits_str, 10, 8)
		if err != nil {
			return &models.TOTPValidate{Status: "Parse digits string to uint error"}, errors.New("Parse digits string to uint error")
		}
		digits = totpPkgConfig.Digits(tmp)
	}
	validateOpts.Digits = digits

	//Extract algorithm
	reg_expr, err = regexp.Compile("algorithm=([\\d\\w]{2,6})")
	if err != nil {
		return &models.TOTPValidate{Status: "Digits match compile error"}, errors.New("Digits match compile error")
	}
	algorithm_str := reg_expr.FindString(url)
	var algorithm totpPkgConfig.Algorithm
	if algorithm_str == "" {
		algorithm = totpPkgConfig.DefaultAlgorithm
	} else {
		algorithm_str = algorithm_str[10:]
		if algorithm_str == "MD5" {
			algorithm = totpPkgConfig.AlgorithmMD5
		} else if algorithm_str == "SHA1" {
			algorithm = totpPkgConfig.AlgorithmSHA1
		} else if algorithm_str == "SHA256" {
			algorithm = totpPkgConfig.AlgorithmSHA256
		} else if algorithm_str == "SHA512" {
			algorithm = totpPkgConfig.AlgorithmSHA512
		} else {
			algorithm = totpPkgConfig.DefaultAlgorithm
		}
	}
	validateOpts.Algorithm = algorithm

	serverCode, err := totpPkg.GenerateCodeCustom(secret, time.Now(), validateOpts)
	if err != nil {
		return &models.TOTPValidate{Status: "Generate code error"}, errors.New("Generate code error")
	}
	if serverCode != code {
		return &models.TOTPValidate{Status: "Incorrect code"}, errors.New("Incorrect")
	}
	return &models.TOTPValidate{Status: "OK"}, nil
}

func (t totpUC) Enable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func (t totpUC) Disable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error) {
	//TODO implement me
	panic("implement me")
}

func NewTOTPUseCase(cfg *config.Config, totpRepo totp.Repository, log logger.Logger) totp.UseCase {
	return &totpUC{cfg: cfg, totpRepo: totpRepo, logger: log}
}
