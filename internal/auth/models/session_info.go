package models

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

type SessionInfo struct {
	UserId    int
	Token     string
	IPAddress string
	Expire    time.Time
}

func (s SessionInfo) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(s) 
}

func (s *SessionInfo) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, s) 
}
