package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/MOZGIII/external-dns-endpoint-server/edns"
	"github.com/MOZGIII/external-dns-endpoint-server/logic"
	"github.com/MOZGIII/external-dns-endpoint-server/update"

	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/external-dns/endpoint"
)

func Run() error {
	ctx := signals.SetupSignalHandler()

	addr := os.Getenv("ADDR")

	ednslc := net.ListenConfig{}
	ednsListener, err := ednslc.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("unable to listen for connector: %v", err)
	}
	defer ednsListener.Close()
	connCh := acceptLoop(ednsListener)

	ipChan := make(chan net.IP, 32)
	updateHandler := &update.Handler{
		IPChan: ipChan,
	}

	endpointsChan := make(chan []endpoint.Endpoint, 32)
	logic := logic.Logic{
		IPChan:       ipChan,
		EnpointsChan: endpointsChan,
	}
	go logic.Run(ctx)

	handler := edns.Handler{
		ConnCh:  connCh,
		StateCh: endpointsChan,
	}
	go handler.Run(ctx)

	return http.ListenAndServe(addr, updateHandler)
}

func acceptLoop(listener net.Listener) <-chan net.Conn {
	ch := make(chan net.Conn, 32)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("got error from accept: %v", err)
				return
			}
			ch <- conn
		}
	}()
	return ch
}
