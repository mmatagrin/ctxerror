package contextualized_error

import (
	"fmt"
	"encoding/json"
	"runtime"
	)

type CtxErrorManager struct {
	context map[string]interface{}
}

type ContextualizedError struct {
	Message string `json:"message"`
	FileName string `json:"file_name"`
	Line int `json:"line"`
	FunctionName string `json:"function_name"`
	Data map[string]interface{} `json:"data"`
	ErrorI error `json:"error"`
}


func (contextualizedError ContextualizedError) Error() string{
	contextualizedErrorBytes, err := json.MarshalIndent(contextualizedError , "", "   ")
	if err != nil{
		return fmt.Sprintf("%v", contextualizedError)
	}

	return string(contextualizedErrorBytes)
}

func (cem CtxErrorManager) Wrap(err error, message string) ContextualizedError{

	contextualizedError := ContextualizedError{
		Data: cem.context,
		Message: message,
		ErrorI: err,
	}

	functionName , fileName, line, ok := runtime.Caller(1)
	//If we can retrieve the runtime informations, we add them to the error
	if ok{
		contextualizedError.FunctionName = runtime.FuncForPC(functionName).Name()
		contextualizedError.FileName = fileName
		contextualizedError.Line = line
	}

	return contextualizedError
}


func (contextualizedError ContextualizedError) GetMessage() string{
	return contextualizedError.Message
}


func SetContext(m map[string]interface{}) CtxErrorManager {
	return CtxErrorManager{context: m}
}

func (cem CtxErrorManager)GetContext() map[string]interface{}{
	return cem.context
}