package testing

import (
	"context"
)

const (
	Address = "127.0.0.1"
	Username = "root"
	Password = "123456"
)

type Service struct {
	context.Context
	Scheme *Scheme
	Health *DeviceServiceHealth
}

func NewService(ctx context.Context) *Service {
	s := &Service{
		Scheme:  &Scheme{
			Address:  Address,
			Username: Username,
			Password: Password,
			Port:     8000,
			Channel:  1,
		},
	}
	s.Health = &DeviceServiceHealth{
		ConnectTime: 5000, // 毫秒
		RecvTimeOut: 5000,
		Reconnect:   5000,
		LogToFile:   "./record/",
		LogLevel:	 3, // INFO/DEBUG/ERROR
	}

	return s
}
