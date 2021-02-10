package edns

import (
	"context"
	"encoding/gob"
	"log"
	"net"
	"sync"

	"sigs.k8s.io/external-dns/endpoint"
)

type Handler struct {
	ConnCh  <-chan net.Conn
	StateCh <-chan []endpoint.Endpoint
}

func (h *Handler) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	var state []endpoint.Endpoint

	// Wait for the state to be initialized first.
	select {
	case <-ctx.Done():
		return
	case newState := <-h.StateCh:
		state = newState
	}

	// Then accept state updates and connections.
	for {
		select {
		case <-ctx.Done():
			return
		case newState := <-h.StateCh:
			state = newState
		case conn := <-h.ConnCh:
			go serveTcp(ctx, conn, state)
		}
	}
}

func serveTcp(ctx context.Context, conn net.Conn, state []endpoint.Endpoint) {
	defer conn.Close()
	enc := gob.NewEncoder(conn)
	if err := enc.Encode(state); err != nil {
		log.Printf("error while writing the state to the socket: %v", err)
	}
}
