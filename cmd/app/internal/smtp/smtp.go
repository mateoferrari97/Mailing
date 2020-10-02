package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type Client struct {
	client *smtp.Client
	config Config
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewClient(config Config) (*Client, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return &Client{
		client: nil,
		config: config,
	}, nil
}

func validateConfig(config Config) error {
	if config.Host == "" {
		return fmt.Errorf("host is required")
	}

	if config.Port == "" {
		return fmt.Errorf("port is required")
	}

	if config.Username == "" {
		return fmt.Errorf("username is required")
	}

	if config.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

func (c *Client) Open() error {
	if c.client != nil {
		if err := c.quit(); err != nil {
			return err
		}
	}

	if err := c.open(); err != nil {
		if err := c.quit(); err != nil {
			return err
		}

		return err
	}

	return nil
}

func (c *Client) open() error {
	addr := fmt.Sprintf("%s:%s", c.config.Host, c.config.Port)

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("opening connection: %v", err)
	}

	c.client = client

	tlsConfig := &tls.Config{ServerName: c.config.Host}

	return c.client.StartTLS(tlsConfig)
}

func (c *Client) Auth() error {
	a := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	if err := c.client.Auth(a); err != nil {
		return fmt.Errorf("authenticating: %v", err)
	}

	return nil
}

func (c *Client) Send(from string, to []string, subject string, message string) error {
	var data []string
	if err := c.client.Mail(from); err != nil {
		return fmt.Errorf("configurating sender (%s): %v", from, err)
	}

	data = append(data, fmt.Sprintf("From:%s\r\n", from))

	for _, r := range to {
		data = append(data, fmt.Sprintf("To:%s\r\n", r))
		if err := c.client.Rcpt(r); err != nil {
			return fmt.Errorf("configurating reciever (%s): %v", r, err)
		}
	}

	data = append(data, fmt.Sprintf("Subject:%s\r\n", subject))

	w, err := c.client.Data()
	if err != nil {
		return fmt.Errorf("configurating message: %v", err)
	}

	defer w.Close()

	data = append(data, fmt.Sprintf("\r\n%s", message))

	if _, err := w.Write([]byte(strings.Join(data, ""))); err != nil {
		return fmt.Errorf("sending message: %v", err)
	}

	return nil
}

func (c *Client) Quit() error {
	return c.quit()
}

func (c *Client) quit() error {
	if err := c.client.Quit(); err != nil {
		return fmt.Errorf("closing connection: %v", err)
	}

	c.client = nil

	return nil
}
