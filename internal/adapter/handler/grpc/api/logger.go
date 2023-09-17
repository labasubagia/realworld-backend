package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/labasubagia/realworld-backend/internal/core/domain"
	"github.com/labasubagia/realworld-backend/internal/core/port"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	reqIDHeader     = "request-id"
	userAgentHeader = "user-agent"
)

func (server *Server) Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	reqID := domain.NewID().String()
	var userAgent string
	metaData, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if reqIDs := metaData.Get(reqIDHeader); len(reqIDs) > 0 {
			userAgent = reqIDs[0]
		}
		if userAgents := metaData.Get(userAgentHeader); len(userAgents) > 0 {
			userAgent = userAgents[0]
		}
	}

	peerInfo, _ := peer.FromContext(ctx)
	clientIP := peerInfo.Addr.String()

	logger := server.logger.NewInstance().Field("request_id", reqID).Logger()
	ctx = context.WithValue(ctx, port.SubLoggerCtxKey, logger)

	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	code := codes.Unknown
	if st, ok := status.FromError(err); ok {
		code = st.Code()
	}

	logEvent := logger.Info()
	if code == codes.Internal {
		logEvent = logger.Error()
		bytes, _ := json.Marshal(req)
		logEvent.Field("request", string(bytes))
	}

	logEvent.
		Field("protocol", "grpc").
		Field("method", info.FullMethod).
		Field("client_ip", clientIP).
		Field("user_agent", userAgent).
		Field("status_code", int(code)).
		Field("status_text", code.String()).
		Field("duration", duration).
		Msg("receive grpc request")

	return result, err
}
