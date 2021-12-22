package testing

type Scheme struct {
	Address  string
	Username string
	Password string
	Port     int
	Channel  int
}

type DeviceServiceHealth struct {
	ConnectTime	int
	RecvTimeOut	int
	Reconnect	int
	LogToFile	string
	LogLevel	int
}
