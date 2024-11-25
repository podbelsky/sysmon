//nolint:godot,godox
package middleware

const (
	TraceIDKey    = "trace_id"
	TraceFlagsKey = "trace_flags"

	TypeKey = "logging/structure"
)

const (
	RPCServerRequestSizeKey  = "rpc.server.request.size"
	RPCServerResponseSizeKey = "rpc.server.response.size"
	RPCServerDurationKey     = "rpc.server.duration"
	ClientAddressKey         = "client.address"

	RPCMethodKey         = "rpc.method"
	RPCGRPCStatusCodeKey = "rpc.grpc.status_code"

	UserAgentOriginalKey = "user_agent.original"
)
