package moitessier

import (
	"github.com/likestripes/pacific"
	"time"
)

type Listener struct {
	Context        *pacific.Context `datastore:"-" sql:"-" json:"-"`
	ActingPersonId int64            `datastore:"-" sql:"-" json:"-"`
	PersonId       int64
	ListenerType   int
	ListenerId     string
	ScopeString    string
}

type ResponderFunc func([]Message, []interface{}) error
type ListenerFunc func(*Listener, []interface{}) error

func (listener *Listener) Close() {
	listener.RemoveSubscriber()
	subscribers := listener.Subscribers()
	if len(subscribers) == 0 {
		err := listener.query().Delete()
		if err != nil {
			listener.Context.Errorf(err.Error())
		}
	}
}

func (listener *Listener) Save() {
	var existing Listener

	err := listener.query().Get(&existing)
	if existing.ListenerId == listener.ListenerId {
		existing.ActingPersonId = listener.ActingPersonId
		return
	}
	err = listener.query().Put(listener)
	if err != nil {
		listener.Context.Errorf(err.Error())
	}

}

func (listener *Listener) ListenUntil(listener_func ListenerFunc, response_func ResponderFunc, since time.Time, interval_in_milliseconds int, duration_in_milliseconds int, args ...interface{}) error {
	increment := 0
	for {
		err := listener_func(listener, args)
		if err != nil {
			return err
		}

		if results := listener.Messages(since); len(results) > 0 {
			err := response_func(results, args)
			if err != nil {
				return err
			}
		}

		since = time.Now()

		increment = increment + interval_in_milliseconds
		if increment >= duration_in_milliseconds {
			break
		}
		time.Sleep(time.Duration(interval_in_milliseconds) * time.Millisecond)
	}

	return nil
}

func (listener *Listener) ancestor() pacific.Ancestor {
	return pacific.Ancestor{
		Kind:      "Listener",
		KeyString: listener.ListenerId,
	}
}

func (listener *Listener) query() pacific.Query {
	return pacific.Query{
		Context:   *listener.Context,
		Kind:      "Listener",
		KeyString: listener.ListenerId,
	}
}
