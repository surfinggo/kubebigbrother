package stores

type ErrNotFound struct {
}

func (e *ErrNotFound) Error() string {
	return "not found"
}

func IsNotFound(err error) bool {
	if _, ok := err.(*ErrNotFound); ok {
		return true
	}
	return false
}

func WrapNotFound(err error) (*ErrNotFound, bool) {
	if e, ok := err.(*ErrNotFound); ok {
		return e, true
	}
	return nil, false
}
