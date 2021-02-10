package update

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type Handler struct {
	IPChan chan<- net.IP
}

var _ http.Handler = (*Handler)(nil)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicf("failed to read the IP from the incoming HTTP request: %v", err)
		return
	}

	ip := net.ParseIP(string(data))
	if ip == nil {
		log.Panicf("failed to parse the IP from the incoming HTTP request: %v", err)
		return
	}

	h.IPChan <- ip
}
