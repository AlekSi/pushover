package pushover

// DefaultClient is a default client for this package.
var DefaultClient Client

// Send sends message.
// Returns either TemporaryError or FatalError.
func Send(message *Message) error {
	return DefaultClient.Send(message)
}

// SendMessage sends message to specified user.
// Returns either TemporaryError or FatalError.
func SendMessage(user, message string) error {
	return DefaultClient.SendMessage(user, message)
}
