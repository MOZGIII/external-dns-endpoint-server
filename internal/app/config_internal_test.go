package app

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/MOZGIII/external-dns-endpoint-server/internal/logic"

	"sigs.k8s.io/external-dns/endpoint"
)

// nolint: gochecknoglobals
var mu sync.Mutex

func TestReadDNSRecordSettings_defaults(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	expectedLabels, _ := endpoint.NewLabelsFromString("heritage=external-dns")
	expected := &logic.DNSRecordSettings{
		DNSName:          "example.com",
		RecordTTL:        300,
		Labels:           expectedLabels,
		ProviderSpecific: endpoint.ProviderSpecific([]endpoint.ProviderSpecificProperty{}),
	}

	os.Setenv("DNS_NAME", "example.com")
	defer os.Unsetenv("DNS_NAME")

	got := readDNSRecordSettings()

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestReadDNSRecordSettings_all(t *testing.T) {
	mu.Lock()
	defer mu.Unlock()

	expectedLabels, _ := endpoint.NewLabelsFromString("heritage=external-dns,external-dns/test=test")
	expected := &logic.DNSRecordSettings{
		DNSName:       "example.com",
		SetIdentifier: "testid",
		RecordTTL:     120,
		Labels:        expectedLabels,
		ProviderSpecific: endpoint.ProviderSpecific([]endpoint.ProviderSpecificProperty{
			{Name: "foo", Value: "bar"},
		}),
	}

	os.Setenv("DNS_NAME", "example.com")
	defer os.Unsetenv("DNS_NAME")

	os.Setenv("RECORD_TTL", "120")
	defer os.Unsetenv("RECORD_TTL")

	os.Setenv("SET_IDENTIFIER", "testid")
	defer os.Unsetenv("SET_IDENTIFIER")

	os.Setenv("LABELS", "heritage=external-dns,external-dns/test=test")
	defer os.Unsetenv("LABELS")

	os.Setenv("PROVIDER_SPECIFIC", `[{"name": "foo", "value": "bar"}]`)
	defer os.Unsetenv("PROVIDER_SPECIFIC")

	got := readDNSRecordSettings()

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
