package socket

import "github.com/devlucas-java/luca-omegle/internal/delivery/socket/dto"

type TypeMessage string

const (
	SUBSCRIBE   TypeMessage = "SUBSCRIBE"
	UNSUBSCRIBE TypeMessage = "UNSUBSCRIBE"
	BRODCAST    TypeMessage = "BRODCAST"
	DISCONNECT  TypeMessage = "DISCONNECT"
	DASHBOARD   TypeMessage = "DASHBOARD"
)

func Subscribe(chatID string) {

}

func Unsubscribe(chatID string) {

}

func Broadcast(chatID string, msg string) {

}

func Next(chatID string) {

}

func Disconnect(chatID string, session *dto.Session) {

}

func Dashboard() {

}
