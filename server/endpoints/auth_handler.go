package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"opentunnel/server/config"
	"opentunnel/server/service"
)

type AuthHandler struct {
	authService       *service.AuthService
	cliSessionService *service.CLISessionService
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService:       service.NewAuthService(cfg),
		cliSessionService: service.GetCLISessionService(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "username and password are required"})
		return
	}

	token, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		Message: "login successful",
	})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing authorization token"})
		return
	}

	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"username": claims.Username,
	})
}

func (h *AuthHandler) StartCLIAuth(c *gin.Context) {
	sessionID, err := h.cliSessionService.CreateSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create auth session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
	})
}

func (h *AuthHandler) PollCLIAuth(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "missing session_id"})
		return
	}

	token, pending, valid := h.cliSessionService.PollSession(sessionID)
	if !valid {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "session not found or expired"})
		return
	}

	if pending {
		c.JSON(http.StatusAccepted, gin.H{"status": "pending"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		Message: "login successful",
	})
}

func (h *AuthHandler) CompleteCLIAuth(c *gin.Context) {
	sessionID := c.Param("session_id")
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "token is required"})
		return
	}

	if ok := h.cliSessionService.CompleteSession(sessionID, req.Token); !ok {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "session not found or expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cli auth completed"})
}
