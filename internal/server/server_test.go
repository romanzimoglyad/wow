package server

import (
	"encoding/json"
	"testing"

	"github.com/romanzimoglyad/wow/internal/model"

	"github.com/stretchr/testify/require"

	"github.com/romanzimoglyad/wow/internal/pow"

	"github.com/stretchr/testify/suite"

	"github.com/romanzimoglyad/wow/internal/config"
	"github.com/romanzimoglyad/wow/internal/storage"
	"github.com/stretchr/testify/assert"
)

type HandleRequestSuite struct {
	suite.Suite
	server    *Server
	cfg       *config.Config
	testEmail string
}

func (h *HandleRequestSuite) SetupTest() {

	cfg, err := config.New(".env")
	if err != nil {
		h.T().Fatal(err)
	}

	h.server = New(cfg, storage.New(), storage.New())
	h.testEmail = "test1@mail.ru"
}

func (h *HandleRequestSuite) Test_ShouldErrorOn0Type() {
	_, err := h.server.handleRequest("0|", h.testEmail)
	assert.Error(h.T(), err)
	assert.Equal(h.T(), "connection quited", err.Error())
}

func (h *HandleRequestSuite) Test_ShouldErrorOnBadRequest() {
	_, err := h.server.handleRequest("IDDQD", h.testEmail)
	assert.Error(h.T(), err)
	assert.Equal(h.T(), "wrong request", err.Error())
}

func (h *HandleRequestSuite) Test_ShouldReceiveChallenge() {
	got, err := h.server.handleRequest("1|", h.testEmail)

	require.NoError(h.T(), err)
	var hashBlock pow.HashBlock
	err = json.Unmarshal([]byte(got.Body), &hashBlock)
	require.NoError(h.T(), err)
	assert.Equal(h.T(), 1, hashBlock.Ver)
	assert.Equal(h.T(), 4, hashBlock.Bits)
	assert.Equal(h.T(), h.testEmail, hashBlock.Resource)
	assert.Equal(h.T(), 0, hashBlock.Counter)
}

func (h *HandleRequestSuite) Test_ShouldErrorOnBadReceiveResource() {
	_, err := h.server.handleRequest("2|", h.testEmail)
	assert.Error(h.T(), err)
}

func (h *HandleRequestSuite) Test_ShouldErrorOnReceiveResourceWithWrongHash() {
	got, err := h.server.handleRequest("1|", h.testEmail)

	require.NoError(h.T(), err)
	var hashBlock pow.HashBlock
	err = json.Unmarshal([]byte(got.Body), &hashBlock)

	require.NoError(h.T(), err)
	msg := model.Message{
		Type: model.GetMessage,
		Body: got.Body,
	}
	_, err = h.server.handleRequest(msg.String(), h.testEmail)
	require.Error(h.T(), err)

	assert.Equal(h.T(), "invalid hashcash", err.Error())
}

func (h *HandleRequestSuite) Test_ShouldErrorOnReceiveResourceWithWrongEmail() {
	got, err := h.server.handleRequest("1|", h.testEmail)

	require.NoError(h.T(), err)
	var hashBlock pow.HashBlock
	err = json.Unmarshal([]byte(got.Body), &hashBlock)
	require.NoError(h.T(), err)
	err = hashBlock.DoWork(1000000)
	require.NoError(h.T(), err)
	body, err := json.Marshal(hashBlock)
	require.NoError(h.T(), err)
	msg := model.Message{
		Type: model.GetMessage,
		Body: string(body),
	}
	_, err = h.server.handleRequest(msg.String(), "wrong email")
	require.Error(h.T(), err)
}

func (h *HandleRequestSuite) Test_ShouldErrorOnReceiveResourceWithSameHash() {
	got, err := h.server.handleRequest("1|", h.testEmail)

	require.NoError(h.T(), err)
	var hashBlock pow.HashBlock
	err = json.Unmarshal([]byte(got.Body), &hashBlock)
	require.NoError(h.T(), err)
	err = hashBlock.DoWork(1000000)
	require.NoError(h.T(), err)
	body, err := json.Marshal(hashBlock)
	require.NoError(h.T(), err)
	msg := model.Message{
		Type: model.GetMessage,
		Body: string(body),
	}
	_, err = h.server.handleRequest(msg.String(), h.testEmail)
	require.NoError(h.T(), err)
	// fake request to add email to cache
	_, err = h.server.handleRequest("1|", h.testEmail)
	require.NoError(h.T(), err)
	_, err = h.server.handleRequest(msg.String(), h.testEmail)
	require.Error(h.T(), err)

	assert.Equal(h.T(), "hash has already been used", err.Error())
}

func (h *HandleRequestSuite) Test_ShouldReceiveAndMakeChallenge() {
	got, err := h.server.handleRequest("1|", h.testEmail)

	require.NoError(h.T(), err)
	var hashBlock pow.HashBlock
	err = json.Unmarshal([]byte(got.Body), &hashBlock)
	require.NoError(h.T(), err)
	err = hashBlock.DoWork(1000000)
	require.NoError(h.T(), err)
	body, err := json.Marshal(hashBlock)
	require.NoError(h.T(), err)
	msg := model.Message{
		Type: model.GetMessage,
		Body: string(body),
	}
	got, err = h.server.handleRequest(msg.String(), h.testEmail)
	require.NoError(h.T(), err)
	assert.NotNil(h.T(), got)
	assert.True(h.T(), isWOW(got.Body))

}

func TestHandleRequestSuite(t *testing.T) {
	suite.Run(t, new(HandleRequestSuite))
}

func isWOW(quote string) bool {
	for _, wow := range model.WordOfWisdomQuotes {
		if quote == wow {
			return true
		}
	}
	return false
}
