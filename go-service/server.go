package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gobwas/ws"
	"github.com/rs/cors"
)

type HttpServer struct {
	Host    string
	clients map[string]*Client
}

func (srv *HttpServer) AddClient(id string, conn net.Conn) {
	srv.RemoveClient(id)

	client := &Client{
		id:   id,
		conn: conn,
	}

	client.Start()
	srv.clients[id] = client
}

func (srv *HttpServer) RemoveClient(id string) {
	client, ok := srv.clients[id]

	if ok {
		delete(srv.clients, id)
		client.Close()
	}
}

func (srv *HttpServer) PushEvent(t string, d interface{}) {
	for id, client := range srv.clients {
		if client.closed {
			srv.RemoveClient(id)
		} else {
			srv.PushEvent(t, d)
		}
	}
}

func StartServer(host string) (*HttpServer, error) {
	srv := &HttpServer{
		Host:    host,
		clients: make(map[string]*Client, 0),
	}
	c := cors.AllowAll()
	api := map[string]func(w http.ResponseWriter, r *http.Request){
		"version": func(w http.ResponseWriter, r *http.Request) {
			end(w, nil, "0.1.0")
		},

		"setup": func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			id := q.Get("id")
			rawLogs := q["log"]

			client, ok := srv.clients[id]
			if !ok {
				end(w, fmt.Errorf("[client] Not found: %s", id), nil)
				return
			}

			logs, err := client.Setup(rawLogs)
			end(w, err, logs)
		},

		"ws": func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			parts := strings.Split(r.URL.Path, "/ws/")
			id := parts[1]

			srv.AddClient(id, conn)
		},

		"404": func(w http.ResponseWriter, r *http.Request) {
			end(w, fmt.Errorf("[api] 404"), nil)
		},
	}

	re := regexp.MustCompile("/([a-z]+)/")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := re.FindStringSubmatch(r.URL.Path)
		method := "404"

		if len(parts) > 0 {
			method = parts[1]
		}

		handle, ok := api[method]

		if ok {
			handle(w, r)
		} else {
			end(w, fmt.Errorf("[server] [api] not found: %s", r.URL.Path), nil)
		}
	})

	fmt.Printf("Start server: http://%s/\n", srv.Host)
	return srv, http.ListenAndServe(srv.Host, c.Handler(handler))
}

func end(w http.ResponseWriter, err error, v interface{}) {
	var jsonBody []byte

	if err == nil {
		jsonBody, err = json.Marshal(v)
	}

	if err != nil {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)

		msg, _ := json.Marshal(err.Error())
		jsonBody = []byte("{\"error\":" + string(msg) + "}")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
