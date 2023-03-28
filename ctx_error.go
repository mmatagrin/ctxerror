package ctxerror

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
)

//By default `password` is hidden
var HiddenFields = []string{
	"password",
	"email",
	"username",
	"mail",
	"first-name",
	"last-name",
	"first_name",
	"last_name",
	"name",
}

type CtxErrorManager struct {
	context map[string]interface{}
}

type CtxErrorTraceI interface {
	Error() string
	ErrorJson() string
	GetMessage() string
	GetTrace() []CtxError
	AddError(error, string) CtxErrorTraceI
	ErrorKind() string
}

type CtxErrorTrace struct {
	Trace []CtxError `json:"trace"`
	kind  *string
	//StackTrace string     `json:"stack_trace"`
}

type CtxError struct {
	Message      string                 `json:"message"`
	FileName     string                 `json:"file_name"`
	Line         int                    `json:"line"`
	FunctionName string                 `json:"function_name"`
	Context      map[string]interface{} `json:"context"`
	ErrorS       string                 `json:"error"`
	ErrorI       error
}

func (cet CtxErrorTrace) ErrorKind() string {
	if cet.kind != nil {
		return *cet.kind
	}
	return ""
}

func (cet CtxErrorTrace) GetMessage() string {
	if cet.Trace != nil && len(cet.Trace) > 0 {
		return cet.Trace[0].GetMessage()
	}

	return ""
}

func (cet CtxErrorTrace) Error() string {
	ctxErrorTraceBytes, err := json.Marshal(cet)
	if err != nil {
		return fmt.Sprintf("%v", cet)
	}

	if len(cet.Trace) > 0 {
		return fmt.Sprintf("%s\n%v", cet.Trace[0].Message, string(ctxErrorTraceBytes))
	}

	return fmt.Sprintf("%v", string(ctxErrorTraceBytes))

}

func (cet CtxErrorTrace) ErrorJson() string {
	ctxErrorTraceBytes, err := json.Marshal(cet)
	if err != nil {
		return fmt.Sprintf("%v", ctxErrorTraceBytes)
	}

	return string(ctxErrorTraceBytes)
}

func (cet CtxErrorTrace) GetTrace() []CtxError {
	if cet.Trace == nil {
		return []CtxError{}
	}

	return cet.Trace
}

func (cet CtxErrorTrace) AddError(err error, message string) CtxErrorTraceI {
	if err == nil {
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

	if valMap, ok := val.(map[string]interface{}); ok {
		val = sanitizeContext(valMap)
	} else {
		for _, hiddenField := range HiddenFields {
			if key == hiddenField {
				val = "hidden"
				break
			}
		}
	}

	cem.context[key] = val
}

func (ctxError CtxError) Error() string {
	contextualizedErrorBytes, err := json.Marshal(ctxError)
	if err != nil {
		return fmt.Sprintf("%v", ctxError)
	}

	return string(contextualizedErrorBytes)
}

func (cem CtxErrorManager) Wrap(err error, message string) CtxErrorTraceI {
	if err == nil {
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

	return CtxErrorTrace{Trace: []CtxError{ctxError} /*, StackTrace: string(debug.Stack())*/}
}

func (cem CtxErrorManager) WrapWithKind(kind string, err error, message string) CtxErrorTraceI {
	if err == nil {
		return nil
	}

	ctxError := getContextualizedError(message, cem.context)

	if errTrace, ok := err.(CtxErrorTrace); ok {
		errTrace.Trace = append([]CtxError{ctxError}, errTrace.Trace...)
		errTrace.kind = &kind
		return errTrace
	}

	if _, ok := err.(CtxError); !ok {
		ctxError.ErrorS = err.Error()
		ctxError.ErrorI = err
	}

	return CtxErrorTrace{Trace: []CtxError{ctxError}, kind: &kind /*, StackTrace: string(debug.Stack())*/}
}

func (cem CtxErrorManager) New(message string) CtxErrorTraceI {
	ctxError := getContextualizedError(message, cem.context)
	return CtxErrorTrace{Trace: []CtxError{ctxError} /*, StackTrace: string(debug.Stack())*/}
}

func (cem CtxErrorManager) NewWithKind(message string, kind string) CtxErrorTraceI {
	ctxError := getContextualizedError(message, cem.context)
	return CtxErrorTrace{Trace: []CtxError{ctxError}, kind: &kind /*, StackTrace: string(debug.Stack())*/}
}

func Wrap(err error, message string) CtxErrorTraceI {
	if err == nil {
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

	return CtxErrorTrace{Trace: []CtxError{ctxError} /*, StackTrace: string(debug.Stack())*/}
}

func WrapWithKind(kind string, err error, message string) CtxErrorTraceI {
	if err == nil {
		return nil
	}

	ctxError := getContextualizedError(message, nil)

	if errTrace, ok := err.(CtxErrorTrace); ok {
		errTrace.Trace = append([]CtxError{ctxError}, errTrace.Trace...)
		errTrace.kind = &kind
		return errTrace
	}

	if _, ok := err.(CtxError); !ok {
		ctxError.ErrorS = err.Error()
		ctxError.ErrorI = err
	}

	return CtxErrorTrace{Trace: []CtxError{ctxError}, kind: &kind /*, StackTrace: string(debug.Stack())*/}
}

func New(message string) CtxErrorTraceI {
	ctxError := getContextualizedError(message, nil)
	return CtxErrorTrace{Trace: []CtxError{ctxError} /*, StackTrace: string(debug.Stack())*/}
}

func NewWithKind(message string, kind string) CtxErrorTraceI {
	ctxError := getContextualizedError(message, nil)
	return CtxErrorTrace{Trace: []CtxError{ctxError}, kind: &kind /*, StackTrace: string(debug.Stack())*/}
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

func sanitizeContext(m map[string]interface{}) (o map[string]interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			o = nil
		}
	}()

	if m == nil {
		return nil
	}

	localCtx := map[string]interface{}{}

	for key, val := range m {
		localCtx[key] = val
		for _, hiddenField := range HiddenFields {
			if key == hiddenField {
				localCtx[key] = "hidden"
				break
			}
		}

		if reflect.TypeOf(val).Kind() == reflect.Map {
			t := reflect.TypeOf(val)

			if t.Key().Kind() != reflect.String {
				continue
			}

			valueOf := reflect.ValueOf(val)
			tmpMap := map[string]interface{}{}

			for _, entry := range valueOf.MapKeys() {
				tmpMap[entry.String()] = valueOf.MapIndex(entry).Interface()
				for _, hiddenField := range HiddenFields {
					if entry.String() == hiddenField {
						tmpMap[entry.String()] = "hidden"
						break
					}
				}
			}
			localCtx[key] = tmpMap
		}
	}

	return localCtx
}

func SetContext(m map[string]interface{}) (c CtxErrorManager) {
	return CtxErrorManager{context: sanitizeContext(m)}
}

func (cem CtxErrorManager) GetContext() map[string]interface{} {
	return cem.context
}

func SetHiddenFields(fields ...string) {
	HiddenFields = fields
}

func AddHiddenFields(fields ...string) {
	HiddenFields = append(HiddenFields, fields...)
}
