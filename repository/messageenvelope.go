package repository

import (
	"fmt"
	"time"

	"github.com/blackhatbrigade/gomessagestore/uuid"
)

//Actual values that come out of the database
type MessageEnvelope struct {
	ID             uuid.UUID `db:"id"`
	StreamName     string    `db:"stream_name"`
	StreamCategory string    `db:"stream_category"`
	MessageType    string    `db:"type"`
	Version        int64     `db:"position"`
	GlobalPosition int64     `db:"global_position"`
	Data           []byte    `db:"data"`
	Metadata       []byte    `db:"metadata"`
	Time           time.Time `db:"time"`
}

func (msgEnv *MessageEnvelope) String() string {
	return fmt.Sprintf("GlobalPosition: %d | ID: %s | MessageType: %s | StreamName: %s | StreamCategory: %s | Data: %+v| Metadata: %+v", msgEnv.GlobalPosition, msgEnv.ID, msgEnv.MessageType, msgEnv.StreamName, msgEnv.StreamCategory, string(msgEnv.Data), string(msgEnv.Metadata))
}
