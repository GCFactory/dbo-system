package http

import (
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"net/http"
)

var (
	ErrorNoUserId   = httpErrors.NewRestError(http.StatusBadRequest, "No user_id field", nil)
	ErrorNoUserName = httpErrors.NewRestError(http.StatusBadRequest, "No user_name field", nil)
)
