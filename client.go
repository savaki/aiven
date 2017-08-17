package aiven

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Client represents an authenticated gateway to aiven
type Client struct {
	token  string
	client *http.Client
}

// Do provides a generic handle for request content from aiven
func (c *Client) Do(ctx context.Context, method, url string, in, out interface{}) error {
	var body io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return errors.Wrapf(err, "unable to json marshal input")
		}
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return errors.Wrapf(err, "unable to create request for url, %v", url)
	}
	req = req.WithContext(ctx)

	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("authorization", "aivenv1 "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "api called failed, %v %v", method, url)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "unable to api content")
	}

	if err := json.Unmarshal(data, out); err != nil {
		return errors.Wrapf(err, "unable to decode response")
	}

	return nil
}

// Get specified url with authentication
func (c *Client) Get(ctx context.Context, url string, out interface{}) error {
	return c.Do(ctx, http.MethodGet, url, nil, out)
}

// Post to specified url with authentication
func (c *Client) Post(ctx context.Context, url string, in, out interface{}) error {
	return c.Do(ctx, http.MethodPost, url, in, out)
}

// Delete to specified url with authentication
func (c *Client) Delete(ctx context.Context, url string, in, out interface{}) error {
	return c.Do(ctx, http.MethodDelete, url, in, out)
}

func (c *Client) authOTP(ctx context.Context, email, password, otp string) error {
	in := map[string]string{
		"email":    email,
		"password": password,
		"otp":      otp,
	}

	out := struct {
		Errors []struct {
			Message string
		}
		Message string
		State   string
		Token   string
	}{}

	if err := c.Post(ctx, "https://api.aiven.io/v1beta/userauth", in, &out); err != nil {
		return errors.Wrapf(err, "unable to authenticate user")
	}

	if len(out.Errors) > 0 {
		return fmt.Errorf("authentication failed: %v", out.Message)
	}

	c.token = out.Token

	return nil
}

// NewOTP accepts credentials plus a one time password to return a new aiven client
func NewOTP(email, password, otp string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create cookie jar")
	}

	c := &Client{
		client: &http.Client{Jar: jar},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()
	if err := c.authOTP(ctx, email, password, otp); err != nil {
		return nil, err
	}

	return c, nil
}

// New returns a new aiven client with the specified email and password
func New(email, password string) (*Client, error) {
	return NewOTP(email, password, "")
}

// EnvAuth constructs a new client from environment variables: AIVEN_EMAIL, AIVEN_PASSWORD, and AIVEN_OTP
func EnvAuth() (*Client, error) {
	return NewOTP(os.Getenv("AIVEN_EMAIL"), os.Getenv("AIVEN_PASSWORD"), os.Getenv("AIVEN_OTP"))
}
