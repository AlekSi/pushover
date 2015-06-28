package pushover

// DefaultClient is a default client for this package.
var DefaultClient Client

// Send sends message.
// It does not retries failed sends.
// Returns either nil, TemporaryError or FatalError.
func Send(message *Message) error {
	return DefaultClient.Send(message)
}

// SendWithRetries sends message.
// It does retries failed sends for temporary errors up to maxRetries times with 5 second delay.
// Specify maxRetries <= 0 for unlimited retries.
// Returns either nil, TemporaryError (if gave up) or FatalError.
func SendWithRetries(message *Message, maxRetries int) error {
	return DefaultClient.SendWithRetries(message, maxRetries)
}

// SendMessage sends message to specified user.
// It does not retries failed sends.
// Returns either nil, TemporaryError or FatalError.
func SendMessage(user, message string) error {
	return DefaultClient.SendMessage(user, message)
}

// SendMessageWithRetries sends message to specified user.
// It does retries failed sends for temporary errors up to maxRetries times with 5 second delay.
// Specify maxRetries <= 0 for unlimited retries.
// Returns either nil, TemporaryError (if gave up) or FatalError.
func SendMessageWithRetries(user, message string, maxRetries int) error {
	return DefaultClient.SendMessageWithRetries(user, message, maxRetries)
}
