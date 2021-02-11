package logic

import (
	"fmt"
	"net"

	"sigs.k8s.io/external-dns/endpoint"
)

type DNSRecordSettings struct {
	DNSName          string
	SetIdentifier    string
	RecordTTL        endpoint.TTL
	Labels           endpoint.Labels
	ProviderSpecific endpoint.ProviderSpecific
}

func (s *DNSRecordSettings) MapIPToEndpoint(ip net.IP) (*endpoint.Endpoint, error) {
	var recordType string
	switch {
	case len(ip.To4()) == net.IPv4len:
		recordType = "A"
	case len(ip) == net.IPv6len:
		recordType = "AAAA"
	default:
		return nil, fmt.Errorf("invalid IP")
	}

	ep := endpoint.Endpoint{
		DNSName:          s.DNSName,
		Targets:          endpoint.NewTargets(ip.String()),
		RecordType:       recordType,
		SetIdentifier:    s.SetIdentifier,
		RecordTTL:        s.RecordTTL,
		Labels:           s.Labels,
		ProviderSpecific: s.ProviderSpecific,
	}
	return &ep, nil
}
