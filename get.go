package gomessagestore

import (
	"context"
	"fmt"

	"github.com/blackhatbrigade/gomessagestore/repository"
	"github.com/sirupsen/logrus"
)

type getOpts struct {
	stream        *string
	category      *string
	sincePosition bool
	sinceVersion  bool
	since         *int64
	converters    []MessageConverter
	batchsize     int
}

//GetOption provide optional arguments to the Get function
type GetOption func(g *getOpts)

func checkGetOptions(opts ...GetOption) *getOpts {
	g := &getOpts{batchsize: 1000}
	for _, option := range opts {
		option(g)
	}
	return g
}

//Get Gets one or more Messages from the message store.
func (ms *msgStore) Get(ctx context.Context, opts ...GetOption) ([]Message, error) {

	if len(opts) == 0 {
		return nil, ErrMissingGetOptions
	}

	getOptions := checkGetOptions(opts...)
	var msgEnvelopes []*repository.MessageEnvelope
	var err error

	if getOptions.stream != nil && getOptions.category != nil {
		return nil, ErrGetMessagesCannotUseBothStreamAndCategory
	} else if getOptions.stream == nil && getOptions.category == nil {
		return nil, ErrGetMessagesRequiresEitherStreamOrCategory
	}

	if getOptions.since != nil {
		if getOptions.stream != nil {
			msgEnvelopes, err = ms.repo.GetAllMessagesInStreamSince(ctx, *getOptions.stream, *getOptions.since, getOptions.batchsize)
		} else {
			msgEnvelopes, err = ms.repo.GetAllMessagesInCategorySince(ctx, *getOptions.category, *getOptions.since, getOptions.batchsize)
		}
	} else {

		if getOptions.stream != nil {
			msgEnvelopes, err = ms.repo.GetAllMessagesInStream(ctx, *getOptions.stream, getOptions.batchsize)
		}

		if getOptions.category != nil {
			msgEnvelopes, err = ms.repo.GetAllMessagesInCategory(ctx, *getOptions.category, getOptions.batchsize)
		}
	}

	if err != nil {
		logrus.WithError(err).Error("Get: Error getting message")

		return nil, err
	}

	return MsgEnvelopesToMessages(msgEnvelopes, getOptions.converters...), nil
}

//CommandStream allows for writing messages using an expected position
func CommandStream(category string) GetOption {
	return func(g *getOpts) {
		stream := fmt.Sprintf("%s:command", category)
		g.stream = &stream
	}
}

//EventStream allows for getting events in a specific stream
func EventStream(category, entityID string) GetOption {
	return func(g *getOpts) {
		stream := fmt.Sprintf("%s-%s", category, entityID)
		g.stream = &stream
	}
}

//Category allows for getting messages by category
func Category(category string) GetOption {
	return func(g *getOpts) {
		g.category = &category
	}
}

//SincePosition allows for getting only more recent messages
func SincePosition(position int64) GetOption {
	return func(g *getOpts) {
		g.since = &position
		g.sincePosition = true
	}
}

//SinceVersion allows for getting only more recent messages
func SinceVersion(version int64) GetOption {
	return func(g *getOpts) {
		g.since = &version
		g.sinceVersion = true
	}
}

//Converter allows for automatic converting of non-Command/Event type messages
func Converter(converter MessageConverter) GetOption {
	return func(g *getOpts) {
		g.converters = append(g.converters, converter)
	}
}

//BatchSize changes how many messages are returned (default 1000)
func BatchSize(batchsize int) GetOption {
	return func(g *getOpts) {
		g.batchsize = batchsize
	}
}
