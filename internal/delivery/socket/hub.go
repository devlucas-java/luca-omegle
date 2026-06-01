package socket

type TypeMessage string

const (
	SUBSCRIBE   TypeMessage = "SUBSCRIBE"
	UNSUBSCRIBE TypeMessage = "UNSUBSCRIBE"
	BRODCAST    TypeMessage = "BRODCAST"
	NEXT        TypeMessage = "NEXT"
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

func Disconnect(chatID string) {

}

func Dashboard() {

}
