package constants

const (
	ContentType         = "Content-Type"
	ApplicationJson     = "application/json"
	ApplicationJsonUtf8 = "application/json; charset=utf-8"
	FormUrlEncoded      = "application/x-www-form-urlencoded"
	CtxTokenKey         = "token"
	CtxClaimsKey        = "claims"
	LangEnUS            = "en-US"

	// SecretKey used to Sign the JWT Token
	// Ideally we would want this in an env var instead.
	SecretKey = "my_secret_key"

	OrderStatusPending = "pending"
	OrderStatusSuccess = "success"
	OrderStatusFailed  = "failed"

	// Viva eventTypeIds - https://developer.viva.com/webhooks-for-payments/
	TransactionPaymentCreated = 1796
	TransactionFailed         = 1798
)
