package protorpclog
import (
	"runtime"
	"time"

	"go.pedge.io/proto/time"
	"go.pedge.io/protolog"

	"github.com/golang/protobuf/proto"
)

// Logger is a logger intended to be used as such:
//
// type apiServer struct {
//   protorpclog.Logger
//   ...
// }
//
// func (a *apiServer) Foo(ctx Context.Context, request *FooRequest) (response *FooResponse, err error) {
//   defer func(start time.Now) { a.Log(request, response, err, time.Since(start)) }(time.Now())
//   ...
// }
type Logger interface {
	Log(request proto.Message, response proto.Message, err error, duration time.Duration)
}

// NewLogger returns a new Logger.
func NewLogger(serviceName string) Logger {
	return newLogger(serviceName)
}

type logger struct {
	serviceName string
}

func newLogger(serviceName string) *logger {
	return &logger{serviceName}
}

func (l *logger) Log(request proto.Message, response proto.Message, err error, duration time.Duration) {
	Log(l.serviceName, getMethodName(2), request, response, err, duration)
}

// Log logs an RPC call at the info level if no error, or at the error level if error.
func Log(serviceName string, methodName string, request proto.Message, response proto.Message, err error, duration time.Duration) {
	if err != nil {
		Error(serviceName, methodName, request, response, err, duration)
	} else {
		Info(serviceName, methodName, request, response, err, duration)
	}
}

// Debug logs an RPC call at the debug level.
func Debug(serviceName string, methodName string, request proto.Message, response proto.Message, err error, duration time.Duration) {
	protolog.Debug(event(serviceName, methodName, request, response, err, duration))
}

// Info logs an RPC call at the info level.
func Info(serviceName string, methodName string, request proto.Message, response proto.Message, err error, duration time.Duration) {
	protolog.Info(event(serviceName, methodName, request, response, err, duration))
}

// Warn logs an RPC call at the warn level.
func Warn(serviceName string, methodName string, request proto.Message, response proto.Message, err error, duration time.Duration) {
	protolog.Warn(event(serviceName, methodName, request, response, err, duration))
}

// Error logs an RPC call at the error level.
func Error(serviceName string, methodName string, request proto.Message, response proto.Message, err error, duration time.Duration) {
	protolog.Error(event(serviceName, methodName, request, response, err, duration))
}

func event(serviceName string, methodName string, request proto.Message, response proto.Message, err error, duration time.Duration) *Call {
	call := &Call{
		Service:  serviceName,
		Method:   methodName,
		Duration: prototime.DurationToProto(duration),
	}
	if request != nil {
		call.Request = request.String()
	}
	if response != nil {
		call.Response = response.String()
	}
	if err != nil {
		call.Error = err.Error()
	}
	return call
}

func getMethodName(depth int) string {
	pc := make([]uintptr, 2+depth)
	runtime.Callers(2+depth, pc)
	return runtime.FuncForPC(pc[0]).Name()
}
