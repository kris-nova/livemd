/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package hackmd

import (
	"bytes"
	"fmt"
	"net/http"
)

const (

	// BearerTokenHeaderKey is the HTTP header key to pass with each request
	//
	// Taken from https://hackmd.io/@hackmd-api/developer-portal
	// Header: Authorization Bearer <token>
	// curl "https://api.hackmd.io/v1/me" -H "Authorization: Bearer <token>"
	BearerTokenHeaderKey string = "Authorization"

	// BearerTokenHeaderValueFormat is the HTTP header value to format with
	// each request.
	//
	// Taken from https://hackmd.io/@hackmd-api/developer-portal
	// Header: Authorization Bearer <token>
	// curl "https://api.hackmd.io/v1/me" -H "Authorization: Bearer <token>"
	BearerTokenHeaderValueFormat string = "Bearer %s"

	// EndpointV1Format is the endpoint formatter to use
	EndpointV1Format string = "https://api.hackmd.io/v1/%s"

	// EnvironmentalVariableHackMDToken is the variable to search for a token
	EnvironmentalVariableHackMDToken string = "HACKMD_TOKEN"

	// EnvironmentalVariableHackMDID is the variable to use for persistent hackmd notes
	EnvironmentalVariableHackMDID string = "HACKMD_ID"
)

type Client struct {

	// bearerToken is the Auth token
	bearerToken string

	// endpointf is the format string to build endpoints with
	endpointf string

	client http.Client
}

func New(token string) *Client {
	return &Client{
		bearerToken: token,
		endpointf:   EndpointV1Format,
		client:      http.Client{},
	}
}

func (c *Client) GET(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(c.endpointf, endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(BearerTokenHeaderKey, fmt.Sprintf(BearerTokenHeaderValueFormat, c.bearerToken))
	return c.client.Do(req)
}

func (c *Client) POST(endpoint string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf(c.endpointf, endpoint), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(BearerTokenHeaderKey, fmt.Sprintf(BearerTokenHeaderValueFormat, c.bearerToken))
	return c.client.Do(req)
}

func (c *Client) PATCH(endpoint string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", fmt.Sprintf(c.endpointf, endpoint), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(BearerTokenHeaderKey, fmt.Sprintf(BearerTokenHeaderValueFormat, c.bearerToken))
	return c.client.Do(req)
}

func (c *Client) DELETE(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf(c.endpointf, endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(BearerTokenHeaderKey, fmt.Sprintf(BearerTokenHeaderValueFormat, c.bearerToken))
	return c.client.Do(req)
}
