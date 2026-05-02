package tcp

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"

	"minesweeper/internal/domain"
	"minesweeper/internal/service"
)

type Handler struct {
	mu       sync.Mutex
	clients  map[string]*ClientConn
	game     *service.GameService
	sessions *service.SessionService
	logger   *slog.Logger
}

func NewHandler(
	game *service.GameService,
	sessions *service.SessionService,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		clients:  make(map[string]*ClientConn),
		game:     game,
		sessions: sessions,
		logger:   logger,
	}
}

func (h *Handler) Handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	_, _ = conn.Write([]byte("Enter name:\n"))

	nameRaw, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	name := strings.TrimSpace(nameRaw)
	if name == "" {
		name = conn.RemoteAddr().String()
	}

	client := &ClientConn{
		Name: name,
		Conn: conn,
	}

	h.addClient(client)
	defer h.removeClient(name)

	h.logger.Info("client connected", "name", name, "addr", conn.RemoteAddr().String())

	h.Broadcast(fmt.Sprintf("[SERVER] %s joined\n", name))
	h.sendHelp(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		cmd := ParseCommand(strings.TrimSpace(msg))
		h.handleCommand(cmd, client)
	}
}

func (h *Handler) handleCommand(cmd Command, client *ClientConn) {
	switch cmd.Name {
	case "":
		return

	case "help":
		h.sendHelp(client.Conn)

	case "create":
		h.handleCreate(cmd, client)

	case "join":
		h.handleJoin(cmd, client)

	case "sessions":
		h.handleSessions(client)

	case "state":
		h.handleState(client)

	case "open":
		h.handleOpen(cmd, client)

	case "flag":
		h.handleFlag(cmd, client)

	case "reset":
		h.handleReset(client)

	case "list":
		h.handleList(client)

	default:
		_, _ = client.Conn.Write([]byte("unknown command. type: help\n"))
	}
}

func (h *Handler) handleCreate(cmd Command, client *ClientConn) {
	if len(cmd.Args) != 1 {
		_, _ = client.Conn.Write([]byte("usage: create session_id\n"))
		return
	}

	sessionID := cmd.Args[0]

	session, err := h.sessions.CreateSession(sessionID)
	if err != nil {
		_, _ = client.Conn.Write([]byte("failed to create session: " + err.Error() + "\n"))
		return
	}

	session, err = h.sessions.JoinSession(session.ID, client.Name)
	if err != nil {
		_, _ = client.Conn.Write([]byte("failed to join session: " + err.Error() + "\n"))
		return
	}

	client.SessionID = session.ID

	_, _ = client.Conn.Write([]byte("created and joined session: " + session.ID + "\n"))
	_, _ = client.Conn.Write([]byte(session.Game.Render()))
}

func (h *Handler) handleState(client *ClientConn) {
	session, err := h.sessions.GetSession(client.SessionID)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	_, _ = client.Conn.Write([]byte(session.Game.Render()))
}

func (h *Handler) handleOpen(cmd Command, client *ClientConn) {
	x, y, ok := ParseXY(cmd.Args)
	if !ok {
		_, _ = client.Conn.Write([]byte("usage: open x y\n"))
		return
	}

	session, err := h.sessions.GetSession(client.SessionID)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	err = session.Game.Open(x, y)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	h.BroadcastToSession(session.ID, session.Game.Render())
	h.broadcastStatusToSession(session.ID, session.Game.Status)
}

func (h *Handler) handleFlag(cmd Command, client *ClientConn) {
	x, y, ok := ParseXY(cmd.Args)
	if !ok {
		_, _ = client.Conn.Write([]byte("usage: flag x y\n"))
		return
	}

	session, err := h.sessions.GetSession(client.SessionID)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	err = session.Game.ToggleFlag(x, y)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	h.BroadcastToSession(session.ID, session.Game.Render())
	h.broadcastStatusToSession(session.ID, session.Game.Status)
}

func (h *Handler) handleReset(client *ClientConn) {
	session, err := h.sessions.GetSession(client.SessionID)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	session.Game.Reset()

	h.BroadcastToSession(session.ID, "[GAME] reset\n")
	h.BroadcastToSession(session.ID, session.Game.Render())
}

func (h *Handler) handleList(client *ClientConn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	names := make([]string, 0, len(h.clients))
	for name := range h.clients {
		names = append(names, name)
	}

	_, _ = client.Conn.Write([]byte("players: " + strings.Join(names, ", ") + "\n"))
}

func (h *Handler) sendHelp(conn net.Conn) {
	_, _ = conn.Write([]byte(
		"Commands:\n" +
			"  create session_id   - create and join session\n" +
			"  join session_id     - join existing session\n" +
			"  sessions            - list sessions\n" +
			"  state               - show current session board\n" +
			"  open x y            - open cell\n" +
			"  flag x y            - toggle flag\n" +
			"  reset               - reset current session\n" +
			"  list                - show connected players\n" +
			"  help                - show help\n",
	))
}

func (h *Handler) addClient(client *ClientConn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.Name] = client
}

func (h *Handler) removeClient(name string) {
	h.mu.Lock()
	delete(h.clients, name)
	h.mu.Unlock()

	h.logger.Info("client disconnected", "name", name)
	h.Broadcast(fmt.Sprintf("[SERVER] %s left\n", name))
}

func (h *Handler) Broadcast(msg string) {
	h.mu.Lock()
	clients := make([]*ClientConn, 0, len(h.clients))
	for _, client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.Unlock()

	for _, client := range clients {
		_, _ = client.Conn.Write([]byte(msg))
	}
}

func (h *Handler) handleJoin(cmd Command, client *ClientConn) {
	if len(cmd.Args) != 1 {
		_, _ = client.Conn.Write([]byte("usage: join session_id\n"))
		return
	}

	sessionID := cmd.Args[0]

	session, err := h.sessions.JoinSession(sessionID, client.Name)
	if err != nil {
		_, _ = client.Conn.Write([]byte("failed to join session: " + err.Error() + "\n"))
		return
	}

	client.SessionID = session.ID

	_, _ = client.Conn.Write([]byte("joined session: " + sessionID + "\n"))
	_, _ = client.Conn.Write([]byte(session.Game.Render()))
}

func (h *Handler) handleSessions(client *ClientConn) {
	_, _ = client.Conn.Write([]byte(h.sessions.ListSessions()))
}

func (h *Handler) BroadcastToSession(sessionID string, msg string) {
	h.mu.Lock()
	clients := make([]*ClientConn, 0, len(h.clients))

	for _, client := range h.clients {
		if client.SessionID == sessionID {
			clients = append(clients, client)
		}
	}

	h.mu.Unlock()

	for _, client := range clients {
		_, _ = client.Conn.Write([]byte(msg))
	}
}

func (h *Handler) broadcastStatusToSession(sessionID string, status domain.GameStatus) {
	switch status {
	case domain.StatusLose:
		h.BroadcastToSession(sessionID, "[GAME OVER] lose\n")
	case domain.StatusWin:
		h.BroadcastToSession(sessionID, "[GAME OVER] win\n")
	}
}
