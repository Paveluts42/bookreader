package delivery

import (
	"github.com/Paveluts42/bookreader/backend/api/apiconnect"
)

type Server struct {
	apiconnect.UnimplementedReaderServiceHandler
}

func NewServer() *Server {
	return &Server{}
}
