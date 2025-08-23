package delivery

import (
	"github.com/Paveluts42/bookreader/backend/api/apiconnect"
)

type Server struct {
	apiconnect.BookServiceHandler
	apiconnect.NoteServiceHandler
	apiconnect.UserServiceHandler
}

func NewServer() *Server {
	return &Server{}
}
