package httprun

import (
	"context"
	"log"
	"net/http"
	"sync"
)

func Run(ctx context.Context, wg *sync.WaitGroup, srv *http.Server) {
	defer wg.Done()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
}
