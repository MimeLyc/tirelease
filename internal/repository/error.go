package repository

type VersionNotFoundError struct {
	Code    int
	Message string
}

func (err VersionNotFoundError) Error() string {
	return err.Message
}
