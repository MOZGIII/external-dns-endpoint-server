package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/MOZGIII/external-dns-endpoint-server/internal/edns"
	"github.com/MOZGIII/external-dns-endpoint-server/internal/httprun"
	"github.com/MOZGIII/external-dns-endpoint-server/internal/logic"
	"github.com/MOZGIII/external-dns-endpoint-server/internal/update"

	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/external-dns/endpoint"
)

func Run() error {
	ctx := signals.SetupSignalHandler()

	httpAddr := readEnv("HTTP_ADDR")
	ednsAddr := readEnv("EDNS_ADDR")

	ednslc := net.ListenConfig{}

	ednsListener, err := ednslc.Listen(ctx, "tcp", ednsAddr)
	if err != nil {
		return fmt.Errorf("unable to listen for connector: %w", err)
	}
	defer ednsListener.Close()
	connCh := acceptLoop(ednsListener)

	ipChan := make(chan net.IP)
	updateHandler := &update.Handler{
		IPChan: ipChan,
	}

	endpointsChan := make(chan []*endpoint.Endpoint)
	logicState := logic.Logic{
		IPChan:       ipChan,
		EnpointsChan: endpointsChan,
	}

	ednsHandler := edns.Handler{
		ConnCh:  connCh,
		StateCh: endpointsChan,
	}
	httpSrv := http.Server{
		Addr:    httpAddr,
		Handler: updateHandler,
	}

	var wg sync.WaitGroup

	go ednsHandler.Run(ctx, &wg)
	wg.Add(1)

	go logicState.Run(ctx, &wg)
	wg.Add(1)

	go httprun.Run(ctx, &wg, &httpSrv)
	wg.Add(1)

	wg.Wait()

	return nil
}

func acceptLoop(listener net.Listener) <-chan net.Conn {
	ch := make(chan net.Conn)

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

func readEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s env var unset", key)
	}

	return val
}
