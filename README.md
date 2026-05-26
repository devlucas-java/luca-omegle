# Go WebRTC Playground

Proyecto experimental para aprender WebRTC, networking realtime y arquitectura de servidores de baja latencia usando Go.

La idea principal es construir una plataforma estilo Omegle:
- comunicación peer-to-peer
- video y audio en tiempo real
- chat aleatorio
- sin autenticación
- matchmaking simple
- señalización WebRTC
- backend performático en Go

## Objetivos del proyecto

Este proyecto existe para aprender:

- WebRTC
- RTP / RTCP
- ICE / STUN / TURN
- signaling servers
- WebSockets
- networking concurrente en Go
- sistemas realtime
- manejo de conexiones masivas
- arquitectura distribuida para multiplayer/chat

## Stack

- Go
- Pion WebRTC
- gnet
- WebSockets
- Docker
- Redis (futuro)
- PostgreSQL (futuro)

## Librerías principales

- :contentReference[oaicite:0]{index=0}
- :contentReference[oaicite:1]{index=1}

## Funcionalidades planeadas

- [ ] signaling server
- [ ] emparejamiento aleatorio
- [ ] video peer-to-peer
- [ ] chat de texto
- [ ] voice chat
- [ ] sistema básico anti-spam
- [ ] reconexión
- [ ] métricas y observabilidad
- [ ] soporte TURN/STUN
- [ ] escalado horizontal

## Estructura futura

```text
/internal
    /signaling
    /matchmaking
    /rtc
    /gateway
    /metrics