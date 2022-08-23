package am

type Subscription interface {
	Unsubscribe() error
}
