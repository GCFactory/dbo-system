package http

import (
	"github.com/GCFactory/dbo-system/platform/pkg/httpErrors"
	"net/http"
)

var (
	ErrorNoUserINN = httpErrors.NewRestError(http.StatusBadRequest, "No user INN field", nil)
)
