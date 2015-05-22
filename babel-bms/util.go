package main

import (
	"time"

	"github.com/twinj/uuid"
)

func CreateUUID() string {
	u := uuid.NewV4()
	return u.String()
}

func ValidateUUID(aUUID string) (uuid.UUID, error) {
	return uuid.ParseUUID(aUUID)
}

func epochToTime(epoch float64) (time.Time, error) {
	epochi := int64(epoch)
	return time.Unix(0, epochi*int64(time.Millisecond)), nil
}
