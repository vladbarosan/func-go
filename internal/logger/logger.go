package logger

// separate package in order to have private eventStream not visible from azure package

import (
	"fmt"

	"github.com/Azure/azure-functions-go-worker/internal/rpc"
)

// Logger exposes the functionality to send logs back to the runtime
type Logger struct {
	invocationID string
	eventStream  rpc.FunctionRpc_EventStreamClient
}

// NewLogger returns a new instance of type Logger to be used in user funcs
func NewLogger(e rpc.FunctionRpc_EventStreamClient, invocationID string) *Logger {
	return &Logger{
		invocationID: invocationID,
		eventStream:  e,
	}
}

// Log sends a log message to the runtime
func (l *Logger) Log(format string, args ...interface{}) error {

	log := &rpc.RpcLog{
		InvocationId: l.invocationID,
		Level:        rpc.RpcLog_Information,
		Message:      fmt.Sprintf(format, args...),
	}

	return l.eventStream.Send(&rpc.StreamingMessage{
		Content: &rpc.StreamingMessage_RpcLog{
			RpcLog: log,
		},
	})
}
