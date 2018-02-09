package main


type Notification struct {
	ApnsID string

	DeviceToken string

	Topic string
	// A byte array containing the JSON-encoded payload of this push notification.
	// Refer to "The Remote Notification Payload" section in the Apple Local and
	// Remote Notification Programming Guide for more info.
	Payload interface{}
}
