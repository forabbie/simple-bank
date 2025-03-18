package api

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	db "github.com/forabbie/vank-app/database/sqlc"
	"github.com/forabbie/vank-app/mail"
	"github.com/forabbie/vank-app/token"
	"github.com/forabbie/vank-app/util"
	"github.com/forabbie/vank-app/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
type TaskProcessor struct {
	mailer mail.EmailSender
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req db.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}

	// Validate request fields
	violations := validator.ValidateCreateUserRequest(&req)
	if len(violations) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": violations,
		})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, gin.H{
					"error":   "validation failed",
					"details": gin.H{"error": "username or email already exists"},
				})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req db.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.FormatValidationErrors(err))
		return
	}
	req.Username = strings.ToLower(req.Username)

	// Validate request fields
	violations := validator.ValidateLoginUserRequest(&req)
	if len(violations) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": violations,
		})
		return
	}

	user, err := server.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid username",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Verify password
	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	// Generate tokens
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	// Create session
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	})
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req db.UpdateUserRequest

	// Bind the ID from the URL and the rest from JSON
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	// Fetch the existing user to verify ownership
	existingUser, err := server.store.GetUserByUsername(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Ensure the authenticated user can only update their own account
	if existingUser.ID != req.ID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var hashedPassword sql.NullString
	if req.Password != nil && *req.Password != "" {
		hashed, err := util.HashPassword(*req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		hashedPassword = sql.NullString{String: hashed, Valid: true}
	} else {
		hashedPassword = sql.NullString{Valid: false} // Keeps existing password
	}

	// Convert FullName and Email to sql.NullString
	fullName := sql.NullString{}
	if req.FullName != nil {
		fullName = sql.NullString{String: *req.FullName, Valid: true}
	}

	email := sql.NullString{}
	if req.Email != nil {
		email = sql.NullString{String: *req.Email, Valid: true}
	}

	arg := db.UpdateUserParams{
		ID:             req.ID,
		HashedPassword: hashedPassword,
		FullName:       fullName,
		Email:          email,
	}

	// Perform the update
	updatedUser, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusConflict, gin.H{
					"error":   "validation failed",
					"details": gin.H{"error": "username or email already exists"},
				})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	rsp := newUserResponse(updatedUser)
	ctx.JSON(http.StatusOK, rsp)
}
