package gomessagestore

import (
	"context"

	"github.com/sirupsen/logrus"
)

//ProcessMessages uses the handlers of the subscriptionWorker to process the messages retrieved from the message store; third process of the polling loop
func (sw *subscriptionWorker) ProcessMessages(ctx context.Context, msgs []Message) (messagesHandled int, positionOfLastHandled int64, err error) {
	log := logrus.WithFields(logrus.Fields{
		"subscriberID": sw.subscriberID,
	})

	for _, msg := range msgs {
		for _, handler := range sw.handlers {
			if handler.Type() == msg.Type() {
				err = handler.Process(ctx, msg)
				if err != nil {
					log.WithError(err).Error("A handler failed to process a message not moving on")
					return
				}

				if !sw.config.stream {
					// category subscriptions care about position
					positionOfLastHandled = msg.Position()
				} else {
					// stream subscriptions care about version
					positionOfLastHandled = msg.Version()
				}
				messagesHandled++
			}
		}
	}
	return
}
