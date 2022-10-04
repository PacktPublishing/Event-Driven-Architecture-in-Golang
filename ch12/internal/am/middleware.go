package am

// MessageStreamWithMiddleware builds a applyMiddleware chain around a stream
//
// Middleware are applied in reverse; this makes the first applyMiddleware
// in the slice the outermost i.e. first to enter, last to exit
// given: stream, A, B, C
// result: A(B(C(stream)))
func MessageStreamWithMiddleware(stream MessageStream, mws ...MessageStreamMiddleware) MessageStream {
	return applyMiddleware(stream, mws...)
}

// MessagePublisherWithMiddleware builds a applyMiddleware chain around a publisher
//
// Middleware are applied in reverse; this makes the first applyMiddleware
// in the slice the outermost i.e. first to enter, last to exit
// given: publisher, A, B, C
// result: A(B(C(publisher)))
func MessagePublisherWithMiddleware(publisher MessagePublisher, mws ...MessagePublisherMiddleware) MessagePublisher {
	return applyMiddleware(publisher, mws...)
}

// MessageHandlerWithMiddleware builds a applyMiddleware chain around a handler
//
// Middleware are applied in reverse; this makes the first applyMiddleware
// in the slice the outermost i.e. first to enter, last to exit
// given: handler, A, B, C
// result: A(B(C(handler)))
func MessageHandlerWithMiddleware(handler MessageHandler, mws ...MessageHandlerMiddleware) MessageHandler {
	return applyMiddleware(handler, mws...)
}

func applyMiddleware[T any, M func(T) T](target T, mws ...M) T {
	h := target
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
