package moitessier

import (
	"github.com/likestripes/pacific"
	"time"
)

type Subscriber struct {
	Context      *pacific.Context `datastore:"-" sql:"-" json:"-"`
	Parent       pacific.Ancestor `datastore:"-" sql:"-" json:"-"`
	SubscriberId int64
	PersonId     int64
	ListenerId   string
	Created      time.Time
}

func (listener *Listener) AddActingAsSubscriber() (Subscriber, error) {

	subscriber := Subscriber{
		Context:      listener.Context,
		Parent:       listener.ancestor(),
		SubscriberId: listener.ActingPersonId,
		PersonId:     listener.ActingPersonId,
		ListenerId:   listener.ListenerId,
		Created:      time.Now(),
	}

	query := pacific.Query{
		Context:   *listener.Context,
		Kind:      "Subscriber",
		Ancestors: []pacific.Ancestor{subscriber.Parent},
		KeyInt:    subscriber.SubscriberId,
	}
	if listener.ListenerType == Shared {
		err := query.Get(&subscriber)
		if err != nil {
			err := query.Put(subscriber)
			if err != nil {
				listener.Context.Errorf(err.Error())
			}
		}
		return subscriber, err
	}

	return subscriber, nil
}

func (listener *Listener) RemoveSubscriber() error {
	if listener.ListenerType == Shared {
		query := pacific.Query{
			Context:   *listener.Context,
			Kind:      "Subscriber",
			Ancestors: []pacific.Ancestor{listener.ancestor()},
			KeyInt:    listener.ActingPersonId,
		}
		err := query.Delete()
		if err != nil {
			listener.Context.Errorf(err.Error())
		}
		return err
	}

	return nil
}

func (listener *Listener) Subscribers() []Subscriber {
	var subscribers []Subscriber

	if listener.ListenerType == Shared {
		query := pacific.Query{
			Context:   *listener.Context,
			Kind:      "Subscriber",
			Ancestors: []pacific.Ancestor{listener.ancestor()},
		}

		err := query.GetAll(&subscribers)
		if err != nil {
			listener.Context.Errorf(err.Error())
		}

	} else {
		subscriber := Subscriber{
			Context:      listener.Context,
			SubscriberId: 0,
			PersonId:     listener.PersonId,
			ListenerId:   listener.ListenerId,
		}
		subscribers = append(subscribers, subscriber)
	}

	return subscribers
}
