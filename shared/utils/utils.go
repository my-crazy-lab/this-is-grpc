package utils

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PbTimestampToISO(timestamp *timestamppb.Timestamp) string {
	if timestamp == nil {
		return ""
	}
	// Convert protobuf Timestamp to time.Time
	t := timestamp.AsTime()
	// Format the time.Time object as ISO 8601 string
	return t.Format(time.RFC3339)
}

func CreateMetadataToken(ctx context.Context) (context.Context, error) {
	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return nil, errors.New("unauthorized: missing token")
	}

	// Create gRPC metadata with the token
	md := metadata.New(map[string]string{"user-authorization": token})
	return metadata.NewOutgoingContext(ctx, md), nil
}
