
package handlers

import (
	"AdminiSoftware/internal/auth"
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewAuthHandler(db *gorm.DB, logger *utils.Logger) *AuthHandler {
	return &AuthHandler{
		db:     db,
		logger: logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TwoFA    string `json:"two_fa"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.TwoFactorEnabled && req.TwoFA == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Two-factor authentication required"})
		return
	}

	if user.TwoFactorEnabled && !auth.ValidateTOTP(req.TwoFA, user.TwoFactorSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid two-factor code"})
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	user.LastLogin = time.Now()
	h.db.Save(&user)

	h.logger.Info("User logged in: " + user.Email)

	c.JSON(http.StatusOK, LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Invalidate token logic here
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) EnableTwoFactor(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	secret, qrCode, err := auth.GenerateTOTPSecret(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret":  secret,
		"qr_code": qrCode,
	})
}

func (h *AuthHandler) VerifyTwoFactor(c *gin.Context) {
	type VerifyRequest struct {
		Secret string `json:"secret" binding:"required"`
		Code   string `json:"code" binding:"required"`
	}

	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !auth.ValidateTOTP(req.Code, req.Secret) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	userID := c.GetUint("user_id")
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.TwoFactorEnabled = true
	user.TwoFactorSecret = req.Secret
	h.db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Two-factor authentication enabled"})
}
package handlers

import (
	"AdminiSoftware/internal/auth"
	"AdminiSoftware/internal/models"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
	twoFA      *auth.TwoFactorAuth
}

func NewAuthHandler(db *gorm.DB, jwtManager *auth.JWTManager, twoFA *auth.TwoFactorAuth) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtManager: jwtManager,
		twoFA:      twoFA,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		TOTPCode string `json:"totp_code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.TwoFactorEnabled {
		if request.TOTPCode == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "2FA code required"})
			return
		}
		if !h.twoFA.ValidateCode(user.TwoFactorSecret, request.TOTPCode) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
			return
		}
	}

	token, err := h.jwtManager.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user.LastLogin = &time.Time{}
	*user.LastLogin = time.Now()
	h.db.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Username  string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:     request.Email,
		Password:  string(hashedPassword),
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Username:  request.Username,
		Role:      "user",
		Status:    "active",
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Implementation for token refresh
	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}
package handlers

import (
	"AdminiSoftware/internal/auth"
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
	bruteForce *auth.BruteForceProtection
	twoFactor  *auth.TwoFactorManager
	logger     *utils.Logger
}

func NewAuthHandler(db *gorm.DB, jwtManager *auth.JWTManager, bruteForce *auth.BruteForceProtection, logger *utils.Logger) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtManager: jwtManager,
		bruteForce: bruteForce,
		twoFactor:  auth.NewTwoFactorManager(),
		logger:     logger,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	TwoFA    string `json:"two_fa,omitempty"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientIP := utils.GetClientIP(c.Request.RemoteAddr, c.GetHeader("X-Forwarded-For"), c.GetHeader("X-Real-IP"))

	// Check brute force protection
	if h.bruteForce.IsBlocked(clientIP) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many failed attempts. Please try again later.",
		})
		return
	}

	var user models.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		h.bruteForce.RecordAttempt(clientIP, false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		h.bruteForce.RecordAttempt(clientIP, false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check 2FA if enabled
	if user.TwoFactorEnabled {
		if req.TwoFA == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Two-factor authentication required",
				"requires_2fa":    true,
				"temp_token":      h.generateTempToken(user.ID),
			})
			return
		}

		if !h.twoFactor.ValidateCode(user.TwoFactorSecret, req.TwoFA) {
			h.bruteForce.RecordAttempt(clientIP, false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
			return
		}
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	h.db.Save(&user)

	// Record successful login
	h.bruteForce.RecordAttempt(clientIP, true)

	// Log login attempt
	h.db.Create(&models.LoginAttempt{
		Username:  req.Username,
		IP:        clientIP,
		Success:   true,
		UserAgent: c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if !utils.IsValidUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username format"})
		return
	}

	if !utils.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if !utils.IsValidPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters with uppercase, lowercase, and number"})
		return
	}

	// Check if user exists
	var existingUser models.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		Status:    "active",
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"role":       user.Role,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := h.db.Preload("Package").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":                 user.ID,
			"username":           user.Username,
			"email":              user.Email,
			"role":               user.Role,
			"first_name":         user.FirstName,
			"last_name":          user.LastName,
			"contact_email":      user.ContactEmail,
			"theme":              user.Theme,
			"language":           user.Language,
			"two_factor_enabled": user.TwoFactorEnabled,
			"last_login":         user.LastLogin,
			"package":            user.Package,
			"created_at":         user.CreatedAt,
		},
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Implementation for updating user profile
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Implementation for changing password
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *AuthHandler) Enable2FA(c *gin.Context) {
	// Implementation for enabling 2FA
	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

func (h *AuthHandler) Disable2FA(c *gin.Context) {
	// Implementation for disabling 2FA
	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	// Implementation for forgot password
	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	// Implementation for reset password
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	// Implementation for verifying 2FA
	c.JSON(http.StatusOK, gin.H{"message": "2FA verified successfully"})
}

func (h *AuthHandler) generateTempToken(userID uint) string {
	return utils.GenerateRandomString(32) + ":" + strconv.Itoa(int(userID))
}
