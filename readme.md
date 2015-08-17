## Moitessier

A golang library for sending messages & navigating channels.

### Warning

*This is v.01 -- it has no tests, performance is probably awful and it's not proven safe in production anywhere.*  But maybe it'll scratch an itch?

### WTF

`Moitessier` is an opinionated message/dispatch pattern that aims to make it easy to aggregate messages into a listener.


### Install / Import

`go get -u github.com/likestripes/moitessier`

```go
import (
	"github.com/likestripes/moitessier"
)
```

### Dependency on `Pacific`

`Moitessier` uses [likestripes/pacific](https://www.github.com/likestripes/pacific) as an opinionated ORM, and [likestripes/kolkata](https://www.github.com/likestripes/kolkata) to handle user identities. `Kolkata` likely doesn't have to be a hard dependency, there's not a lot special about it relative to `Moitessier` if you're so inclined.

`Pacific` currently supports AppEngine and Postgres:

Google AppEngine: `goapp serve` works out of the box (they include the buildtag for you)

Postgres: `go run -tags 'postgres' main.go` -- details in the [pacific/Readme](https://github.com/likestripes/pacific/blob/master/readme.md).


### Getting Started

#### Listener types

```go

moitessier.Shared // support many subscribers on a listener; only deprecate/delete when the last subscriber has terminated
moitessier.Private // listener's owner's eyes only
moitessier.Handler // up to you & your func

```


#### Invoking a Listener

```go
conditions := moitessier.Conditions{state.Context, src_person, moitessier.Shared, "scoped string"}
listener := conditions.NewListener()

since := time.Now()
interval := 1000 //in ms
time_out := 150000 //in ms
args := ws_state{conn} //OPTIONAL args for the listener_func & responder_fund

listener.ListenUntil(listener_func, responder_func, since, interval, time_out, args)
```

Use your `listener_func` & `responder_func` to control/act in the listener's loop.

```go

func listener_func(listener *moitessier.Listener, args []interface{}) error {
// make sure the listener is still active (i.e. is the websocket connection still open?)
}

func responder_func(msgs []moitessier.Message, args []interface{}) error {
// do something with the message (...post it to the websocket connection)
}
```


#### Dispatching messages

Use `moitessier.Dispatch` to create messages for each `Listener` that's shares the scope and type.  Note that you can use the same `moitessier.Dispacth` for basic messages and custom handlers by including all arguments.


#### `moitessier.Private` or `moitessier.Shared` listeners

```go

moitessier.Dispatch(&state.Context, state.Person, "scoped string", "response string (required for Private and Shared listeners)")

```

#### `moitessier.Handler` listeners

```go

type Handler func(string, Message) interface{}

func saver_func(listener_id string, message moitessier.Message) (interface{}) {

  return "saver was called for: "+listener_id
}

messages, _, _ := moitessier.Dispatch(&state.Context, state.Person, "scoped string", "response string (leave blank if this should be Handler only)", saver_func)
```

### TODO:
- [ ] scope as interface{}, i.e. string, []string, func
- [ ] listeners should be unique to client instead of uniq to person
- [ ] out of band GC
- [ ] tests, docs.
