package pkgEmail

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/gomail.v2"
)

type mockDialer struct {
	sendFunc func(...*gomail.Message) error
}

func (m *mockDialer) DialAndSend(msgs ...*gomail.Message) error {
	return m.sendFunc(msgs...)
}

func TestNew_ValidConfig(t *testing.T) {
	cfg := &Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		Username: "user",
		Password: "pass",
		From:     "from@example.com",
	}

	client, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNew_InvalidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
	}{
		{"missing host", &Config{Username: "u", Password: "p", From: "f"}},
		{"missing username", &Config{SMTPHost: "h", Password: "p", From: "f"}},
		{"missing password", &Config{SMTPHost: "h", Username: "u", From: "f"}},
		{"missing from", &Config{SMTPHost: "h", Username: "u", Password: "p"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.cfg)
			assert.Error(t, err)
			assert.Nil(t, client)
		})
	}
}

func TestClient_SendEmail_Success(t *testing.T) {
	cfg := &Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		Username: "user",
		Password: "pass",
		From:     "from@example.com",
	}

	client, _ := New(cfg)

	origNewDialer := newDialer
	newDialer = func(_ string, _ int, _ string, _ string) Dialer {
		return &mockDialer{
			sendFunc: func(_ ...*gomail.Message) error { return nil },
		}
	}
	defer func() { newDialer = origNewDialer }()

	err := client.SendEmail([]string{"to@example.com"}, "subject", "body", false)
	assert.NoError(t, err)
}

func TestClient_SendEmail_Failure(t *testing.T) {
	cfg := &Config{
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		Username: "user",
		Password: "pass",
		From:     "from@example.com",
	}

	client, _ := New(cfg)

	origNewDialer := newDialer
	newDialer = func(_ string, _ int, _ string, _ string) Dialer {
		return &mockDialer{
			sendFunc: func(_ ...*gomail.Message) error {
				return errors.New("smtp error")
			},
		}
	}
	defer func() { newDialer = origNewDialer }()

	err := client.SendEmail([]string{"to@example.com"}, "subject", "body", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "smtp error")
}
