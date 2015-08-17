package moitessier

import (
	"github.com/likestripes/kolkata"
	"github.com/likestripes/pacific"
)

const (
	Shared  = 1000
	Private = 1001
	Handler = 1002
)

func Dispatch(context *pacific.Context, person *kolkata.Person, scope, text string, args...interface{}) ([]interface{}, []Listener, error) {

	message := Message{
		Context:        context,
		Scope:          scope,
		Text:           text,
		ActingPersonId: person.PersonId,
	}

	if len(args) > 0 {
		message.Handler, message.Arguments = args[0].(MessageHandlerFunc), args[1:]
	}

	return message.Dispatch()
}

type Conditions struct {
	*pacific.Context
	*kolkata.Person
	ListenerType int
	ScopeString  string
}

func (conditions *Conditions) NewListener() Listener {

	var id string
	scope := conditions.ScopeString

	switch conditions.ListenerType {
	case Private:
	case Handler:
		id = conditions.privateKey()
	case Shared:
		id = conditions.sharedKey()
	}

	listener := Listener{conditions.Context, conditions.Person.PersonId, conditions.Person.PersonId, conditions.ListenerType, id, scope}

	listener.Save()
	listener.AddActingAsSubscriber()

	return listener
}

func (conditions *Conditions) privateKey() string {
	return conditions.PersonIdStr + "/" + conditions.ScopeString
}

func (conditions *Conditions) sharedKey() string {
	return conditions.ScopeString
}
