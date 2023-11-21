package ailiyunmsm

type Code interface {
	Send(singerName string, iphonenums ...string) error
}
