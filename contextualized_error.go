package contextualized_error

import (
	"fmt"
	"encoding/json"
	"runtime"
	)

type ContextualizedErrorManager struct {
	context map[string]interface{}
}

type ContextualizedError struct {
	Message string `json:"message"`
	FileName string `json:"file_name"`
	Line int `json:"line"`
	FunctionName string `json:"function_name"`
	Data map[string]interface{} `json:"data"`
	Previous *ContextualizedError `json:"previous"`
	ErrorMsg string `json:"error"`
}


func (contextualizedError ContextualizedError) Error() string{
	contextualizedErrorBytes, err := json.MarshalIndent(contextualizedError , "", "   ")
	if err != nil{
		return fmt.Sprintf("%v", contextualizedError)
	}

	return string(contextualizedErrorBytes)
}

func (cem ContextualizedErrorManager) NewContextualizedError(err interface{}, message string) ContextualizedError{

	contextualizedError := ContextualizedError{
		Data: cem.context,
		Message: message,
	}

	functionName , fileName, line, ok := runtime.Caller(1)
	//If we can retrieve the runtime informations, we add them to the error
	if ok{
		contextualizedError.FunctionName = runtime.FuncForPC(functionName).Name()
		contextualizedError.FileName = fileName
		contextualizedError.Line = line
	}

	if previous, ok := err.(ContextualizedError); ok{
		contextualizedError.Previous = &previous
	} else if previous, ok := err.(*ContextualizedError); ok{
		contextualizedError.Previous = previous
	} else {
		contextualizedError.ErrorMsg = fmt.Sprintf("%v", err)
	}

	return contextualizedError
}


func (contextualizedError ContextualizedError) GetMessage() string{
	return contextualizedError.Message
}


func SetErrorContext(m map[string]interface{}) ContextualizedErrorManager{
	return ContextualizedErrorManager{context: m}
}

func (cem ContextualizedErrorManager)GetContext() map[string]interface{}{
	return cem.context
}