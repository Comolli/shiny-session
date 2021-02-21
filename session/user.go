package session

type User interface {
	// GetID returns the user's unique ID.
	GetID() interface{}
}
