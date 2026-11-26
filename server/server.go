package server

import "context"

type Server struct{}

var _ StrictServerInterface = Server{}

func (Server) GetHealth(context.Context, GetHealthRequestObject) (GetHealthResponseObject, error) {
	return GetHealth200JSONResponse{Status: Ok}, nil
}

func NewServer() ServerInterface {
	return NewStrictHandler(Server{}, nil)
}
