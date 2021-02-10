package logic

import (
	"context"
	"log"
	"net"
	"sync"

	"sigs.k8s.io/external-dns/endpoint"
)

type DNSRecordSettings struct {
	DNSName       string
	SetIdentifier string
}

type Logic struct {
	DNSRecordSettings DNSRecordSettings

	IPChan       <-chan net.IP
	EnpointsChan chan<- []endpoint.Endpoint
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

func (l *Logic) mapIPToEndpoints(ip net.IP) []endpoint.Endpoint {
	var recordType string
	switch len(ip) {
	case net.IPv4len:
		recordType = "A"
	case net.IPv6len:
		recordType = "AAAA"
	default:
		log.Printf("invalid IP")
		return []endpoint.Endpoint{}
	}

	ep := endpoint.Endpoint{
		DNSName:       l.DNSRecordSettings.DNSName,
		Targets:       endpoint.NewTargets(ip.String()),
		RecordType:    recordType,
		SetIdentifier: l.DNSRecordSettings.SetIdentifier,
	}

	return []endpoint.Endpoint{ep}
}
