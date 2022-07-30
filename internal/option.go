package internal

type Options struct {
	AppName           string
	ListenAddressHTTP string
	DatabaseHost      string
	DatabaseUser      string
	DatabasePass      string
	DatabaseName      string
	DatabasePort      int
	DatabaseSSLMode   string
	DatabaseTimezone  string
}
