package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shubhangcs/agromart-server/internal/hub"
	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	// Allow all origins — tighten this in production.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// wsIncoming is the JSON envelope a client sends over WebSocket.
type wsIncoming struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

// ChatHandler handles direct-message endpoints.
type ChatHandler struct {
	chatStore store.ChatStore
	hub       *hub.Hub
	logger    *slog.Logger
}

func NewChatHandler(chatStore store.ChatStore, h *hub.Hub, logger *slog.Logger) *ChatHandler {
	return &ChatHandler{
		chatStore: chatStore,
		hub:       h,
		logger:    logger,
	}
}

// HandleWebSocket godoc
// @Summary      Real-time chat WebSocket
// @Description  Upgrades the connection to WebSocket. Pass ?user_id=<uuid>. Send JSON {"receiver_id":"...","content":"..."}; receive JSON Message objects in real-time.
// @Tags         chat
// @Param        user_id query string true "Authenticated user ID"
// @Security     BearerAuth
// @Router       /chat/ws [get]
func (ch *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "user_id is required"})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ch.logger.Error("ws upgrade", "error", err)
		return
	}

	client := &hub.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	ch.hub.Register(client)

	// Write pump — forwards messages from the Send channel to the WebSocket.
	go func() {
		defer conn.Close()
		for payload := range client.Send {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				ch.logger.Error("ws write", "user_id", userID, "error", err)
				return
			}
		}
	}()

	// Read pump — blocks until the connection closes.
	defer func() {
		ch.hub.Unregister(client)
		conn.Close()
	}()

	conn.SetReadLimit(4096)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Ping ticker to keep the connection alive.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ch.logger.Error("ws read", "user_id", userID, "error", err)
			}
			return
		}
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		var incoming wsIncoming
		if err := json.Unmarshal(raw, &incoming); err != nil {
			ch.logger.Warn("ws bad payload", "user_id", userID, "error", err)
			continue
		}
		if incoming.ReceiverID == "" || incoming.Content == "" {
			continue
		}
		if incoming.ReceiverID == userID {
			continue
		}

		msg := &models.Message{
			SenderID:   userID,
			ReceiverID: incoming.ReceiverID,
			Content:    incoming.Content,
		}
		if err := ch.chatStore.SaveMessage(msg); err != nil {
			ch.logger.Error("ws save message", "user_id", userID, "error", err)
			continue
		}

		payload, _ := json.Marshal(msg)

		// Deliver to receiver if online.
		ch.hub.Deliver(incoming.ReceiverID, payload)

		// Echo back to sender so they see the saved message with its ID/timestamp.
		select {
		case client.Send <- payload:
		default:
		}
	}
}

// HandleSendMessage godoc
// @Summary      Send a message
// @Description  Sends a text message from one user to another and persists it in the database
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        body body models.SendMessageRequest true "Message payload"
// @Success      201 {object} map[string]interface{} "Saved message with id and created_at"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /chat/send [post]
func (ch *ChatHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ch.logger.Error("send message", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if req.SenderID == req.ReceiverID {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "cannot send a message to yourself"})
		return
	}

	msg := &models.Message{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}
	if err := ch.chatStore.SaveMessage(msg); err != nil {
		ch.logger.Error("send message", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	payload, _ := json.Marshal(msg)
	// If receiver is online over WS, push the message immediately.
	ch.hub.Deliver(req.ReceiverID, payload)

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "message sent successfully", "data": msg})
}

// HandleGetChatHistory godoc
// @Summary      Get chat history
// @Description  Returns all messages exchanged between two users, ordered oldest to newest
// @Tags         chat
// @Produce      json
// @Param        user1_id query string true  "First user ID"
// @Param        user2_id query string true  "Second user ID"
// @Param        page     query int    false "Page number (default 1)"
// @Param        limit    query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /chat/history [get]
func (ch *ChatHandler) HandleGetChatHistory(w http.ResponseWriter, r *http.Request) {
	user1ID := r.URL.Query().Get("user1_id")
	user2ID := r.URL.Query().Get("user2_id")

	if user1ID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "user1_id is required"})
		return
	}
	if user2ID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "user2_id is required"})
		return
	}

	pg := utils.ReadPaginationParams(r)
	messages, err := ch.chatStore.GetChatHistory(user1ID, user2ID, pg.Limit, pg.Offset())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusOK, utils.Envelope{
				"message":    "no messages found",
				"messages":   []models.Message{},
				"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
			})
			return
		}
		ch.logger.Error("get chat history", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if messages == nil {
		messages = []models.Message{}
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "chat history fetched successfully",
		"messages":   messages,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleMarkAsRead godoc
// @Summary      Mark messages as read
// @Description  Marks all unread messages from sender_id to receiver_id as read
// @Tags         chat
// @Produce      json
// @Param        sender_id   query string true "Sender user ID"
// @Param        receiver_id query string true "Receiver user ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /chat/read [put]
func (ch *ChatHandler) HandleMarkAsRead(w http.ResponseWriter, r *http.Request) {
	senderID := r.URL.Query().Get("sender_id")
	receiverID := r.URL.Query().Get("receiver_id")

	if senderID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "sender_id is required"})
		return
	}
	if receiverID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "receiver_id is required"})
		return
	}

	if err := ch.chatStore.MarkAsRead(senderID, receiverID); err != nil {
		ch.logger.Error("mark as read", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "messages marked as read"})
}
