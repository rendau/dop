package sms

type Sms interface {
	Send(phone string, msg string) bool
	SendAsync(phone string, msg string) bool
}
