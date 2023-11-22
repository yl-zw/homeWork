package ailiyunmsm

type Code interface {
	Send(singerName, code string, phoneNumber ...string) error
}
