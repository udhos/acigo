package aci

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebsocketOpen opens websocket for receiving subscription information.
func (c *Client) WebsocketOpen() error {
	api := "/socket" + c.loginToken
	url := c.getURLws(api)
	header := http.Header{}

	d := websocket.Dialer{
		TLSClientConfig: tlsConfig(),
	}

	c.debugf("WebsocketOpen: url=%s", url)

	conn, _, errDial := d.Dial(url, header)
	if errDial != nil {
		return errDial
	}

	c.socket = conn

	return nil
}

// WebsocketReadJson reads subscription message from websocket.
func (c *Client) WebsocketReadJson(v interface{}) error {
	if c.socket == nil {
		return fmt.Errorf("websocket not open")
	}
	return c.socket.ReadJSON(v)
}
