package utils

import (
	"time"

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
