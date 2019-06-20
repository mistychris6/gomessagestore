package gomessagestore_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	. "github.com/blackhatbrigade/gomessagestore"
	"github.com/blackhatbrigade/gomessagestore/repository"
	mock_repository "github.com/blackhatbrigade/gomessagestore/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestSubscriberGetsMessages(t *testing.T) {
	messageHandler := &msgHandler{}

	tests := []struct {
		name             string
		expectedError    error
		handlers         []MessageHandler
		expectedPosition int64
		expectedStream   string
		expectedCategory string
		opts             []SubscriberOption
		messageEnvelopes []*repository.MessageEnvelope
		repoReturnError  error
		expectedHandled  []string
		positionEnvelope *repository.MessageEnvelope
	}{{
		name:             "When subscriber is called with SubscribeToEntityStream() option, repository is called correctly",
		expectedStream:   "some category-some id1",
		handlers:         []MessageHandler{messageHandler},
		expectedPosition: 5,
		opts: []SubscriberOption{
			SubscribeToEntityStream("some category", "some id1"),
		},
	}, {
		name:             "When subscriber is called with SubscribeToCategory() option, repository is called correctly",
		expectedCategory: "some category",
		handlers:         []MessageHandler{messageHandler},
		expectedPosition: 5,
		opts: []SubscriberOption{
			SubscribeToCategory("some category"),
		},
	}, {
		name:           "When subscriber is called with SubscribeToEntityStream() option, repository is called correctly",
		handlers:       []MessageHandler{messageHandler},
		expectedStream: "some category:command",
		opts: []SubscriberOption{
			SubscribeToCommandStream("some category"),
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mock_repository.NewMockRepository(ctrl)

			if test.expectedStream != "" {
				mockRepo.
					EXPECT().
					GetAllMessagesInStreamSince(ctx, test.expectedStream, test.expectedPosition, 1000).
					Return(test.messageEnvelopes, test.repoReturnError)
			}
			if test.expectedCategory != "" {
				mockRepo.
					EXPECT().
					GetAllMessagesInCategorySince(ctx, test.expectedCategory, test.expectedPosition, 1000).
					Return(test.messageEnvelopes, test.repoReturnError)
			}

			myMessageStore := NewMessageStoreFromRepository(mockRepo)

			mySubscriber, err := myMessageStore.CreateSubscriber(
				"some id",
				test.handlers,
				test.opts...,
			)

			if err != nil {
				t.Errorf("Failed on CreateSubscriber() Got: %s\n", err)
				return
			}

			_, err = mySubscriber.GetMessages(ctx, test.expectedPosition)
			if err != test.expectedError {
				t.Errorf("Failed to get expected error from GetMessages()\nExpected: %s\n and got: %s\n", test.expectedError, err)
			}
		})
	}
}

func TestSubscriberProcessesMessages(t *testing.T) {

	tests := []struct {
		name             string
		subscriberID     string
		expectedError    error
		handlers         []MessageHandler
		expectedStream   string
		expectedCategory string
		opts             []SubscriberOption
		messages         []Message
		repoReturnError  error
		expectedHandled  []string
		positionEnvelope *repository.MessageEnvelope
	}{{
		name: "Subscriber Poll processes a message in the registered handler with command stream",
		handlers: []MessageHandler{
			&msgHandler{class: "Command MessageType 1"},
			&msgHandler{class: "Command MessageType 2"},
		},
		expectedHandled: []string{
			"Command MessageType 1",
			"Command MessageType 2",
		},
		expectedStream: "category:command",
		opts: []SubscriberOption{
			SubscribeToCommandStream("category"),
		},
		messages: commandsToMessageSlice(getSampleCommands()),
	}, {
		name: "Subscriber Poll processes a message in the registered handler with entity stream",
		handlers: []MessageHandler{
			&msgHandler{class: "Event MessageType 1"},
			&msgHandler{class: "Event MessageType 2"},
		},
		expectedHandled: []string{
			"Event MessageType 1",
			"Event MessageType 2",
		},
		expectedStream: "category-someid",
		opts: []SubscriberOption{
			SubscribeToEntityStream("category", "someid"),
		},
		messages: eventsToMessageSlice(getSampleEvents()),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mock_repository.NewMockRepository(ctrl)

			myMessageStore := NewMessageStoreFromRepository(mockRepo)

			mySubscriber, err := myMessageStore.CreateSubscriber(
				"some id",
				test.handlers,
				test.opts...,
			)

			if err != nil {
				t.Errorf("Failed on CreateSubscriber() Got: %s\n", err)
				return
			}

			_, _, err = mySubscriber.ProcessMessages(ctx, test.messages)
			if err != test.expectedError {
				t.Errorf("Failed to get expected error from ProcessMessages()\nExpected: %s\n and got: %s\n", test.expectedError, err)
			}

			handled := make([]string, 0, len(test.expectedHandled))
			for _, handlerI := range test.handlers {
				handler := handlerI.(*msgHandler)
				if !handler.Called {
					t.Error("Handler was not called")
				}
				handled = append(handled, handler.Handled...) // cause variable names are hard
			}
			if !reflect.DeepEqual(handled, test.expectedHandled) {
				t.Errorf("Handler was called for the wrong messages, \nCalled: %s\nExpected: %s\n", handled, test.expectedHandled)
			}
		})
	}
}

func TestSubscriberGetsPosition(t *testing.T) {

	tests := []struct {
		name             string
		subscriberID     string
		expectedError    error
		handlers         []MessageHandler
		expectedPosition int64
		expectedStream   string
		expectedCategory string
		opts             []SubscriberOption
		messages         []Message
		repoReturnError  error
		expectedHandled  []string
		positionEnvelope *repository.MessageEnvelope
	}{{
		name:             "When GetPosition is called subscriber returns a position that matches the expected position",
		expectedPosition: 0,
		handlers:         []MessageHandler{&msgHandler{}},
		subscriberID:     "some id",
		opts: []SubscriberOption{
			SubscribeToEntityStream("some category", "1234"),
			SubscribeBatchSize(1),
		},
	}, {
		name:             "When GetPosition is called subscriber returns a position that matches the expected position",
		expectedPosition: 400,
		handlers:         []MessageHandler{&msgHandler{}},
		subscriberID:     "some id",
		opts: []SubscriberOption{
			SubscribeToEntityStream("some category", "1234"),
			SubscribeBatchSize(1),
		},
		positionEnvelope: &repository.MessageEnvelope{
			ID:             "some-id-goes-here",
			StreamName:     "I_am_subscriber_id+position",
			StreamCategory: "I_am_subscriber_id+position",
			MessageType:    "CommittedPosition",
			Version:        5,
			GlobalPosition: 500,
			Data:           []byte("{\"position\":400}"),
			Time:           time.Unix(1, 5),
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			mockRepo := mock_repository.NewMockRepository(ctrl)

			mockRepo.
				EXPECT().
				GetLastMessageInStream(ctx, "some id+position").
				Return(&repository.MessageEnvelope{}, nil)

			myMessageStore := NewMessageStoreFromRepository(mockRepo)

			mySubscriber, err := myMessageStore.CreateSubscriber(
				test.subscriberID,
				test.handlers,
				test.opts...,
			)

			if err != nil {
				t.Errorf("Failed on CreateSubscriber() Got: %s\n", err)
				return
			}

			pos, err := mySubscriber.GetPosition(ctx)

			if err != nil {
				t.Errorf("Failed on GetPosition() because of %v", err)
			}

			if pos != test.expectedPosition {
				t.Errorf("Failed on GetPosition()\n Expected%d\n Got: %d", test.expectedPosition, pos)
			}
		})
	}
}
