package cstx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const metadataKey = "cstx-id"

func GetCstxIDFromGrpcMetadata(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		val := md[metadataKey]
		if len(val) > 0 {
			return val[0]
		}
	}
	return ""
}

func AddCstxIDToGrpcMetadata(ctx context.Context, cstxID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, metadataKey, cstxID)
}
