package portal

const packageVersion = 0.1

// Emulab/Portal response codes (for convenience).
const (
	ResponseSuccess       = 0
	ResponseBadArgs       = 1
	ResponseError         = 2
	ResponseForbidden     = 3
	ResponseBadVersion    = 4
	ResponseServerError   = 5
	ResponseTooBig        = 6
	ResponseRefused       = 7
	ResponseTimedOut      = 8
	ResponseSearchFailed  = 12
	ResponseAlreadyExists = 17
)

type EmulabResponse struct {
	Code   int         `xmlrpc:"code"   json:"code"`
	Value  interface{} `xmlrpc:"value"  json:"value"`
	Output string      `xmlrpc:"output" json:"output"`
}
