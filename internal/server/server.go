package server

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/romanzimoglyad/wow/internal/config"
	"github.com/romanzimoglyad/wow/internal/model"
	"github.com/romanzimoglyad/wow/internal/pow"
	"github.com/rs/zerolog/log"
)

type Server struct {
	ip           string
	port         string
	listener     net.Listener
	quit         chan interface{}
	wg           *sync.WaitGroup
	zeroNumber   int
	hashStorage  Storage
	emailStorage Storage
}

type Storage interface {
	Add(string) error
	Get(string) error
	Delete(hash string) error
}

func New(conf *config.Config, hashStorage, emailStorage Storage) *Server {
	return &Server{
		ip:           conf.Ip,
		port:         conf.Port,
		zeroNumber:   conf.ZeroNumber,
		quit:         make(chan interface{}),
		wg:           &sync.WaitGroup{},
		hashStorage:  hashStorage,
		emailStorage: emailStorage,
	}
}

func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", net.JoinHostPort(s.ip, s.port))
	log.Debug().Msgf("listening: %s", net.JoinHostPort(s.ip, s.port))
	if err != nil {
		return fmt.Errorf("error on listen %w", err)
	}
	s.wg.Add(1)
	go s.serve()
	return nil
}

func (s *Server) serve() {
	defer s.wg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Err(err)
			}
		} else {
			s.wg.Add(1)
			go func(conn net.Conn) {
				s.handleConnection(conn)
				s.wg.Done()
			}(conn)
		}
	}
}
func (s *Server) handleConnection(conn net.Conn) {
	log.Debug().Msgf("new client: %v", conn.RemoteAddr())
	defer conn.Close()
	rdr := bufio.NewReader(conn)

	for {
		req, err := rdr.ReadString('\n')
		if err != nil {
			log.Err(err).Msg("connection read error")
			return
		}
		msg, err := s.handleRequest(req, conn.RemoteAddr().String())
		if err != nil {

			log.Err(err).Msg("processing request error")
			msg := &model.Message{
				Type: model.Error,
				Body: err.Error(),
			}
			err = sendMsg(*msg, conn)
			if err != nil {
				log.Err(err).Msg("sending response error")
			}
			return
		}
		if msg != nil {

			err = sendMsg(*msg, conn)
			if err != nil {
				log.Err(err).Msg("sending response error")
			}
		}
	}

}

func (s *Server) handleRequest(req string, email string) (*model.Message, error) {
	request, err := model.FromString(req)
	if err != nil {
		return nil, err
	}
	switch request.Type {
	case model.None:
		return nil, fmt.Errorf("connection quited")
	case model.GetChallenge:
		log.Debug().Msg("GetChallenge request")
		hashBlock := pow.HashBlock{
			Ver:      pow.Ver,
			Bits:     s.zeroNumber,
			Date:     time.Now().Unix(),
			Resource: email,
			Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(10000000)))),
			Counter:  0,
		}
		hashCashMarshaled, err := json.Marshal(&hashBlock)
		if err != nil {
			return nil, fmt.Errorf("err marshal hashBlock: %v", err)
		}
		err = s.emailStorage.Add(hashBlock.Resource)
		if err != nil {
			return nil, fmt.Errorf("err add email to cache: %v", err)
		}
		msg := model.Message{
			Body: string(hashCashMarshaled),
		}
		return &msg, nil
	case model.GetMessage:
		log.Debug().Msg("GetMessage request")
		var hashBlock pow.HashBlock
		err := json.Unmarshal([]byte(request.Body), &hashBlock)
		if err != nil {
			return nil, fmt.Errorf("err unmarshal hashBlock: %w", err)
		}

		hash := pow.CalculateHash(hashBlock.String())
		// First check
		if !pow.CheckValidHash(hash, hashBlock.Bits) {
			return nil, fmt.Errorf("invalid hashcash")
		}
		// Second check
		if time.Since(time.Unix(hashBlock.Date, 0)) > 48*time.Hour {
			return nil, fmt.Errorf("invalid date")
		}
		// Third check
		if err = s.emailStorage.Get(email); err != nil {
			return nil, fmt.Errorf("unknown email")
		}
		err = s.emailStorage.Delete(email)
		if err != nil {
			log.Err(err).Msg("error while delete email from storage")
		}
		// Forth check
		if err = s.hashStorage.Get(hash); err == nil {
			return nil, fmt.Errorf("hash has already been used")
		}
		err = s.hashStorage.Add(hash)

		if err != nil {
			return nil, fmt.Errorf("err add to cache: %w", err)
		}
		return &model.Message{
			Type: 0,
			Body: model.WordOfWisdomQuotes[rand.Intn(len(model.WordOfWisdomQuotes))],
		}, nil
	}
	return &model.Message{}, nil
}

func (s *Server) Shutdown() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}
func sendMsg(msg model.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.String())
	_, err := conn.Write([]byte(msgStr))
	return err
}
