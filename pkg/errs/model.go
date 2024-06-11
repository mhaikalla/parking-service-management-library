package errs

import "os"

var (
	// ErrorDebugLocation to include error location.
	ErrorDebugLocation = os.Getenv("ERROR_DEBUG_LOCATION")
	// MaskingErrorMessage to masking error message passing to client.
	MaskingErrorMessage = os.Getenv("MASKING_ERROR_MESSAGE")
)

const (
	// HTTPClientRequestErr when creating request or when requesting.
	// This error mostly happen when connect/write timeout or failed to get bearer token from generator.
	// A few Error when ivalid value provided to `http.Request`.
	HTTPClientRequestErr = 211

	// HTTPCLientRequestBodyErr when creating body request.
	HTTPCLientRequestBodyErr = 212

	// HTTPClientResponseErr when receiving response.
	// This error mostly happen when reponse contains non 200 http code.
	HTTPClientResponseErr = 213

	// HTTPClientResponseBodyErr when read/parse body response.
	// This error mostly happen when body response failed to read/parse.
	HTTPClientResponseBodyErr = 214

	// BearerGenClientError when generate bearer token.
	BearerGenClientError = 241

	// BearerGenClientRequestError when requesting bearer token to upstreams.
	BearerGenClientRequestError = 242

	// BearerGenClientResponseError when cannot read/parse reponse contains bearer token.
	BearerGenClientResponseError = 243

	// BearerGenClientTokenNotFound when bearer token not found on parsed response.
	BearerGenClientTokenNotFound = 244

	// BearerClientTokenExpired when client's bearer token expired.
	// Deprecated use `132` for all problem regarding client's bearer token.
	BearerClientTokenExpired = 245

	// RequestLimited request client limited.
	RequestLimited = 121

	// RequestLimitReached request client reached limit.
	RequestLimitReached = 122

	// RequestBodyMalformed request body client's request not valid.
	RequestBodyMalformed = 111

	// RequestMissingAPIKey client's request missing/invalid `x-api-key` header.
	RequestMissingAPIKey = 131

	// RequestMissingBearer client's request missing/invalid `authorization` bearer token header.
	// Now, this error code will used on all error regarding client's bearer token.
	RequestMissingBearer = 132

	// RequestBearerMalformed client's bearer token cannot be parsed.
	// Deprecated use `132` for all problem regarding client's bearer token.
	RequestBearerMalformed = 133

	// RequestBearerNotValid client's bearer token is invalid/expired.
	// Deprecated use `132` for all problem regarding client's bearer token.
	RequestBearerNotValid = 134

	// RequestNotAllowed request not allowed.
	RequestNotAllowed = 135

	// RequestAmountNotMatch ...
	RequestAmountNotMatch = 136

	// RequestProductNotFound ...
	RequestProductNotFound = 137

	// RequestTransactionIDNotFound ...
	RequestTransactionIDNotFound = 138

	// RequestTransactionIDExpired ...
	RequestTransactionIDExpired = 139

	// RTOPaymentWithBalance ...
	RTOPaymentWithBalance = 140

	// ResponseEncounteredError error when try to respond client's request.
	ResponseEncounteredError = 141

	// GeneralErrorHandler error cannot be mapped.
	GeneralErrorHandler = 151

	// RequestTokenInvalid error if request token is different than configured token.
	RequestTokenInvalid = 161

	// RequestTimestampInvalid error if request timestampt is expired.
	RequestTimestampInvalid = 162

	// RedisError all error thrown by redis.
	RedisError = 411

	// DatabaseError all error thrown by database.
	DatabaseError = 412

	// ResponseNotSuccess ...
	ResponseNotSuccess = 404

	// NotFoundResponse ...
	NotFoundResponse = 404

	// BadRequest ...
	BadRequest = 400

	// Forbidden ...
	Forbidden = 403

	// InvalidToken ...
	InvalidToken = 401

	// RedisTokenNotFound ...
	RedisTokenNotFound = "BR_NOT_FOUND"

	// Not Found
	NotFound = 404

	// Internal Server Error
	InternalServerError = 500

	// Conflict response
	Conflict = 409
)

var (
	// MessageByCode is mapping message to masking original error message.
	MessageByCode = map[string]string{
		"111": "REQUEST_BODY_MALFORMED",
		"121": "REQUEST_LIMITED",
		"122": "REQUEST_LIMIT_REACHED",
		"131": "REQUEST_MISSING_API_KEY",
		"132": "REQUEST_MISSING_BEARER",
		"135": "REQUEST_NOT_ALLOWED",
		"136": "REQUEST_AMOUNT_NOT_MATCH",
		"137": "REQUEST_PRODUCT_NOT_FOUND",
		"138": "REQUEST_TRX_ID_NOT_FOUND",
		"139": "REQUEST_TRX_ID_EXPIRED",
		"140": "REQUEST_TIMEOUT",
		"141": "RESPONSE_ENCOUNTERED_ERROR",
		"151": "GENERAL_ERROR_HANDLER",
		"161": "REQUEST_INVALID",
		"162": "REQUEST_INVALID",

		"211": "HTTP_CLIENT_REQUEST_ERR",
		"212": "HTTP_CLIENT_REQUEST_BODY_ERR",
		"213": "HTTP_CLIENT_RESPONSE_ERR",
		"214": "HTTP_CLIENT_RESPONSE_BODY_ERR",
		"241": "BEARER_GEN_CLIENT_ERROR",
		"242": "BEARER_GEN_CLIENT_REQUEST_ERROR",
		"243": "BEARER_GEN_CLIENT_RESPONSE_ERROR",
		"244": "BEARER_GEN_CLIENT_TOKEN_NOT_FOUND",
		"245": "BEARER_TOKEN_CLIENT_EXPIRED",

		"404": "RESPONSE_NOT_SUCCESS",
		"403": "FORBIDDEN",
		"409": "CONFLICT",
		"411": "REDIS_ERROR",
		"412": "DATABASE_ERROR",
		"500": "INTERNAL_SERVER_ERROR",
	}
)
