package customerrors

import (
	"errors"
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
	ErrRunAppView            = errors.New("application view error")
	ErrCreateRequest         = errors.New("creating request error")
	ErrSendRequest           = errors.New("sending request error")
	ErrReadResponseBody      = errors.New("error reading response body")
	ErrServerStatus          = errors.New("error get server status")
	ErrAuthHeaderReq         = errors.New("the request does not contain the \"Authorization\" header")
	ErrAuthHeaderResp        = errors.New("the response does not contain the \"Authorization\" header")
	ErrInvalidAuthHeaderReq  = errors.New("the request contains the invalid \"Authorization\" header")
	ErrInvalidAuthHeaderResp = errors.New("the response contains the invalid \"Authorization\" header")
	ErrBearer                = errors.New("the \"Authorization\" header does not contain \"Bearer\"")
	ErrAccessToken           = errors.New("the \"Authorization\" header does not contain a access token")
	ErrRefreshToken          = errors.New("the response does not contain the \"Refresh-Token\" header")
	ErrGenPasswordHash       = errors.New("generate password hash error")
	ErrGenHash               = errors.New("generate hash error")
	ErrGenRefreshTokenHash   = errors.New("generate a refresh token hash error")
	ErrSaveRefreshToken      = errors.New("save refresh token error")
	ErrNoActiveAccount       = errors.New("no active account")
	ErrSaveActiveAccount     = errors.New("save active account")
	ErrDeleteActiveAccount   = errors.New("delete active account error")
	ErrLoadActiveAccount     = errors.New("load active account error")
	ErrLoadEmail             = errors.New("load email error")
	ErrLogout                = errors.New("logout error")
	ErrAddGroup              = errors.New("group is already exists")
	ErrAddEmailInGroup       = errors.New("add e-email in group error")
	ErrAddNote               = errors.New("add note error")
	ErrAddPass               = errors.New("add password error")
	ErrAddCard               = errors.New("add card error")
	ErrAddFile               = errors.New("add file error")
	ErrGetNote               = errors.New("get note error")
	ErrGetPass               = errors.New("get password error")
	ErrGetCard               = errors.New("get card error")
	ErrGetFile               = errors.New("get file error")
	ErrGetFileForSave        = errors.New("get file for save error")
	ErrDecryptEmptyBody      = errors.New(("empty encrypted body"))
)

const (
	ServerMsg = "SERVER"
)

const (
	AccountErr        = "account model"
	HandlerErr        = "transport HTTP"
	AccountServiceErr = "account service"
	ManageServiceErr  = "manage service"
	GroupsServiceErr  = "groups service"
	DataServiceErr    = "data service"
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
	ErrDebug                      = errors.New("error for debug")
)
