package multicast

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"reflect"
	"time"

	"golang.org/x/net/ipv4"
)

type MessageEnvelope struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func SendWithEnvelope(ctx context.Context, addr string, mtu int, typ string, payload interface{}) error {
	return RunFragmentedSender(ctx, addr, mtu, MessageEnvelope{
		Type:    typ,
		Payload: payload,
	})
}

// RunFragmentedSenderRequest sends a fragmented request message over UDP using multiple interfaces. (한번만 전송)
func RunFragmentedSender(ctx context.Context, addr string, mtu int, data any) error {
	hostname, _ := os.Hostname()
	msgID := fmt.Sprintf("%s-%s-%d", reflect.TypeOf(data).Name(), hostname, time.Now().UnixNano())

	msgBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("invalid data for marshalling: %w", err)
	}

	maxPayloadSize := mtu - 100
	totalFragments := int(math.Ceil(float64(len(msgBytes)) / float64(maxPayloadSize)))

	type Fragment struct {
		MessageID string `json:"id"`
		Seq       int    `json:"seq"`
		Total     int    `json:"total"`
		Data      []byte `json:"data"`
	}

	var fragments [][]byte
	for i := 0; i < totalFragments; i++ {
		start := i * maxPayloadSize
		end := start + maxPayloadSize
		if end > len(msgBytes) {
			end = len(msgBytes)
		}

		fragment := Fragment{
			MessageID: msgID,
			Seq:       i + 1,
			Total:     totalFragments,
			Data:      msgBytes[start:end],
		}

		j, err := json.Marshal(fragment)
		if err != nil {
			return fmt.Errorf("failed to marshal fragment %d: %w", i+1, err)
		}
		fragments = append(fragments, j)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to resolve address %s: %w", addr, err)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		hasIPv4 := false
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				hasIPv4 = true
				break
			}
		}
		if !hasIPv4 {
			continue
		}

		log.Printf("Sending fragmented message via interface: %s", iface.Name)

		go func(iface net.Interface, fragments [][]byte) {
			conn, err := net.ListenPacket("udp4", "")
			if err != nil {
				log.Printf("[%s] failed to create UDP socket: %v", iface.Name, err)
				return
			}
			defer conn.Close()

			p := ipv4.NewPacketConn(conn)
			if err := p.SetMulticastInterface(&iface); err != nil {
				log.Printf("[%s] failed to set multicast interface: %v", iface.Name, err)
				return
			}

			// 메시지를 한 번만 전송
			select {
			case <-ctx.Done():
				log.Printf("[%s] sender canceled", iface.Name)
				return
			default:
				for i := 0; i < 3; i++ {
					for _, fragment := range fragments {
						_, err := p.WriteTo(fragment, nil, udpAddr)
						time.Sleep(10 * time.Millisecond)
						if err != nil {
							log.Printf("[%s] send fragment failed: %v", iface.Name, err)
						} else {
							log.Printf("[%s] sent fragment: %s", iface.Name, fragment)
						}
						time.Sleep(300 * time.Millisecond)
					}
				}
			}
		}(iface, fragments)
	}

	return nil
}

// RunFragmentedSender sends a fragmented message over UDP using multiple interfaces. (반복적으로 전송 특정 초 입력)
func RunFragmentedSenderCicle(ctx context.Context, addr string, mtu int, data any, second time.Duration) error {
	hostname, _ := os.Hostname()
	msgID := fmt.Sprintf("%s-%d", hostname, time.Now().UnixNano())

	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %w", err)
	}

	msgBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("invalid data for marshalling: %w", err)
	}

	maxPayloadSize := mtu - 100
	totalFragments := int(math.Ceil(float64(len(msgBytes)) / float64(maxPayloadSize)))

	type Fragment struct {
		MessageID string `json:"id"`
		Seq       int    `json:"seq"`
		Total     int    `json:"total"`
		Data      []byte `json:"data"`
	}

	var fragments [][]byte
	for i := 0; i < totalFragments; i++ {
		start := i * maxPayloadSize
		end := start + maxPayloadSize
		if end > len(msgBytes) {
			end = len(msgBytes)
		}

		fragment := Fragment{
			MessageID: msgID,
			Seq:       i + 1,
			Total:     totalFragments,
			Data:      msgBytes[start:end],
		}

		j, err := json.Marshal(fragment)
		if err != nil {
			return fmt.Errorf("failed to marshal fragment %d: %w", i+1, err)
		}
		fragments = append(fragments, j)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to resolve address %s: %w", addr, err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		hasIPv4 := false
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				hasIPv4 = true
				break
			}
		}
		if !hasIPv4 {
			continue
		}

		log.Printf("Sending fragmented message via interface: %s", iface.Name)

		go func(iface net.Interface, fragments [][]byte) {
			conn, err := net.ListenPacket("udp4", "")
			if err != nil {
				log.Printf("[%s] failed to create UDP socket: %v", iface.Name, err)
				return
			}
			defer conn.Close()

			p := ipv4.NewPacketConn(conn)
			if err := p.SetMulticastInterface(&iface); err != nil {
				log.Printf("[%s] failed to set multicast interface: %v", iface.Name, err)
				return
			}

			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					log.Printf("[%s] sender stopped", iface.Name)
					return
				case <-ticker.C:
					for _, fragment := range fragments {
						_, err := p.WriteTo(fragment, nil, udpAddr)
						if err != nil {
							log.Printf("[%s] send fragment failed: %v", iface.Name, err)
						} else {
							log.Printf("[%s] sent fragment: %s", iface.Name, fragment)
						}
						time.Sleep(10 * time.Millisecond)
					}
				}
			}
		}(iface, fragments)
	}

	return nil
}

func RunFragmentedSenderHostInfo(ctx context.Context, addr string, mtu int) error {
	type hostInfoSender struct {
		Hostname string   `json:"hostname"`
		IPs      []string `json:"ips"`
	}

	hostname, _ := os.Hostname()
	msgID := fmt.Sprintf("%s-%d", hostname, time.Now().UnixNano())

	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %w", err)
	}

	var localIPs []string
	for _, ifaceAddr := range ifaces {
		addrsI, _ := ifaceAddr.Addrs()
		for _, a := range addrsI {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				if ipnet.IP.String() == "127.0.0.1" {
					continue
				}
				localIPs = append(localIPs, ipnet.IP.String())
			}
		}
	}

	hostInfo := hostInfoSender{
		Hostname: hostname,
		IPs:      localIPs,
	}

	msgBytes, err := json.Marshal(hostInfo)
	if err != nil {
		return fmt.Errorf("invalid data for marshalling: %w", err)
	}

	maxPayloadSize := mtu - 100
	totalFragments := int(math.Ceil(float64(len(msgBytes)) / float64(maxPayloadSize)))

	type Fragment struct {
		MessageID string `json:"id"`
		Seq       int    `json:"seq"`
		Total     int    `json:"total"`
		Data      []byte `json:"data"`
	}

	var fragments [][]byte
	for i := 0; i < totalFragments; i++ {
		start := i * maxPayloadSize
		end := start + maxPayloadSize
		if end > len(msgBytes) {
			end = len(msgBytes)
		}

		fragment := Fragment{
			MessageID: msgID,
			Seq:       i + 1,
			Total:     totalFragments,
			Data:      msgBytes[start:end],
		}

		j, err := json.Marshal(fragment)
		if err != nil {
			return fmt.Errorf("failed to marshal fragment %d: %w", i+1, err)
		}
		fragments = append(fragments, j)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to resolve address %s: %w", addr, err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 {
			continue
		}

		addrs, _ := iface.Addrs()
		hasIPv4 := false
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				hasIPv4 = true
				break
			}
		}
		if !hasIPv4 {
			continue
		}

		log.Printf("Sending fragmented message via interface: %s", iface.Name)

		go func(iface net.Interface, fragments [][]byte) {
			conn, err := net.ListenPacket("udp4", "")
			if err != nil {
				log.Printf("[%s] failed to create UDP socket: %v", iface.Name, err)
				return
			}
			defer conn.Close()

			p := ipv4.NewPacketConn(conn)
			if err := p.SetMulticastInterface(&iface); err != nil {
				log.Printf("[%s] failed to set multicast interface: %v", iface.Name, err)
				return
			}

			select {
			case <-ctx.Done():
				log.Printf("[%s] sender canceled", iface.Name)
				return
			default:
				for i := 0; i < 3; i++ {
					for _, fragment := range fragments {
						_, err := p.WriteTo(fragment, nil, udpAddr)
						time.Sleep(10 * time.Millisecond)
						if err != nil {
							log.Printf("[%s] send fragment failed: %v", iface.Name, err)
						} else {
							log.Printf("[%s] sent fragment: %s", iface.Name, fragment)
						}
						time.Sleep(300 * time.Millisecond)
					}
				}
			}
		}(iface, fragments)
	}

	return nil
}
