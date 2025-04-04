package usecase

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/GCFactory/dbo-system/service/totp/internal/totp"
	totpPkg "github.com/GCFactory/dbo-system/service/totp/pkg"
	totpErrors "github.com/GCFactory/dbo-system/service/totp/pkg/errors"
	totpPkgConfig "github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type totpUC struct {
	cfg       *config.Config
	totpRepo  totp.Repository
	totpLogic totpPkg.Totp
	logger    logger.Logger
}

func (t totpUC) Enroll(ctx context.Context, totpConfig models.TOTPConfig) (*models.TOTPEnroll, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enroll")
	defer span.Finish()

	if totpConfig.AccountName == "" {
		return nil, totpErrors.NoUserName
	}
	if totpConfig.UserId == uuid.Nil {
		return nil, totpErrors.EmptyUserId
	}

	_, err := t.totpRepo.GetActiveConfig(ctxWithTrace, totpConfig.UserId)
	if err == nil {
		return nil, totpErrors.ActiveTotp
	}

	totpId := uuid.New()

	totpConfig.Id = totpId
	totpConfig.IsActive = true
	totpConfig.Issuer = "dbo.gcfactory.space"

	secret, url, err := t.totpLogic.Generate(totpPkgConfig.GenerateOpts{
		Issuer:      totpConfig.Issuer,
		AccountName: totpConfig.AccountName,
		SecretSize:  totpPkgConfig.DefaultSecretLength,
		Algorithm:   totpPkgConfig.DefaultAlgorithm,
	})

	if err != nil {
		return nil, ErrorGenTotp
	}
	totpConfig.Secret = *secret
	totpConfig.URL = *url

	err = t.totpRepo.CreateConfig(ctxWithTrace, totpConfig)
	if err != nil {
		return nil, ErrorCreateConfig
	}

	return &models.TOTPEnroll{
		TotpId:     totpConfig.Id.String(),
		TotpSecret: *secret,
		TotpUrl:    *url,
	}, nil
}

func (t totpUC) Verify(ctx context.Context, url string) (*models.TOTPVerify, error) {
	algorithm := strings.Contains(url, "algorithm")
	if !algorithm {
		return &models.TOTPVerify{Status: totpErrors.NoAlgorithmField.Error()}, totpErrors.NoAlgorithmField
	}
	digits := strings.Contains(url, "digits")
	if !digits {
		return &models.TOTPVerify{Status: totpErrors.NoDigitsField.Error()}, totpErrors.NoDigitsField
	}
	issuer := strings.Contains(url, "issuer")
	if !issuer {
		return &models.TOTPVerify{Status: totpErrors.NoIssuerField.Error()}, totpErrors.NoIssuerField
	}
	period := strings.Contains(url, "period")
	if !period {
		return &models.TOTPVerify{Status: totpErrors.NoPeriodField.Error()}, totpErrors.NoPeriodField
	}
	secret := strings.Contains(url, "secret")
	if !secret {
		return &models.TOTPVerify{Status: totpErrors.NoSecretField.Error()}, totpErrors.NoSecretField
	}
	return &models.TOTPVerify{Status: "OK"}, nil
}

func (t totpUC) Validate(ctx context.Context, id uuid.UUID, code string, time time.Time) (*models.TOTPValidate, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Validate")
	defer span.Finish()

	activeConfig, err := t.totpRepo.GetActiveConfig(ctxWithTrace, id)
	if err != nil {
		return &models.TOTPValidate{Status: totpErrors.NoUserId.Error()}, totpErrors.NoUserId
	}

	url := activeConfig.URL

	validateOpts := totpPkgConfig.ValidateOpts{}

	//Extract secret
	reg_expr, err := regexp.Compile("secret=([\\d\\w]{32})")
	if err != nil {
		return nil, ErrorRegexCompile
	}

	secret := reg_expr.FindString(url)
	secret = secret[7:]

	//Extract period
	reg_expr, err = regexp.Compile("period=(\\d{1,2})")
	if err != nil {
		return nil, ErrorRegexCompile
	}

	var period uint

	period_str := reg_expr.FindString(url)
	if period_str == "" {
		period = totpPkgConfig.DefaultPeriod
	} else {
		period_str = period_str[7:]
		tmp, err := strconv.ParseUint(period_str, 10, 32)
		if err != nil {
			return nil, ErrorRegexCompile
		}
		period = uint(tmp)
	}

	validateOpts.Period = period

	//Extract digits
	reg_expr, err = regexp.Compile("digits=(6|8)")
	if err != nil {
		return nil, ErrorRegexCompile
	}

	var digits totpPkgConfig.Digits

	digits_str := reg_expr.FindString(url)
	if digits_str == "" {
		digits = totpPkgConfig.DefaultDigits
	} else {
		digits_str = digits_str[7:]
		tmp, err := strconv.ParseInt(digits_str, 10, 8)
		if err != nil {
			return nil, ErrorRegexCompile
		}
		digits = totpPkgConfig.Digits(tmp)
	}

	validateOpts.Digits = digits

	//Extract algorithm
	reg_expr, err = regexp.Compile("algorithm=([\\d\\w]{2,6})")
	if err != nil {
		return nil, ErrorRegexCompile
	}

	var algorithm totpPkgConfig.Algorithm

	algorithm_str := reg_expr.FindString(url)
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
		}
	}

	validateOpts.Algorithm = algorithm

	serverCode, err := t.totpLogic.GenerateCodeCustom(secret, time, validateOpts)
	if err != nil {
		return nil, ErrorGenCodeCustom
	}
	if serverCode != code {
		return &models.TOTPValidate{Status: totpErrors.WrongTotpCode.Error()}, totpErrors.WrongTotpCode
	}
	return &models.TOTPValidate{Status: "OK"}, nil
}

func (t totpUC) Enable(ctx context.Context, totpId uuid.UUID, userId uuid.UUID) (*models.TOTPEnable, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Enable")
	defer span.Finish()

	totpCfg := &models.TOTPConfig{
		Id: uuid.Nil,
	}
	var err error

	if totpId != uuid.Nil && userId != uuid.Nil {
		totpCfg, err = t.totpRepo.GetConfigByTotpId(ctxWithTrace, totpId)
		if err != nil || totpCfg.UserId != userId {
			return &models.TOTPEnable{Status: totpErrors.NoId.Error()}, totpErrors.NoId
		}
	}
	if totpId != uuid.Nil {
		if totpCfg.Id == uuid.Nil {
			totpCfg, err = t.totpRepo.GetConfigByTotpId(ctxWithTrace, totpId)
		}
		if err != nil {
			return &models.TOTPEnable{Status: totpErrors.NoTotpId.Error()}, totpErrors.NoTotpId
		} else {
			if totpCfg.IsActive == true {
				return &models.TOTPEnable{Status: totpErrors.TotpIsActive.Error()}, totpErrors.TotpIsActive
			}
			activeConfig, err := t.totpRepo.GetActiveConfig(ctxWithTrace, totpCfg.UserId)
			if err == nil && activeConfig.IsActive == true {
				return &models.TOTPEnable{Status: totpErrors.TotpIsActive.Error()}, totpErrors.TotpIsActive
			}
			err = t.totpRepo.UpdateTotpActivityByTotpId(ctxWithTrace, totpId, true)
			if err != nil {
				return nil, ErrorUpdateActivityByTotpId
			}
			return &models.TOTPEnable{Status: "OK"}, nil
		}
	}
	if userId != uuid.Nil {
		totpCfg, err = t.totpRepo.GetConfigByUserId(ctxWithTrace, userId)
		if err != nil {
			return &models.TOTPEnable{Status: totpErrors.NoUserId.Error()}, totpErrors.NoUserId
		}
		totpCfg, err = t.totpRepo.GetActiveConfig(ctxWithTrace, userId)
		if err == nil {
			return &models.TOTPEnable{Status: totpErrors.TotpIsActive.Error()}, totpErrors.TotpIsActive
		}
		totpCfg, err = t.totpRepo.GetLastDisabledConfig(ctxWithTrace, userId)
		if err != nil {
			return &models.TOTPEnable{Status: totpErrors.TotpIsActive.Error()}, totpErrors.TotpIsActive
		}
		err = t.totpRepo.UpdateTotpActivityByTotpId(ctxWithTrace, totpCfg.Id, true)
		if err != nil {
			return nil, ErrorUpdateActivityByTotpId
		}
		return &models.TOTPEnable{Status: "OK"}, nil
	}
	return &models.TOTPEnable{Status: totpErrors.NoId.Error()}, totpErrors.NoId
}

func (t totpUC) Disable(ctx context.Context, totpId uuid.UUID, userId uuid.UUID) (*models.TOTPDisable, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
	defer span.Finish()

	totpCfg := &models.TOTPConfig{
		Id: uuid.Nil,
	}

	var err error

	if totpId != uuid.Nil && userId != uuid.Nil {
		totpCfg, err = t.totpRepo.GetConfigByTotpId(ctxWithTrace, totpId)
		if err != nil || totpCfg.UserId != userId {
			return &models.TOTPDisable{Status: totpErrors.NoId.Error()}, totpErrors.NoId
		}
	}
	if totpId != uuid.Nil {
		if totpCfg.Id == uuid.Nil {
			totpCfg, err = t.totpRepo.GetConfigByTotpId(ctxWithTrace, totpId)
		}
		if err != nil {
			return &models.TOTPDisable{Status: totpErrors.NoTotpId.Error()}, totpErrors.NoTotpId
		} else {
			if totpCfg.IsActive == false {
				return &models.TOTPDisable{Status: totpErrors.TotpIsDisabled.Error()}, totpErrors.TotpIsDisabled
			}
			err = t.totpRepo.UpdateTotpActivityByTotpId(ctxWithTrace, totpId, false)
			if err != nil {
				return nil, ErrorUpdateActivityByTotpId
			}
			return &models.TOTPDisable{Status: "OK"}, nil
		}
	}
	if userId != uuid.Nil {
		totpCfg, err = t.totpRepo.GetConfigByUserId(ctxWithTrace, userId)
		if err != nil {
			return &models.TOTPDisable{Status: totpErrors.NoTotpId.Error()}, totpErrors.NoTotpId
		} else {
			totpCfg, err = t.totpRepo.GetActiveConfig(ctxWithTrace, userId)
			if err != nil {
				return &models.TOTPDisable{Status: totpErrors.TotpIsDisabled.Error()}, totpErrors.TotpIsDisabled
			}
			err = t.totpRepo.UpdateTotpActivityByTotpId(ctxWithTrace, totpCfg.Id, false)
			if err != nil {
				return nil, ErrorUpdateActivityByTotpId
			}
			return &models.TOTPDisable{Status: "OK"}, nil
		}
	}
	return &models.TOTPDisable{Status: totpErrors.NoId.Error()}, totpErrors.NoId
}

func (t totpUC) Url(ctx context.Context, userId uuid.UUID) (*models.TOTPEnroll, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "totpUC.Disable")
	defer span.Finish()

	totpInfo, err := t.totpRepo.GetActiveConfig(ctxWithTrace, userId)
	if err != nil {
		return nil, err
	}

	return &models.TOTPEnroll{
		TotpUrl: totpInfo.URL,
	}, nil

}

func NewTOTPUseCase(cfg *config.Config, totpRepo totp.Repository, totpLogic totpPkg.Totp, log logger.Logger) totp.UseCase {
	return &totpUC{cfg: cfg, totpRepo: totpRepo, totpLogic: totpLogic, logger: log}
}
