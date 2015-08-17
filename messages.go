package moitessier

import (
	"github.com/likestripes/pacific"
	"math/rand"
	"time"
)

type MessageHandlerFunc func(string, Message) string

type Message struct {
	Context        *pacific.Context `datastore:"-" sql:"-" json:"-"`
	Listener       Listener         `datastore:"-" sql:"-" json:"-"`
	ActingPersonId int64            `datastore:"-" sql:"-" json:"-"`
	ListenerId     string
	Scope          string
	Created        time.Time
	PersonId       int64
	MessageId      int64
	Text           string
	Handler        MessageHandlerFunc `datastore:"-" sql:"-" json:"-"`
	Arguments      []interface{} `datastore:"-" sql:"-" json:"-"`
}

func (message Message) Dispatch() ([]interface{}, []Listener, error) {

	query := pacific.Query{
		Context: *message.Context,
		Kind:    "Listener",
		Filters: map[string]interface{}{
			"ScopeString =": message.Scope,
		},
	}

	var listeners []Listener
	err := query.GetAll(&listeners)

	var results []interface{}
	for _, listener := range listeners {
		var result interface{}
		message.setListener(listener)
		if listener.ListenerType == Handler && message.Handler != nil {
			message.Text = message.Handler(listener.ListenerId, message)
		}
		if message.Text != "" {
			message.Save()
			result = message
			results = append(results, result)
		}
	}

	return results, listeners, err
}

func (message Message) Save() {

	message.MessageId = rand.Int63()
	message.Created = time.Now()

	query := pacific.Query{
		Context:   *message.Context,
		Kind:      "Message",
		Ancestors: []pacific.Ancestor{message.Listener.ancestor()},
		KeyInt:    message.MessageId,
	}
	err := query.Put(&message)
	if err != nil {
		message.Context.Errorf(err.Error())
	}

}

func (listener *Listener) Messages(since time.Time) []Message {

	query := pacific.Query{
		Context:   *listener.Context,
		Kind:      "Message",
		Ancestors: []pacific.Ancestor{listener.ancestor()},
		Filters: map[string]interface{}{
			"Created >": since,
		},
	}

	if listener.ListenerType == Private {
		query.Filters["PersonId ="] = listener.PersonId
	}

	var messages []Message
	err := query.GetAll(&messages)
	if err != nil {
		listener.Context.Errorf(err.Error())
	}
	return messages
}

func (message *Message) setListener(listener Listener) {
	message.Listener = listener
	message.ListenerId = listener.ListenerId
	message.PersonId = message.ActingPersonId
}
