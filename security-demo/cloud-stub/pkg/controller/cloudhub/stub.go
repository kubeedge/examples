package cloudhub

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"

	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/examples/security-demo/cloud-stub/cmd/config"
)

type StubCloudHub struct {
	//context *context.Context
	wsConn *websocket.Conn
	config *config.CloudStubConfig
}

func (tm *StubCloudHub) eventReadLoop(conn *websocket.Conn, stop chan bool) {
	for {
		var event interface{}
		err := conn.ReadJSON(&event)
		if err != nil {
			fmt.Println("read error, connection will be closed: %v", err)
			stop <- true
			return
		}
		fmt.Println("cloud hub receive message %+v", event)
	}
}

func (tm *StubCloudHub) serveEvent(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("fail to upgrade towebsocket connection: %v", err)
		http.Error(w, "fail to upgrade to websocket protocol", http.StatusInternalServerError)
		return
	}
	tm.wsConn = conn
	stop := make(chan bool, 1)
	fmt.Println("edge connected")
	go tm.eventReadLoop(conn, stop)
	<-stop
	tm.wsConn = nil
	fmt.Println("edge disconnected")
}

func (tm *StubCloudHub) deviceHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println("failed to read request body with error %s", err.Error())
			w.Write([]byte("failed to read request body"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("request body is %s\n", string(body))
		var device interface{}

		err = json.Unmarshal(body, &device)
		if err != nil {
			fmt.Println("unmarshal request body error %v", err)
			w.Write([]byte("unmarshal request body error"))
		}
		var operation string
		switch req.Method {
		case "POST":
			operation = model.InsertOperation
		case "DELETE":
			operation = model.DeleteOperation
		case "PUT":
			operation = model.UpdateOperation
		}

		msgReq := model.NewMessage("").BuildRouter("edgemgr", "twin", "membership", operation).FillBody(device)

		if tm.wsConn == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("invalid websocket connection")
			return
		}

		err = tm.wsConn.WriteJSON(*msgReq)
		fmt.Println("send message to edgehub is %+v\n", *msgReq)

		if err != nil {
			fmt.Println("failed to send message to edge hub with error - %s", err.Error())
		}
		io.WriteString(w, "OK\n")

		respMsg := model.Message{}
		err = tm.wsConn.ReadJSON(&respMsg)
		if err != nil {
			fmt.Println("failed to read response message from edge hub with error - %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)

	}
}

func (tm *StubCloudHub) podHandler(w http.ResponseWriter, req *http.Request) {
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println("read body error %v", err)
			w.Write([]byte("read request body error"))
			return
		}
		fmt.Println("request body is %s\n", string(body))

		var pod v1.Pod
		if err = json.Unmarshal(body, &pod); err != nil {
			fmt.Println("unmarshal request body error %v", err)
			w.Write([]byte("unmarshal request body error"))
			return
		}
		var msgReq *model.Message
		switch req.Method {
		case "POST":
			msgReq = model.NewMessage("").BuildRouter("controller", "resource",
				"node/fake_node_id/pod/"+string(pod.UID), model.InsertOperation).FillBody(pod)
		case "DELETE":
			msgReq = model.NewMessage("").BuildRouter("controller", "resource",
				"node/fake_node_id/pod/"+string(pod.UID), model.DeleteOperation).FillBody(pod)

		case "GET":
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		if tm.wsConn == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("invalid websocket connection")
			return
		}

		err = tm.wsConn.WriteJSON(*msgReq)
		fmt.Println("send message to edgehub is %+v\n", *msgReq)

		if err != nil {
			fmt.Println("failed to send message to edge hub with error - %s", err.Error())
		}
		io.WriteString(w, "OK\n")

		respMsg := model.Message{}
		err = tm.wsConn.ReadJSON(&respMsg)
		if err != nil {
			fmt.Println("failed to read response message from edge hub with error - %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
	return
}

func (tm *StubCloudHub) placementHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Println(" Config placement url : %s", tm.config.PlacementURL)
		w.Write([]byte(tm.config.PlacementURL))
		w.WriteHeader(http.StatusOK)
	}
	return
}

func (tm *StubCloudHub) wsHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("at wshandler")
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, req, req.Header)
	if err != nil {
		fmt.Println("websocket upgrade failed with error - %s", err.Error())
		return
	}

	if tm.wsConn == nil {
		tm.wsConn = conn
	}

	for {
		msg := model.Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("read message failed with error - %s", err.Error())
			break
		}
		fmt.Println(" received message on websocket == %v", msg)
	}
}

func (tm *StubCloudHub) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/placement_external/message_queue", tm.placementHandler) // for cloudhub url placement handler
	router.HandleFunc("/{group_id}/events", tm.serveEvent)                         // for edge-hub
	router.HandleFunc("/{project_id}/{node_id}/events", tm.wsHandler)              // for cloudhub url placement handler
	s := http.Server{
		Addr:    "127.0.0.1:20000",
		Handler: router,
	}

	fmt.Println("Start cloud hub service")
	err := s.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe: %v", err)
	}
}

func (tm *StubCloudHub) PlacementServer() {
	fmt.Println("started placement server")
	router := mux.NewRouter()
	router.HandleFunc("/pod", tm.podHandler) // for pod test
	router.HandleFunc("/device", tm.deviceHandler)

	s := http.Server{
		Addr:    "127.0.0.1:30000",
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe: %v", err)
	}
}

func NewCloudStub(config *config.CloudStubConfig) *StubCloudHub {
	return &StubCloudHub{
		config: config,
	}
}
