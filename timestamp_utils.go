package gproxy

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"time"
)

func MustTimestampOrNow(ts *timestamp.Timestamp) time.Time {
	if ts == nil {
		return time.Now().UTC()
	}
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		log.Print("MustTimestampOrNow: ", err)
	}
	return t
}
