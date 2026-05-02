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
	mu      sync.Mutex
	clients map[string]*ClientConn
	game    *service.GameService
	logger  *slog.Logger
}

func NewHandler(game *service.GameService, logger *slog.Logger) *Handler {
	return &Handler{
		clients: make(map[string]*ClientConn),
		game:    game,
		logger:  logger,
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
		h.handleCreate(client)

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

func (h *Handler) handleCreate(client *ClientConn) {
	board, err := h.game.Create()
	if err != nil {
		_, _ = client.Conn.Write([]byte("failed to create game: " + err.Error() + "\n"))
		return
	}

	h.Broadcast("[GAME] new game created\n")
	h.Broadcast(board)
}

func (h *Handler) handleState(client *ClientConn) {
	board, err := h.game.State()
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + ". type: create\n"))
		return
	}

	_, _ = client.Conn.Write([]byte(board))
}

func (h *Handler) handleOpen(cmd Command, client *ClientConn) {
	x, y, ok := ParseXY(cmd.Args)
	if !ok {
		_, _ = client.Conn.Write([]byte("usage: open x y\n"))
		return
	}

	board, status, err := h.game.Open(x, y)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	h.Broadcast(board)
	h.broadcastStatus(status)
}

func (h *Handler) handleFlag(cmd Command, client *ClientConn) {
	x, y, ok := ParseXY(cmd.Args)
	if !ok {
		_, _ = client.Conn.Write([]byte("usage: flag x y\n"))
		return
	}

	board, status, err := h.game.ToggleFlag(x, y)
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	h.Broadcast(board)
	h.broadcastStatus(status)
}

func (h *Handler) handleReset(client *ClientConn) {
	board, err := h.game.Reset()
	if err != nil {
		_, _ = client.Conn.Write([]byte("error: " + err.Error() + "\n"))
		return
	}

	h.Broadcast("[GAME] reset\n")
	h.Broadcast(board)
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

func (h *Handler) broadcastStatus(status domain.GameStatus) {
	switch status {
	case domain.StatusLose:
		h.Broadcast("[GAME OVER] lose\n")
	case domain.StatusWin:
		h.Broadcast("[GAME OVER] win\n")
	}
}

func (h *Handler) sendHelp(conn net.Conn) {
	_, _ = conn.Write([]byte(
		"Commands:\n" +
			"  create        - create new game\n" +
			"  state         - show board\n" +
			"  open x y      - open cell\n" +
			"  flag x y      - toggle flag\n" +
			"  reset         - reset current game\n" +
			"  list          - show connected players\n" +
			"  help          - show help\n",
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
