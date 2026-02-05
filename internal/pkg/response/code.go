package response

import "gitlab.com/skeleton-golang/internal/pkg/constants"

const (
	SuccessCode      string = constants.STATUS_SUCCESS
	ErrorInvalidJson string = constants.STATUS_INVALID_REQUEST_FORMAT
	GeneralError     string = constants.STATUS_UNKNOWN_ERROR
)
