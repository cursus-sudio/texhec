package httperrors

import "errors"

var (
	ErrHttp = errors.New("Http error")
)

// services errors
var (
	// Bad Request
	Err400 = errors.Join(ErrHttp, errors.New("400 Bad Request"))
	// Unauthorized
	Err401 = errors.Join(Err400, errors.New("401 Unauthorized"))
	// Payment Required
	Err402 = errors.Join(Err400, errors.New("402 Payment Required"))
	// Forbidden
	Err403 = errors.Join(Err400, errors.New("403 Forbidden"))
	// Not Found
	Err404 = errors.Join(Err400, errors.New("404 Not Found"))
	// Method Not Allowed
	Err405 = errors.Join(Err400, errors.New("405 Method Not Allowed"))
	// Not Acceptable
	Err406 = errors.Join(Err400, errors.New("406 Not Acceptable"))
	// Proxy Authentication Required
	Err407 = errors.Join(Err400, errors.New("407 Proxy Authentication Required"))
	// Request Timeout
	Err408 = errors.Join(Err400, errors.New("408 Request Timeout"))
	// Conflict
	Err409 = errors.Join(Err400, errors.New("409 Conflict"))
	// Gone
	Err410 = errors.Join(Err400, errors.New("410 Gone"))
	// Length Required
	Err411 = errors.Join(Err400, errors.New("411 Length Required"))
	// Precondition Failed
	Err412 = errors.Join(Err400, errors.New("412 Precondition Failed"))
	// Payload Too Large
	Err413 = errors.Join(Err400, errors.New("413 Payload Too Large"))
	// URI Too Long
	Err414 = errors.Join(Err400, errors.New("414 URI Too Long"))
	// Unsupported Media Type
	Err415 = errors.Join(Err400, errors.New("415 Unsupported Media Type"))
	// Range Not Satisfiable
	Err416 = errors.Join(Err400, errors.New("416 Range Not Satisfiable"))
	// Expectation Failed
	Err417 = errors.Join(Err400, errors.New("417 Expectation Failed"))
	// Teapot (RFC 2324)
	Err418 = errors.Join(Err400, errors.New("418 I'm a teapot"))
	// Misdirected Request
	Err421 = errors.Join(Err400, errors.New("421 Misdirected Request"))
	// Unprocessable Entity
	Err422 = errors.Join(Err400, errors.New("422 Unprocessable Entity"))
	// Locked
	Err423 = errors.Join(Err400, errors.New("423 Locked"))
	// Failed Dependency
	Err424 = errors.Join(Err400, errors.New("424 Failed Dependency"))
	// Too Early
	Err425 = errors.Join(Err400, errors.New("425 Too Early"))
	// Upgrade Required
	Err426 = errors.Join(Err400, errors.New("426 Upgrade Required"))
	// Precondition Required
	Err428 = errors.Join(Err400, errors.New("428 Precondition Required"))
	// Too Many Requests
	Err429 = errors.Join(Err400, errors.New("429 Too Many Requests"))
	// Request Header Fields Too Large
	Err431 = errors.Join(Err400, errors.New("431 Request Header Fields Too Large"))
	// Unavailable For Legal Reasons
	Err451 = errors.Join(Err400, errors.New("451 Unavailable For Legal Reasons"))
)

var (
	// Internal Server Error
	Err500 = errors.Join(ErrHttp, errors.New("500 Internal Server Error"))
	// Not Implemented
	Err501 = errors.Join(Err500, errors.New("501 Not Implemented"))
	// Bad Gateway
	Err502 = errors.Join(Err500, errors.New("502 Bad Gateway"))
	// Service Unavailable
	Err503 = errors.Join(Err500, errors.New("503 Service Unavailable"))
	// Gateway Timeout
	Err504 = errors.Join(Err500, errors.New("504 Gateway Timeout"))
	// HTTP Version Not Supported
	Err505 = errors.Join(Err500, errors.New("505 HTTP Version Not Supported"))
	// Variant Also Negotiates
	Err506 = errors.Join(Err500, errors.New("506 Variant Also Negotiates"))
	// Insufficient Storage (WebDAV)
	Err507 = errors.Join(Err500, errors.New("507 Insufficient Storage"))
	// Loop Detected (WebDAV)
	Err508 = errors.Join(Err500, errors.New("508 Loop Detected"))
	// Not Extended
	Err510 = errors.Join(Err500, errors.New("510 Not Extended"))
	// Network Authentication Required
	Err511 = errors.Join(Err500, errors.New("511 Network Authentication Required"))
)

func IsHttpError(err error) bool {
	return errors.Is(err, ErrHttp)
}

func IsClientError(err error) bool {
	return errors.Is(err, Err400)
}

func IsServiceError(err error) bool {
	return errors.Is(err, Err500)
}
