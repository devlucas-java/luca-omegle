package socket

import (
	"encoding/json"
	"net/http"

	"github.com/coder/websocket"
	"github.com/devlucas-java/luca-omegle/pkg/logger"
)

type Dashboard struct {
	Users int `json:"users"`
	Rooms int `json:"rooms"`
	log   *logger.Logger
}

func NewDashboard(l *logger.Logger) *Dashboard {
	return &Dashboard{
		Users: 0,
		Rooms: 0,
		log:   l,
	}
}

func (t *Dashboard) WSDashboard(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for {

		var dashboard Dashboard
		data, err := json.Marshal(dashboard)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			break
		}

		err = conn.Write(r.Context(), websocket.MessageText, data)
		if err != nil {
			break
		}
	}
}

func (t *Dashboard) GetClients(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]int{
		"users": t.Users,
		"rooms": t.Rooms,
	})
}
