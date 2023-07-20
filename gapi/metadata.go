package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcgatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extraMetadata(ctx context.Context) *Metadata {
	mtd := &Metadata{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {

		if userAgents := md.Get(grpcgatewayUserAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}
		if clientIps := md.Get(xForwardedForHeader); len(clientIps) > 0 {
			mtd.ClientIp = clientIps[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIp = p.Addr.String()
	}
	return mtd
}
