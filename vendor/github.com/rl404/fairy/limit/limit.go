package limit

// Limiter is rate limiting interface.
//
// See usage example in example folder.
type Limiter interface {
	// To add count to rate limiter.
	// Will block until the rate in below the
	// limit.
	Take()
}
