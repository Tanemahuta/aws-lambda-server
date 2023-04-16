package errorx

// Fns is a slice of Fn.
type Fns []Fn

// Run the Fns, returning the first error. If the error is a skip, nil will be returned.
func (fns Fns) Run() error {
	var err error
	for _, fn := range fns {
		if err = fn(); err != nil {
			break
		}
	}
	return err
}
