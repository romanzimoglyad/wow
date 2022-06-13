package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/romanzimoglyad/wow/internal/config"
	"github.com/romanzimoglyad/wow/internal/model"
	"github.com/romanzimoglyad/wow/internal/pow"
	"github.com/rs/zerolog/log"
)

const (
	WriteDeadlineMs = 200
)

type Client struct {
	ip                     string
	port                   string
	requestNumber          int
	maxHashCountIterations int
	sendIntervalMs         int
}

func NewClient(conf *config.Config) *Client {
	return &Client{
		ip:                     conf.Ip,
		port:                   conf.Port,
		requestNumber:          conf.ClientRequestNumber,
		maxHashCountIterations: conf.MaxHashCountIterations,
		sendIntervalMs:         conf.MaxHashCountIterations,
	}
}
func (c *Client) Start() error {

	conn, err := net.Dial("tcp", net.JoinHostPort(c.ip, c.port))

	if err != nil {
		return fmt.Errorf("error on listen %w", err)
	}
	defer conn.Close()

	for i := 0; i < c.requestNumber; i++ {
		answer, err := c.doConnection(conn, conn)
		if err != nil {
			log.Err(err).Msg("error on client")
			continue
		}
		log.Info().Msgf("Word of Wisdom: %s", answer)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (c *Client) doConnection(readerConn io.Reader, writerConn io.Writer) (string, error) {
	reader := bufio.NewReader(readerConn)
	request := model.Message{
		Type: model.GetChallenge,
	}

	// 1. request service

	err := sendMsg(request, writerConn)
	if err != nil {
		return "", fmt.Errorf("error send request: %w", err)
	}

	// 2. get response and challenge

	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error read response: %w", err)
	}
	msg, err := model.FromString(answer)
	if err != nil {
		return "", fmt.Errorf("error in response format: %w", err)
	}
	var hashBlock pow.HashBlock
	err = json.Unmarshal([]byte(msg.Body), &hashBlock)
	if err != nil {
		return "", fmt.Errorf("error unmarshal response: %w", err)
	}

	err = hashBlock.DoWork(c.maxHashCountIterations)
	if err != nil {
		return "", fmt.Errorf("error counting hash: %w; counter: %d", err, hashBlock.Counter)
	}
	log.Debug().Msgf("Hash Found!Counter: %d", hashBlock.Counter)
	data, err := json.Marshal(&hashBlock)
	if err != nil {
		return "", fmt.Errorf("error marshal hash: %w", err)
	}
	request = model.Message{
		Type: model.GetMessage,
		Body: string(data),
	}

	// 3. send response

	err = sendMsg(request, writerConn)
	if err != nil {
		return "", fmt.Errorf("error send request: %w", err)
	}

	// 4. get word of wisdom
	answer, err = reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error read response: %w", err)
	}

	msg, err = model.FromString(answer)
	if err != nil {
		return "", fmt.Errorf("error in response format: %w", err)
	}

	return msg.Body, nil
}
func sendMsg(msg model.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.String())
	_, err := conn.Write([]byte(msgStr))
	return err
}
