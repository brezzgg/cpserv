package service

import (
	context "context"
	"fmt"

	"github.com/brezzgg/cpserv/clipboard"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(clip clipboard.Clipboard) *Server {
	return &Server{
		clipboard: clip,
	}
}

type Server struct {
	UnimplementedClipboardServiceServer
	clipboard clipboard.Clipboard
}

func (s *Server) Read(context.Context, *Auth) (*Clipboard, error) {
	text, err := s.clipboard.Read()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to read clipboard")
	}
	fmt.Printf("read len=%d\n", len(text))
	return &Clipboard{
		Text: text,
	}, nil
}

func (s *Server) Write(ctx context.Context, req *WriteReq) (*Empty, error) {
	err := s.clipboard.Write(req.Clipboard.Text)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to write clipboard")
	}
	fmt.Printf("written len=%d\n", len(req.Clipboard.Text))
	return nil, nil
}
