package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type UploadRequest struct {
	AgentID string
}

type UploadResponse struct {
	Conn string
}

func main() {
	httpClient := http.Client{}

	reqData, err := json.Marshal(UploadRequest{
		AgentID: "some_agent_id",
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:7777/upload", bytes.NewBuffer(reqData))
	if err != nil {
		panic(err)
	}

	rawResp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	var resp UploadResponse

	err = fromBody(rawResp.Body, &resp)
	if err != nil {
		panic(err)
	}

	udpDestAddr, err := net.ResolveUDPAddr("udp4", resp.Conn)
	if err != nil {
		panic(err)
	}

	socket, err := net.DialUDP("udp", nil, udpDestAddr)
	if err != nil {
		panic(err)
	}

	sent, err := socket.Write([]byte("message from client not server"))
	if err != nil {
		panic(err)
	}
	fmt.Println("LAL N bytes sent:", sent, udpDestAddr.String())
}

func fromBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return errors.WithStack(decoder.Decode(&dest))
}
