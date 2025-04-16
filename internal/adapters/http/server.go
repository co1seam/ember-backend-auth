package http

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	Server *fiber.App
}

func NewServer() *Server {
	return &Server{
		Server: fiber.New(fiber.Config{
			Prefork:     false,
			AppName:     "Auth",
			JSONDecoder: json.Unmarshal,
			JSONEncoder: json.Marshal,
		})}
}

func (s *Server) Run(port string) error {
	return s.Server.Listen(":" + port)
}

func (s *Server) Shutdown() error {
	return s.Server.Shutdown()
}
