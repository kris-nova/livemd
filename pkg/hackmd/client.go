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

	// EnvironmentalVariableToken is the variable to search for a token
	EnvironmentalVariableToken string = "HACKMD_TOKEN"
)

type Client struct {

	// bearerToken is the Auth token
	bearerToken string

	// endpointf is the format string to build endpoints with
	endpointf string
}

func New(token string) *Client {
	return &Client{
		bearerToken: token,
		endpointf:   EndpointV1Format,
	}
}

func (c *Client) GET(endpoint string) (*http.Response, error) {
	return http.Get(fmt.Sprintf(c.endpointf, endpoint))
}
