package multicast

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type MessageHandler func(payload json.RawMessage, addr string) error

var handlerRegistry = make(map[string]MessageHandler)

func RegisterHandler(msgType string, handler MessageHandler) {
	handlerRegistry[msgType] = handler
}

var hostData = make(map[string]HostInfoReceiver)
var hostDataLock sync.RWMutex

func Init() {
	RegisterHandler("hostinfoSend", handleHostInfoSend)
	RegisterHandler("hostinfo", handleHostInfo)

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Failed to get hostname: %v", err)
		return
	}

	ips := getLocalIPs()
	hostDataLock.Lock()
	hostData[hostname] = HostInfoReceiver{
		Version:      "",
		BuildDate:    "",
		Revision:     "",
		Hostname:     hostname,
		IPs:          ips,
		Endpoint:     "",
		EndpointPort: 0,
	}
	hostDataLock.Unlock()
}

func RunReceivers(addr string) error {
	if len(handlerRegistry) == 0 {
		return fmt.Errorf("handler registry is empty — did you forget to call multicast.Init()?")
	}

	// mcast addr
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to resolve multicast address: %w", err)
	}

	// Interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to get interfaces: %w", err)
	}

	// for each interface
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 {
			continue
		}

		//go RunReceiver(cancel, udpAddr, &iface)
		log.Printf("Starting receiver on interface: %s [%s]", iface.Name, iface.HardwareAddr)
		go RunReceiverWithTimeoutCleanup(udpAddr, &iface, addr)
	}

	return nil
}

func RunReceiverWithTimeoutCleanup(addr *net.UDPAddr, iface *net.Interface, multicastaddr string) error {
	conn, err := net.ListenMulticastUDP("udp", iface, addr)
	if err != nil {
		return fmt.Errorf("failed to listen on multicast: %w", err)
	}
	defer conn.Close()

	if err := conn.SetReadBuffer(2048); err != nil {
		log.Printf("Warning: failed to set read buffer: %v", err)
	}

	type Fragment struct {
		MessageID string `json:"id"`
		Seq       int    `json:"seq"`
		Total     int    `json:"total"`
		Data      []byte `json:"data"`
	}

	type MessageBuffer struct {
		fragments map[int][]byte
		received  int
		total     int
		createdAt time.Time
	}

	cache := make(map[string]*MessageBuffer)
	var mu sync.Mutex

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				for id, entry := range cache {
					if time.Since(entry.createdAt) > 15*time.Second {
						delete(cache, id)
					}
				}
				mu.Unlock()
			}
		}
	}()

	buf := make([]byte, 2048)

	for {
		select {
		default:
			conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}
				return fmt.Errorf("UDP read failed: %w", err)
			}

			var frag Fragment
			if err := json.Unmarshal(buf[:n], &frag); err != nil {
				log.Printf("Invalid fragment JSON: %v", err)
				continue
			}

			mu.Lock()
			entry, exists := cache[frag.MessageID]
			if !exists {
				entry = &MessageBuffer{
					fragments: make(map[int][]byte),
					total:     frag.Total,
					createdAt: time.Now(),
				}
				cache[frag.MessageID] = entry
			}
			if _, ok := entry.fragments[frag.Seq]; !ok {
				entry.fragments[frag.Seq] = frag.Data
				entry.received++
			}

			if entry.received == entry.total {
				success := true
				totalLen := 0
				for i := 1; i <= entry.total; i++ {
					part, ok := entry.fragments[i]
					if !ok {
						log.Printf("Missing fragment %d in message %s", i, frag.MessageID)
						success = false
						break
					}
					totalLen += len(part)
				}

				if success {
					full := make([]byte, totalLen)
					offset := 0
					for i := 1; i <= entry.total; i++ {
						copy(full[offset:], entry.fragments[i])
						offset += len(entry.fragments[i])
					}

					var generic GenericMessage
					if err := json.Unmarshal(full, &generic); err != nil {
						log.Printf("Invalid generic message: %s", err)
					} else {
						handler, ok := handlerRegistry[generic.Type]
						if !ok {
							log.Printf("No handler for type: %s", generic.Type)
						} else {
							err := handler(generic.Payload, multicastaddr)
							if err != nil {
								log.Printf("Handler error: %s", err)
							}
						}
					}
				}

				delete(cache, frag.MessageID)
			}
			mu.Unlock()
		}
	}
}

func handleHostInfoSend(payload json.RawMessage, addr string) error {
	// trigger 기능만 수행
	log.Println("✅ Received OK message (triggered)")

	go SendWithEnvelope(addr, 1500, "hostinfo", payload)
	return nil
}

func handleHostInfo(payload json.RawMessage, addr string) error {
	var info HostInfoReceiver
	if err := json.Unmarshal(payload, &info); err != nil {
		return fmt.Errorf("failed to decode host info: %w", err)
	}

	log.Printf("✅ Received full message from %s: %+v", info.Hostname, info.IPs)

	hostDataLock.Lock()
	defer hostDataLock.Unlock()

	existing, found := hostData[info.Hostname]
	if !found || !equalIPs(existing.IPs, info.IPs) || existing.Endpoint != info.Endpoint || existing.EndpointPort != info.EndpointPort || existing.Version != info.Version || existing.BuildDate != info.BuildDate || existing.Revision != info.Revision {
		hostData[info.Hostname] = info
		log.Printf("📥 Updated host data for %s", info.Hostname)
	} else {
		log.Printf("🧩 Duplicate host data for %s ignored", info.Hostname)
	}

	return nil
}

func equalIPs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[string]bool)
	for _, ip := range a {
		set[ip] = true
	}
	for _, ip := range b {
		if !set[ip] {
			return false
		}
	}
	return true
}

func GetHostData() map[string]HostInfoReceiver {
	hostDataLock.RLock()
	defer hostDataLock.RUnlock()
	copied := make(map[string]HostInfoReceiver)
	for k, v := range hostData {
		copied[k] = v
	}
	return copied
}

func getLocalIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Printf("Failed to get interface addresses: %v", err)
		return ips
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.String())
		}
	}

	return ips
}
