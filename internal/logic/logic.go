package logic

import (
	"context"
	"log"
	"net"
	"sync"

	"sigs.k8s.io/external-dns/endpoint"
)

type Logic struct {
	DNSRecordSettings DNSRecordSettings

	IPChan       <-chan net.IP
	EnpointsChan chan<- []*endpoint.Endpoint
}

func (l *Logic) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case ip := <-l.IPChan:
			l.EnpointsChan <- l.mapIPToEndpoints(ip)
		case <-ctx.Done():
			return
		}
	}
}

func (l *Logic) mapIPToEndpoints(ip net.IP) []*endpoint.Endpoint {
	ep, err := l.DNSRecordSettings.MapIPToEndpoint(ip)
	if err != nil {
		log.Printf("unable to map the IP address to an endpoint: %v", err)
		return []*endpoint.Endpoint{}
	}
	return []*endpoint.Endpoint{ep}
}
