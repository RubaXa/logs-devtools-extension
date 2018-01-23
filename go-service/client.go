package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Client struct {
	id     string
	conn   net.Conn
	logs   map[string]tailTarget
	closed bool
}

type tailTarget struct {
	N    int               `json:"n"`
	File string            `json:"file"`
	Sub  *StreamSubscriber `json:"-"`
}

type wsEvent struct {
	Type   string      `json:"type"`
	Detail interface{} `json:"detail"`
}

func (c *Client) Start() {
	fmt.Println("[client] Started:", c.id)

	// Ping/Pong
	go func() {
		defer c.conn.Close()

		for {
			if c.closed {
				break
			}

			err := wsutil.WriteServerMessage(c.conn, ws.OpPing, nil)
			if err != nil {
				fmt.Println("[client] [ws] Ping error:", c.id, err)
				break
			}

			time.Sleep(time.Second)
		}
	}()
}

func (c *Client) Setup(rawLogs []string) (map[string]tailTarget, error) {
	for _, v := range rawLogs {
		parts := strings.Split(v, " ")
		n, _ := strconv.Atoi(parts[0])
		target := tailTarget{
			N:    n,
			File: strings.Join(parts[1:], " "),
		}

		if _, ok := c.logs[target.File]; !ok {
			c.logs[target.File] = target
			tail, lines, err := StartTail(target.File, target.N)

			if err != nil {
				return nil, err
			}

			for _, v := range lines {
				c.PushEvent("record", v)
			}

			target.Sub = tail.Subscribe()
		}
	}

	return c.logs, nil
}

func (c *Client) PushEvent(t string, d interface{}) {
	if !c.closed {
		fmt.Println("[client] [ws] Push unacceptably, because closed", c.id)
		return
	}

	msg, err := json.Marshal(wsEvent{t, d})
	if err != nil {
		fmt.Println("[client] [ws] Push failed", c.id, err)
		return
	}

	wsutil.WriteServerMessage(c.conn, ws.OpText, msg)
}

func (c *Client) Close() {
	c.closed = true
	fmt.Println("[client] Closed:", c.id)
}
