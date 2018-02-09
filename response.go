package main


type APNSError struct {
	Reason    string
}

func (e *APNSError) Error() string {
	return e.Reason
}

type Response struct {
	StatusCode     int
	ApnsID string
}