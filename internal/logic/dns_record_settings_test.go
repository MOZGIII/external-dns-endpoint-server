package logic_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/MOZGIII/external-dns-endpoint-server/internal/logic"

	"sigs.k8s.io/external-dns/endpoint"
)

func makeTestLabels(t *testing.T) endpoint.Labels {
	testLabels, err := endpoint.NewLabelsFromString("heritage=external-dns,external-dns/foo=bar")
	if err != nil {
		t.Fatalf("invalid test labels: %v", err)
	}

	return testLabels
}

func TestMapIPToEndpoint(t *testing.T) {
	testLabels := makeTestLabels(t)
	testCases := []struct {
		desc     string
		settings logic.DNSRecordSettings
		ip       net.IP
		endpoint *endpoint.Endpoint
		err      bool
	}{
		{
			desc: "empty IP",
			settings: logic.DNSRecordSettings{
				DNSName:          "example.com",
				SetIdentifier:    "testid",
				RecordTTL:        300,
				Labels:           testLabels.DeepCopy(),
				ProviderSpecific: nil,
			},
			ip:       []byte{},
			endpoint: nil,
			err:      true,
		},
		{
			desc: "IPv4",
			settings: logic.DNSRecordSettings{
				DNSName:          "example.com",
				SetIdentifier:    "testid",
				RecordTTL:        300,
				Labels:           testLabels.DeepCopy(),
				ProviderSpecific: nil,
			},
			ip: net.IPv4(127, 0, 0, 1),
			endpoint: &endpoint.Endpoint{

				DNSName:          "example.com",
				Targets:          endpoint.NewTargets("127.0.0.1"),
				RecordType:       "A",
				SetIdentifier:    "testid",
				RecordTTL:        300,
				Labels:           testLabels.DeepCopy(),
				ProviderSpecific: nil,
			},
			err: false,
		},
		{
			desc: "IPv6",
			settings: logic.DNSRecordSettings{
				DNSName:          "example.com",
				SetIdentifier:    "testid",
				RecordTTL:        300,
				Labels:           testLabels.DeepCopy(),
				ProviderSpecific: nil,
			},
			ip: net.IPv6loopback,
			endpoint: &endpoint.Endpoint{

				DNSName:          "example.com",
				Targets:          endpoint.NewTargets("::1"),
				RecordType:       "AAAA",
				SetIdentifier:    "testid",
				RecordTTL:        300,
				Labels:           testLabels.DeepCopy(),
				ProviderSpecific: nil,
			},
			err: false,
		},
	}

	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			ep, err := tC.settings.MapIPToEndpoint(tC.ip)
			if tC.err {
				if err == nil {
					t.Error("expected an error but didn't get one")
				}
			} else {
				if !reflect.DeepEqual(ep, tC.endpoint) {
					t.Errorf(
						"endpoint returned do not match the expectation:\ngot %v, expected %v",
						ep, tC.endpoint,
					)
				}
			}
		})
	}
}
