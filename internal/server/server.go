package server

import (
	http_v1 "cw/internal/controller/http/v1"
	"cw/internal/controller/websocket"
	"cw/internal/service"

	"go.uber.org/zap"
)

func MustRun() error {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	service := &service.Service{}

	hub := websocket.NewHub(zapLogger)
	go hub.Run()

	handler := http_v1.NewHandler(zapLogger, service, hub)

	app := handler.InitRoutes()

	return app.Listen(":8080")
}
