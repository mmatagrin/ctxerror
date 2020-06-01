package ctxerror

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
)

var HiddenFields = []string{}

type CtxErrorManager struct {
	context map[string]interface{}
}

type CtxErrorTraceI interface {
	Error() string
	ErrorJson() string
	GetMessage() string
	GetTrace() []CtxError
	AddError(error, string) CtxErrorTraceI
}

type CtxErrorTrace struct {
	Trace      []CtxError `json:"trace"`
	StackTrace string     `json:"stack_trace"`
}

type CtxError struct {
	Message      string                 `json:"message"`
	FileName     string                 `json:"file_name"`
	Line         int                    `json:"line"`
	FunctionName string                 `json:"function_name"`
	Context      map[string]interface{} `json:"context"`
	ErrorS       string                 `json:"error"`
	ErrorI   	 error
}

func (cet CtxErrorTrace) GetMessage() string {
	if cet.Trace != nil && len(cet.Trace) > 0{
		return cet.Trace[0].GetMessage()
	}

	return ""
}

func (cet CtxErrorTrace) Error() string {
	ctxErrorTraceBytes, err := json.MarshalIndent(cet, "", "   ")
	if err != nil {
		return string(ctxErrorTraceBytes)
	}

	if len(cet.Trace) > 0 {
		return fmt.Sprintf("%s\n%v", cet.Trace[0].Message, string(ctxErrorTraceBytes))
	}

	return fmt.Sprintf("%v", string(ctxErrorTraceBytes))

}

func (cet CtxErrorTrace) ErrorJson() string {
	ctxErrorTraceBytes, err := json.MarshalIndent(cet, "", "   ")
	if err != nil {
		return fmt.Sprintf("%v", ctxErrorTraceBytes)
	}

	return string(ctxErrorTraceBytes)
}

func  (cet CtxErrorTrace) GetTrace() []CtxError{
	if cet.Trace == nil {
		return []CtxError{}
	}

	return cet.Trace
}

func (cet CtxErrorTrace) AddError(err error, message string) CtxErrorTraceI {
	if err == nil{
		return cet
	}

	if errTrace, ok := err.(CtxErrorTraceI); ok {
		cet.Trace = append(errTrace.GetTrace(), cet.Trace...)
		return cet
	}

	ctxError := getContextualizedError(message, nil)
	ctxError.ErrorS = err.Error()
	ctxError.ErrorI = err
	cet.Trace = append([]CtxError{ctxError}, cet.Trace...)

	return cet
}

func (cem CtxErrorManager) AddContext(key string, val interface{}) {
	if cem.context == nil {
		cem.context = make(map[string]interface{})
	}

	for _, hiddenField := range HiddenFields {
		if key == hiddenField {
			cem.context[key] = "hidden"
			return
		}
	}

	cem.context[key] = val
}

func (ctxError CtxError) Error() string {
	contextualizedErrorBytes, err := json.MarshalIndent(ctxError, "", "   ")
	if err != nil {
		return fmt.Sprintf("%v", ctxError)
	}

	return string(contextualizedErrorBytes)
}

func (cem CtxErrorManager) Wrap(err error, message string) CtxErrorTraceI {
	if err == nil{
		return nil
	}

	ctxError := getContextualizedError(message, cem.context)

	if errTrace, ok := err.(CtxErrorTrace); ok {
		errTrace.Trace = append([]CtxError{ctxError}, errTrace.Trace...)
		return errTrace
	}

	if _, ok := err.(CtxError); !ok {
		ctxError.ErrorS = err.Error()
		ctxError.ErrorI = err
	}

	return CtxErrorTrace{Trace: []CtxError{ctxError}, StackTrace: string(debug.Stack())}
}

func (cem CtxErrorManager) New(message string) CtxErrorTraceI {
	ctxError := getContextualizedError(message, cem.context)
	return CtxErrorTrace{Trace: []CtxError{ctxError}, StackTrace: string(debug.Stack())}
}

func Wrap(err error, message string) CtxErrorTraceI {

	if err == nil{
		return nil
	}

	ctxError := getContextualizedError(message, nil)

	if errTrace, ok := err.(CtxErrorTrace); ok {
		errTrace.Trace = append([]CtxError{ctxError}, errTrace.Trace...)
		return errTrace
	}

	if _, ok := err.(CtxError); !ok {
		ctxError.ErrorS = err.Error()
		ctxError.ErrorI = err
	}

	return CtxErrorTrace{Trace: []CtxError{ctxError}, StackTrace: string(debug.Stack())}
}


func New(message string) CtxErrorTraceI {
	ctxError := getContextualizedError(message, nil)
	return CtxErrorTrace{Trace: []CtxError{ctxError}, StackTrace: string(debug.Stack())}
}

func getContextualizedError(message string, context map[string]interface{}) CtxError {
	ctxError := CtxError{
		Message: message,
		Context: context,
	}

	functionName, fileName, line, ok := runtime.Caller(2)
	//If we can retrieve the runtime informations, we add them to the error
	if ok {
		ctxError.FunctionName = runtime.FuncForPC(functionName).Name()
		ctxError.FileName = fileName
		ctxError.Line = line
	}

	return ctxError
}

func (contextualizedError CtxError) GetMessage() string {
	return contextualizedError.Message
}

func SetContext(m map[string]interface{}) CtxErrorManager {
	if m == nil {
		return CtxErrorManager{context: m}
	}

	for key := range m {
		for _, hiddenField := range HiddenFields {
			if key == hiddenField {
				m[key] = "hidden"
				break
			}
		}
	}

	return CtxErrorManager{context: m}
}

func (cem CtxErrorManager) GetContext() map[string]interface{} {
	return cem.context
}

func SetHiddenFields(fields ...string)  {
	HiddenFields = append(HiddenFields, fields...)
}
