package errs

// Strip will strip the error from a function call and exit if the error is not nil.
func Strip[T any](t T, err error) T {
	MaybeExit(err)
	return t
}
