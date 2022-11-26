package cstx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const metadataKey = "cstx-serialized"

func GetSerializedCstxFromGrpcMetadata(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		val := md[metadataKey]
		if len(val) > 0 {
			return val[0]
		}
	}
	return ""
}

func AddSerializedCstxToGrpcMetadata(ctx context.Context, cstx string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, metadataKey, cstx)
}
