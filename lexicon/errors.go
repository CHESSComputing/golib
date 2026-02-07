package lexicon

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
)

// GenericErr represents generic lexicon error
var GenericErr = errors.New("lexicon error")

// InvalidParamErr represents generic error for invalid input parameter
var InvalidParamErr = errors.New("invalid parameter(s)")

// PatterErr represents generic error for invalid input parameter
var PatternErr = errors.New("invalid pattern")

// ValidationErr represents generic validation error
var ValidationErr = errors.New("validation error")

// NotImplementedApiErr represents generic not implemented api error
var NotImplementedApiErr = errors.New("not implemented api error")

// InvalidRequestErr represents generic invalid request error
var InvalidRequestErr = errors.New("invalid request error")

// Lexicon Error codes provides static representation of Lexicon errors, they cover 1xx range
const (
	GenericErrorCode   = iota + 10 // generic Lexicon error
	ReaderErrorCode                // 11 reader error
	UnmarshalErrorCode             // 12 unmarshal error
	MarshalErrorCode               // 13 marshal error
	PatternErrorCode               // 14 pattern error
	ValidateErrorCode              // 15 validate error
)

// helper function to create Lexicon error
func Error(err error, code int, msg, function string) error {
	reason := "nil"
	if err != nil {
		reason = err.Error()
	}
	stackSlice := make([]byte, 1024*4)
	s := runtime.Stack(stackSlice, false)
	return &LexiconError{
		Reason:     reason,
		Message:    msg,
		Code:       code,
		Function:   function,
		Stacktrace: fmt.Sprintf("\n%s", stackSlice[0:s]),
	}
}

// LexiconError represents common structure for Lexicon errors
type LexiconError struct {
	Reason     string `json:"reason"`     // error string
	Message    string `json:"message"`    // additional message describing the issue
	Function   string `json:"function"`   // Lexicon function
	Code       int    `json:"code"`       // Lexicon error code
	Stacktrace string `json:"stacktrace"` // Go stack trace
}

// Error function implements details of Lexicon error message
func (e *LexiconError) Error() string {
	sep := ": "
	if strings.Contains(e.Reason, "LexiconError") { // nested error
		sep += "nested "
	}
	msg := fmt.Sprintf(
		"LexiconError Code:%d Description:%s Function:%s Message:%s Error%s%v",
		e.Code, e.Explain(), e.Function, e.Message, sep, e.Reason)
	log.Println(msg)
	rmsg := fmt.Sprintf(
		"LexiconError Code:%d Reason:%s",
		e.Code, e.Message)
	if strings.Contains(sep, "nested") {
		arr := strings.Split(e.Reason, sep)
		rrr := strings.Split(arr[len(arr)-1], "Reason:")
		rmsg = rrr[len(rrr)-1]
	}
	return rmsg
}

// ErrorStacktrace function implements details of Lexicon error message and stacktrace
func (e *LexiconError) ErrorStacktrace() string {
	sep := ": "
	if strings.Contains(e.Reason, "LexiconError") { // nested error
		sep += "nested "
	}
	return fmt.Sprintf(
		"LexiconError Code:%d Description:%s Function:%s Message:%s Error%s%v Stacktrace: %v",
		e.Code, e.Explain(), e.Function, e.Message, sep, e.Reason, e.Stacktrace)
}

// Explain provides explanation about the error
func (e *LexiconError) Explain() string {
	switch e.Code {
	case ReaderErrorCode:
		return "Lexicon reader I/O error, e.g. unable to read HTTP POST payload"
	case UnmarshalErrorCode:
		return "Lexicon unable to parse JSON record"
	case MarshalErrorCode:
		return "Lexicon unable to convert record to JSON"
	case GenericErrorCode:
		return "Generic Lexicon error"
	case PatternErrorCode:
		return "Lexicon validation error when wrong pattern is provided"
	case ValidateErrorCode:
		return "Lexicon validation error, e.g. input parameter does not match lexicon rules"
	default:
		return "Not defined"
	}
	return "Not defined"
}
