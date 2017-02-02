package xiaomi

import (
	"encoding/json"
	"log"
	"net"
)

type (
	Response struct {
		Cmd   string `json:"cmd"`
		Model string `json:"model"`
		Sid   string `json:"sid"`
		//ShortId int `json:"short_id,integer"`
		Token string      `json:"token,omitempty"`
		IP    string      `json:"ip,omitempty"`
		Port  string      `json:"port,omitempty"`
		Data  interface{} `json:"data"`
	}

	Request struct {
		Cmd string `json:"cmd"`
		Sid string `json:"sid,omitempty"`
	}

	Gateway struct {
		Addr  *net.UDPAddr
		Sid   string
		Token string
	}

	object map[string]interface{}
)

var (
	conn     *net.UDPConn
	gateways map[string]Gateway
)

const (
	multicastIp   = "224.0.0.50"
	multicastPort = "9898"
	//multicastAddr         = "239.255.255.250:9898"
	maxDatagramSize = 8192
)

func init() {
	go serveMulticastUDP(multicastIp+":"+multicastPort, connHandler, msgHandler)
}

func GetStatus() string { //test
	return ""
}

func (this *Gateway) sendMessage(msg string) {
	sendMessage(this.Addr, msg)
}

func sendMessage(addr *net.UDPAddr, msg string) {
	req, err := json.Marshal(Request{Cmd: msg})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(req))
	log.Println(addr)
	conn.WriteMsgUDP([]byte(req), nil, addr)
}

func connHandler() {
	pingAddr, err := net.ResolveUDPAddr("udp", multicastIp+":4321")
	if err != nil {
		log.Fatal(err)
	}
	sendMessage(pingAddr, "whois")
}

func msgHandler(resp *Response) {
	switch resp.Cmd {
	case "iam":
		log.Printf("%+v", resp)
		_, err := net.ResolveUDPAddr("udp", resp.IP+":"+multicastPort)
		if err != nil {
			log.Fatal(err)
		}

	case "heartbeat":
		log.Printf("%+v", resp)
	}
}

func serveMulticastUDP(a string, connectedHandler func(), msgHandler func(resp *Response)) {
	addr, err := net.ResolveUDPAddr("udp", a)
	if err != nil {
		log.Fatal(err)
	}
	conn, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Panic(err)
	}

	conn.SetReadBuffer(maxDatagramSize)
	for {
		b := make([]byte, maxDatagramSize)
		n, _, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		//log.Println(n, "bytes read from", src)
		//log.Println(string(b[:n]))

		resp := Response{}
		err = json.Unmarshal(b[:n], &resp)
		if err != nil {
			log.Fatal(err)
		}
		msgHandler(&resp)
	}
}
