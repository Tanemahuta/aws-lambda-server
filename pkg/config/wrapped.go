package config

// Wrapped interface.
type Wrapped interface {
	// Unwrap the wrapped type.
	Unwrap() interface{}
}

// Unwrap a value.
func Unwrap(src interface{}) interface{} {
	for wrapped, ok := src.(Wrapped); ok; wrapped, ok = src.(Wrapped) {
		src = wrapped.Unwrap()
	}
	return src
}
