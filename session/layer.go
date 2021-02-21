package session

type ExtendablePersistenceLayer struct {
	LoadSessionFunc   func(id string) (*Session, error)
	SaveSessionFunc   func(id string, session *Session) error
	DeleteSessionFunc func(id string) error
	UserSessionsFunc  func(userID interface{}) ([]string, error)
	LoadUserFunc      func(id interface{}) (User, error)
}

// LoadSession delegates to LoadSessionFunc or returns a nil session.
func (p *ExtendablePersistenceLayer) LoadSession(id string) (*Session, error) {
	if p.LoadSessionFunc != nil {
		return p.LoadSessionFunc(id)
	}
	return nil, nil
}

// SaveSession delegates to SaveSessionFunc or does nothing.
func (p *ExtendablePersistenceLayer) SaveSession(id string, session *Session) error {
	if p.SaveSessionFunc != nil {
		return p.SaveSessionFunc(id, session)
	}
	return nil
}

// DeleteSession delegates to DeleteSessionFunc or does nothing.
func (p *ExtendablePersistenceLayer) DeleteSession(id string) error {
	if p.DeleteSessionFunc != nil {
		return p.DeleteSessionFunc(id)
	}
	return nil
}

// UserSessions delegates to UserSessionsFunc or returns nil.
func (p *ExtendablePersistenceLayer) UserSessions(userID interface{}) ([]string, error) {
	if p.UserSessionsFunc != nil {
		return p.UserSessionsFunc(userID)
	}
	return nil, nil
}

// LoadUser delegates to LoadUserFunc or returns a nil user.
func (p *ExtendablePersistenceLayer) LoadUser(id interface{}) (User, error) {
	if p.LoadUserFunc != nil {
		return p.LoadUserFunc(id)
	}
	return nil, nil
}
