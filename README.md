# luca-omegle

Chat aleatorio en tiempo real estilo Omegle, construido en Go con WebSockets y un protocolo de mensajería inspirado en STOMP.

## Stack

- Go + `net/http`
- [`coder/websocket`](https://github.com/coder/websocket)
- Redis (sesiones y rooms via cache)
- Docker / Docker Compose

## Protocolo de mensajes

El protocolo es inspirado en STOMP. Cada frame tiene la forma:

```
COMMAND
header1:value1
header2:value2

body\x00
```

### Comandos cliente → servidor

| Comando      | Headers requeridos | Descripción                                      |
|--------------|--------------------|--------------------------------------------------|
| `SUBSCRIBE`  | —                  | Entra a la cola de espera para ser emparejado    |
| `NEXT`       | —                  | Sale del room actual y vuelve a la cola          |
| `SEND`       | `room-id`          | Envía un mensaje al room                         |
| `DISCONNECT` | —                  | Cierra la sesión limpiamente                     |
| `DASHBOARD`  | —                  | Solicita datos del dashboard (pendiente)         |

### Comandos servidor → cliente

| Comando          | Headers               | Descripción                          |
|------------------|-----------------------|--------------------------------------|
| `CONNECTED`      | `room-id`             | Confirmación de emparejamiento       |
| `MESSAGE`        | `room-id`, `username` | Mensaje de otro usuario del room     |
| `ERROR`          | `message`             | Error en el frame o en la operación  |
| `DASHBOARD_DATA` | —                     | Datos del dashboard (pendiente)      |

### Ejemplo de sesión

```
# Cliente se conecta vía WebSocket a /ws?nickname=Alice

# Cliente entra a la cola
SUBSCRIBE

\x00

# Servidor empareja con otro usuario
CONNECTED
room-id:f3a1b2c4-...

\x00

# Cliente envía mensaje
SEND
room-id:f3a1b2c4-...

Hola!\x00

# El otro usuario recibe
MESSAGE
room-id:f3a1b2c4-...
username:Alice

Hola!\x00
```

## Arquitectura

```
cmd/server/main.go              → arranque, DI, graceful shutdown
internal/
  delivery/socket/
    ws_handler.go               → acepta conexiones, parsea frames, despacha al hub
    hub.go                      → goroutines de matchmaking y eventos
    dto/
      frame.go                  → Frame struct, ParseFrame, Encode
      session.go                → Session (conn + user), Send, SendError
  application/service/
    user_service.go
    room_service.go
  domain/entity/
    user.go / room.go / message.go
  infrastructure/
    cache/                      → Redis (user_ch, room_ch)
    repository/                 → interfaces
configs/config.go               → viper + redis + postgres init
pkg/logger/logger.go            → logger con niveles y componente
```

## Levantar en local

```bash
docker compose up -d        # Redis
go run ./cmd/server
```

La UI de prueba está en `http://localhost:<SERVER_PORT>/`.

## Variables de entorno

Ver `default.env`.

| Variable      | Ejemplo     |
|---------------|-------------|
| `SERVER_PORT` | `8080`      |
| `CH_HOST`     | `localhost` |
| `CH_PORT`     | `6379`      |
| `CH_PASSWORD` | `""`        |
| `CH_DB`       | `0`         |

## Funcionalidades pendientes

- [ ] Dashboard con métricas en tiempo real
- [ ] Persistencia de mensajes en PostgreSQL
- [ ] Anti-spam básico
- [ ] Reconexión automática
- [ ] Soporte TURN/STUN para video P2P
