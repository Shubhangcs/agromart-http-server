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
	"github.com/shubhangcs/agromart-server/internal/tokens"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
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
	return &ChatHandler{chatStore: chatStore, hub: h, logger: logger}
}

// claimsFromCtx extracts the authenticated user's claims from the request context.
func claimsFromCtx(r *http.Request) *tokens.Token {
	claims, _ := r.Context().Value("claims").(*tokens.Token)
	return claims
}

// HandleWebSocket godoc
// @Summary      Real-time chat WebSocket
// @Description  Upgrades the connection to WebSocket. Pass the JWT as ?token=<jwt> query param — most WS clients cannot set Authorization headers during the handshake. Send JSON {"receiver_id":"...","content":"..."}; receive JSON Message objects in real-time.
// @Tags         chat
// @Param        token query string true "JWT access token"
// @Router       /chat/ws [get]
func (ch *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	rawToken := r.URL.Query().Get("token")
	if rawToken == "" {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "missing token"})
		return
	}
	claims, err := tokens.ValidateToken(rawToken)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or expired token"})
		return
	}
	userID := claims.UserID

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

	// Write pump — forwards queued messages to the WebSocket connection.
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

	// Ping ticker — keeps the connection alive.
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

	// Read pump — blocks until the connection is closed.
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
		// Prevent messaging yourself.
		if incoming.ReceiverID == userID {
			continue
		}

		msg := &models.Message{
			SenderID:   userID, // always from JWT, never from client
			ReceiverID: incoming.ReceiverID,
			Content:    incoming.Content,
		}
		if err := ch.chatStore.SaveMessage(msg); err != nil {
			ch.logger.Error("ws save message", "user_id", userID, "error", err)
			continue
		}

		payload, _ := json.Marshal(msg)

		// Push to receiver if they are online.
		ch.hub.Deliver(incoming.ReceiverID, payload)

		// Echo back to sender with the DB-assigned ID and timestamp.
		select {
		case client.Send <- payload:
		default:
		}
	}
}

// HandleSendMessage godoc
// @Summary      Send a message (REST fallback)
// @Description  Sends a text message to another user. The sender identity is taken from the JWT token — the request body only needs receiver_id and content.
// @Tags         chat
// @Accept       json
// @Produce      json
// @Param        body body models.SendMessageRequest true "Message payload"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /chat/send [post]
func (ch *ChatHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r)
	if claims == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}
	senderID := claims.UserID

	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid request payload")
		return
	}
	if err := validator.Validate(&req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if req.ReceiverID == senderID {
		badRequest(w, "cannot send a message to yourself")
		return
	}

	msg := &models.Message{
		SenderID:   senderID, // always from JWT
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}
	if err := ch.chatStore.SaveMessage(msg); err != nil {
		serverError(w, ch.logger, "send message", err)
		return
	}

	payload, _ := json.Marshal(msg)
	// Push to receiver if they are online via WebSocket.
	ch.hub.Deliver(req.ReceiverID, payload)

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "message sent successfully", "data": msg})
}

// HandleGetChatHistory godoc
// @Summary      Get chat history
// @Description  Returns all messages exchanged between the authenticated user and another user (with_user_id). Only the two participants can fetch this conversation.
// @Tags         chat
// @Produce      json
// @Param        with_user_id query string true  "The other participant's user ID"
// @Param        page         query int    false "Page number (default 1)"
// @Param        limit        query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /chat/history [get]
func (ch *ChatHandler) HandleGetChatHistory(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r)
	if claims == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}
	myID := claims.UserID

	withUserID := r.URL.Query().Get("with_user_id")
	if withUserID == "" {
		badRequest(w, "with_user_id is required")
		return
	}

	pg := utils.ReadPaginationParams(r)
	messages, err := ch.chatStore.GetChatHistory(myID, withUserID, pg.Limit, pg.Offset())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusOK, utils.Envelope{
				"messages":   []models.Message{},
				"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
			})
			return
		}
		serverError(w, ch.logger, "get chat history", err)
		return
	}

	if messages == nil {
		messages = []models.Message{}
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"messages":   messages,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleMarkAsRead godoc
// @Summary      Mark messages as read
// @Description  Marks all unread messages from sender_id (the other user) to the authenticated user as read.
// @Tags         chat
// @Produce      json
// @Param        sender_id query string true "The user who sent the messages"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /chat/read [put]
func (ch *ChatHandler) HandleMarkAsRead(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r)
	if claims == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return
	}
	receiverID := claims.UserID // the authenticated user is the one marking their messages as read

	senderID := r.URL.Query().Get("sender_id")
	if senderID == "" {
		badRequest(w, "sender_id is required")
		return
	}

	if err := ch.chatStore.MarkAsRead(senderID, receiverID); err != nil {
		serverError(w, ch.logger, "mark as read", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "messages marked as read"})
}
