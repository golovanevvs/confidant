package customerrors

import (
	"errors"
)

const (
	OrderAlredyUploadedThisUser200  = "номер заказа уже был загружен этим пользователем"
	EmptyOrder204                   = "нет данных для ответа"
	EmptyWithdrawals204             = "нет ни одного списания"
	ASOrderNotRegistered204         = "заказ не зарегистрирован в системе расчёта"
	InvalidRequest400               = "неверный формат запроса"
	JWTParseError401                = "ошибка при чтении JWT"
	NotEnoughPoints402              = "на счету недостаточно средств"
	OrderAlredyUploadedOtherUser409 = "номер заказа уже был загружен другим пользователем"
	InvalidOrderNumber422           = "Неверный формат номера заказа"
	InvalidOrderNumberNotInt422     = "Неверный формат номера заказа: не соответствует типу int"
	ASTooManyRequests429            = "превышено количество запросов к сервису"
	InternalServerError500          = "внутренняя ошибка сервера"
	DecodeJSONError500              = "ошибка десериализации JSON"
	EncodeJSONError500              = "ошибка сериализации JSON"
	ResponseBodyError500            = "ошибка при чтении тела ответа"
	AtoiError500                    = "ошибка преобразования строки в число"
	ClientError500                  = "ошибка при отправке запроса"
	ASError                         = "сервис по взаимодействию с системой расчёта начислений баллов"
)

const (
	ClientMsg = "CLIENT"
)

const (
	ClientAppViewErr = "application view"
	ClientHTTPErr    = "transport HTTP"
	ClientServiceErr = "service"
)

var (
	ErrRunAppView               = errors.New("application view error")
	ErrCreateRequest            = errors.New("creating request error")
	ErrSendRequest              = errors.New("sending request error")
	ErrReadResponseBody         = errors.New("error reading response body")
	ErrServerStatus             = errors.New("error get server status")
	ErrAuthHeaderReq            = errors.New("the request does not contain the \"Authorization\" header")
	ErrAuthHeaderResp           = errors.New("the response does not contain the \"Authorization\" header")
	ErrInvalidAuthHeaderReq     = errors.New("the request contains the invalid \"Authorization\" header")
	ErrInvalidAuthHeaderResp    = errors.New("the response contains the invalid \"Authorization\" header")
	ErrBearer                   = errors.New("the \"Authorization\" header does not contain \"Bearer\"")
	ErrAccessToken              = errors.New("the \"Authorization\" header does not contain a access token")
	ErrRefreshToken             = errors.New("the response does not contain the \"Refresh-Token\" header")
	ErrGenPasswordHash          = errors.New("generate password hash error")
	ErrGenHash                  = errors.New("generate hash error")
	ErrGenRefreshTokenHash      = errors.New("generate a refresh token hash error")
	ErrSaveRefreshToken         = errors.New("save refresh token error")
	ErrNoActiveRefreshToken     = errors.New("no active refresh token")
	ErrSaveActiveRefreshToken   = errors.New("save active refresh token error")
	ErrDeleteActiveRefreshToken = errors.New("delete active refresh token error")
)

const (
	ServerMsg = "SERVER"
)

const (
	AccountErr        = "account model"
	HandlerErr        = "transport HTTP"
	AccountServiceErr = "account service"
	ManageServiceErr  = "manage service"
	DBErr             = "DB"
)

var (
	ErrDBBusyEmail409             = errors.New("e-mail is already busy")
	ErrDBEmailNotFound401         = errors.New("there is no account with this email address")
	ErrDBWrongPassword401         = errors.New("wrong password")
	ErrDBInternalError500         = errors.New("DB internal error")
	ErrAccountValidateEmail422    = errors.New("invalid e-mail")
	ErrAccountValidatePassword422 = errors.New("invalid password")
	ErrContentType400             = errors.New("invalid Content-Type")
	ErrDecodeJSON400              = errors.New("deserializing JSON error")
	ErrEncodeJSON500              = errors.New("JSON serialization error")
	ErrInternalServerError500     = errors.New("internal server error")
	ErrTokenSignedString          = errors.New("signed string token error")
	ErrJWTWrongSingingMethod401   = errors.New("invalid signature method")
	ErrJWTInvalidToken401         = errors.New("invalid token")
	ErrJWTExpiredToken401         = errors.New("expired token")
)
