package ctxerror

import (
	"fmt"
	"encoding/json"
	"runtime"
	"runtime/debug"
)

type CtxErrorManager struct {
	context map[string]interface{}
}

type CtxErrorTrace struct {
	Trace []CtxError `json:"trace"`
	StackTrace string `json:"stack_trace"`
}

type CtxError struct {
	Message string `json:"message"`
	FileName string `json:"file_name"`
	Line int `json:"line"`
	FunctionName string `json:"function_name"`
	Context map[string]interface{} `json:"data"`
	ErrorI string `json:"error"`
}

func (cet CtxErrorTrace) Error() string{
	ctxErrorTraceBytes, err := json.MarshalIndent(cet , "", "   ")
	if err != nil{
		if len(cet.Trace) > 0{
			return fmt.Sprintf("%s\n%v", cet.Trace[0].Message, ctxErrorTraceBytes)
		}
		return fmt.Sprintf("%v", ctxErrorTraceBytes)
	}

	return string(ctxErrorTraceBytes)
}

func (cet CtxErrorTrace) ErrorJson() string{
	ctxErrorTraceBytes, err := json.MarshalIndent(cet , "", "   ")
	if err != nil{
		if len(cet.Trace) > 0{
			return fmt.Sprintf("%s\n%v", cet.Trace[0].Message, ctxErrorTraceBytes)
		}
		return fmt.Sprintf("%v", ctxErrorTraceBytes)
	}

	return string(ctxErrorTraceBytes)
}

func (cem CtxErrorManager) AddContext(key string, val interface{}) {
	if cem.context == nil{
		cem.context = make(map[string]interface{})
	}

	cem.context[key] = val
}

func (ctxError CtxError) Error() string{
	contextualizedErrorBytes, err := json.MarshalIndent(ctxError , "", "   ")
	if err != nil{
		return fmt.Sprintf("%v", ctxError)
	}

	return string(contextualizedErrorBytes)
}

func (cem CtxErrorManager) Wrap(err error, message string) CtxErrorTrace {

	ctxError := getContextualizedError(message, cem.context)

	if errTrace, ok := err.(CtxErrorTrace); ok {
		errTrace.Trace = append([]CtxError{ctxError}, errTrace.Trace...)
		return errTrace
	}

	if _, ok := err.(CtxError); !ok{
		ctxError.ErrorI = err.Error()
	}

	return CtxErrorTrace{Trace:[]CtxError{ctxError}, StackTrace: string(debug.Stack())}
}

func (cem CtxErrorManager)(message string) CtxErrorTrace{
	ctxError := getContextualizedError(message, cem.context)

	return CtxErrorTrace{Trace:[]CtxError{ctxError}, StackTrace: string(debug.Stack())}
}

func Wrap(err error, message string) CtxErrorTrace {
	ctxError := getContextualizedError(message, nil)

	if errTrace, ok := err.(CtxErrorTrace); ok {
		errTrace.Trace = append([]CtxError{ctxError}, errTrace.Trace...)
		return errTrace
	}

	if _, ok := err.(CtxError); !ok{
		ctxError.ErrorI = err.Error()
	}

	return CtxErrorTrace{Trace:[]CtxError{ctxError}, StackTrace: string(debug.Stack())}
}

func getContextualizedError(message string, context map[string]interface{}) CtxError{
	ctxError := CtxError{
		Message: message,
		Context: context,
	}

	functionName , fileName, line, ok := runtime.Caller(2)
	//If we can retrieve the runtime informations, we add them to the error
	if ok{
		ctxError.FunctionName = runtime.FuncForPC(functionName).Name()
		ctxError.FileName = fileName
		ctxError.Line = line
	}

	return ctxError
}

func (contextualizedError CtxError) GetMessage() string{
	return contextualizedError.Message
}

func SetContext(m map[string]interface{}) CtxErrorManager {
	return CtxErrorManager{context: m}
}

func (cem CtxErrorManager)GetContext() map[string]interface{}{
	return cem.context
}
