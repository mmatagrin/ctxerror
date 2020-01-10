package contextualized_error

import (
	"fmt"
	"encoding/json"
	"runtime"
	)

type CtxErrorManager struct {
	context map[string]interface{}
}

type CtxErrorTrace struct {
	Trace []CtxError `json:"trace"`
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
		return fmt.Sprintf("%v", ctxErrorTraceBytes)
	}

	return string(ctxErrorTraceBytes)
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

	return CtxErrorTrace{Trace:[]CtxError{ctxError}}
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

	return CtxErrorTrace{Trace:[]CtxError{ctxError}}
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
