package http

import (
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"net/http"
)

var (
	ErrorNoUserId   = httpErrors.NewRestError(http.StatusBadRequest, "No user_id field", nil)
	ErrorNoUserName = httpErrors.NewRestError(http.StatusBadRequest, "No user_name field", nil)
	ErrorNoUrl      = httpErrors.NewRestError(http.StatusBadRequest, "No totp_url field", nil)
	ErrorNoTotpCode = httpErrors.NewRestError(http.StatusBadRequest, "No totp_code field", nil)
)
