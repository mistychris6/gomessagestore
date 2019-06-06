package repository

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (r postgresRepo) WriteMessage(ctx context.Context, message *MessageEnvelope) error {
	return r.writeMessageEitherWay(ctx, message)
}

func (r postgresRepo) WriteMessageWithExpectedPosition(ctx context.Context, message *MessageEnvelope, position int64) error {
	return r.writeMessageEitherWay(ctx, message, position)
}

func (r postgresRepo) writeMessageEitherWay(ctx context.Context, message *MessageEnvelope, position ...int64) error {
	if message == nil {
		return ErrNilMessage
	}

	if message.MessageID == "" {
		return ErrMessageNoID
	}

	if message.Stream == "" {
		return ErrInvalidStreamID
	}

	// our return channel for our goroutine that will either finish or be cancelled
	retChan := make(chan error, 1)
	go func() {
		// last thing we do is ensure our return channel is populated
		defer func() {
			retChan <- nil
		}()

		eventideMetadata := &eventideMessageMetadata{
			CorrelationID: message.CorrelationID,
			CausedByID:    message.CausedByID,
			UserID:        message.UserID,
		}
		eventideMessage := &eventideMessageEnvelope{
			ID:          message.MessageID,
			MessageType: message.Type,
			StreamName:  message.Stream,
			Data:        message.Data,
			Position:    message.Position,
		}

		if metadata, err := json.Marshal(eventideMetadata); err == nil {
			eventideMessage.Metadata = metadata
		} else {
			logrus.WithError(err).Error("Failure to marshal metadata in repo_postgres.go::WriteMessage")
			retChan <- err
			return
		}

		/*"write_message(
			_id varchar,
			_stream_name varchar,
			_type varchar,
			_data jsonb,
			_metadata jsonb DEFAULT NULL,
			_expected_version bigint DEFAULT NULL
		)"*/
		if len(position) > 0 {
			if position[0] < -1 {
				retChan <- ErrInvalidPosition
				return
			}

			// with _expected_version passed in
			query := "SELECT write_message($1, $2, $3, $4, $5, $6)"
			if _, err := r.dbx.ExecContext(ctx, query, eventideMessage.ID, eventideMessage.StreamName, eventideMessage.MessageType, eventideMessage.Data, eventideMessage.Metadata, position[0]); err != nil {
				logrus.WithError(err).Error("Failure in repo_postgres.go::WriteMessageWithExpectedPosition")
				retChan <- err
				return
			}
		} else {
			// without _expected_version passed in
			query := "SELECT write_message($1, $2, $3, $4, $5)"
			logrus.WithFields(logrus.Fields{
				"query":                      query,
				"eventideMessageID":          eventideMessage.ID,
				"eventideMessageStreamName":  eventideMessage.StreamName,
				"eventideMessageMessageType": eventideMessage.MessageType,
				"eventideMessageData":        eventideMessage.Data,
				"eventideMessageMetadata":    eventideMessage.Metadata,
			}).Debug("about to write message")
			if _, err := r.dbx.ExecContext(ctx, query, eventideMessage.ID, eventideMessage.StreamName, eventideMessage.MessageType, eventideMessage.Data, eventideMessage.Metadata); err != nil {
				logrus.WithError(err).Error("Failure in repo_postgres.go::WriteMessage")
				retChan <- err
				return
			}
		}
	}()

	// wait for our return channel or the context to cancel
	select {
	case retval := <-retChan:
		return retval
	case <-ctx.Done():
		return nil
	}
}