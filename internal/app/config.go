package app

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/MOZGIII/external-dns-endpoint-server/internal/logic"

	"sigs.k8s.io/external-dns/endpoint"
)

func readEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%s env var unset", key)
	}

	return val
}

func readEnvWithDefault(key, defval string) string {
	val := os.Getenv(key)
	if val == "" {
		return defval
	}

	return val
}

func readDNSRecordSettings() *logic.DNSRecordSettings {
	dnsName := readEnv("DNS_NAME")
	setIdentifier := readEnvWithDefault("SET_IDENTIFIER", "")
	recordTTLStr := readEnvWithDefault("RECORD_TTL", "300")
	labelsStr := readEnvWithDefault("LABELS", "heritage=external-dns")
	providerSpecificStr := readEnvWithDefault("PROVIDER_SPECIFIC", "[]")

	recordTTLInt, err := strconv.Atoi(recordTTLStr)
	if err != nil {
		log.Fatalf("unable to parse record TTL: %v", err)
	}

	recordTTL := endpoint.TTL(recordTTLInt)

	labels, err := endpoint.NewLabelsFromString(labelsStr)
	if err != nil {
		log.Fatalf("unable to parse labels: %v", err)
	}

	var providerSpecific endpoint.ProviderSpecific
	if err := json.Unmarshal([]byte(providerSpecificStr), &providerSpecific); err != nil {
		log.Fatalf("unable to parse provider specific: %v", err)
	}

	return &logic.DNSRecordSettings{
		DNSName:          dnsName,
		SetIdentifier:    setIdentifier,
		RecordTTL:        recordTTL,
		Labels:           labels,
		ProviderSpecific: providerSpecific,
	}
}
