package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/luismingati/buymeacoffee/internal/service"
	"github.com/luismingati/buymeacoffee/internal/store/pg"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"password"`
}

func (sr *SignupRequest) Validate() error {
	if strings.TrimSpace(sr.Username) == "" {
		return errors.New("username é obrigatório")
	}
	if len(sr.Username) < 3 {
		return errors.New("o username deve ter pelo menos 3 caracteres")
	}
	if strings.TrimSpace(sr.Email) == "" {
		return errors.New("email é obrigatório")
	}
	if _, err := mail.ParseAddress(sr.Email); err != nil {
		return errors.New("email inválido")
	}
	if len(sr.Password) < 8 {
		return errors.New("a senha deve ter pelo menos 8 caracteres")
	}
	return nil
}

func (sr *SigninRequest) Validate() error {
	if strings.TrimSpace(sr.Email) == "" {
		return errors.New("email é obrigatório")
	}
	if _, err := mail.ParseAddress(sr.Email); err != nil {
		return errors.New("email inválido")
	}
	if strings.TrimSpace(sr.Password) == "" {
		return errors.New("senha é obrigatória")
	}
	if len(sr.Password) < 8 {
		return errors.New("a senha deve ter pelo menos 8 caracteres")
	}
	return nil
}

func (fr *ForgotPasswordRequest) Validate() error {
	if strings.TrimSpace(fr.Email) == "" {
		return errors.New("email é obrigatório")
	}
	if _, err := mail.ParseAddress(fr.Email); err != nil {
		return errors.New("Por favor, informe um email válido.")
	}
	return nil
}

func (rr *ResetPasswordRequest) Validate() error {
	if strings.TrimSpace(rr.Token) == "" {
		return errors.New("token é obrigatório")
	}
	if strings.TrimSpace(rr.NewPassword) == "" {
		return errors.New("nova senha é obrigatória")
	}
	if len(rr.NewPassword) < 8 {
		return errors.New("a senha deve ter pelo menos 8 caracteres")
	}
	return nil
}

func (cfg *apiConfig) SignupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body SignupRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error("Failed to decode request body: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	if err := body.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	hashedPwd, err := service.HashPassword(body.Password)
	if err != nil {
		slog.Error("Failed to hash password: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	userId, err := cfg.q.InsertUser(r.Context(), pg.InsertUserParams{
		ID:       uuid.New(),
		Username: body.Username,
		Email:    body.Email,
		Password: hashedPwd,
	})
	if err != nil {
		slog.Error("Failed to insert user: ", err.Error(), err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				respondWithError(w, http.StatusConflict, "Usuário ou email já cadastrados.", "Por favor, tente outro nome de usuário.")
				return
			}
		}
		respondWithInternalServerError(w)
		return
	}

	claims := jwt.MapClaims{
		"user_id":  userId.String(),
		"username": body.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token, err := cfg.jwt.Sign(claims)
	if err != nil {
		slog.Error("Failed to sign token: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (cfg *apiConfig) SigninHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Failed to decode request body: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	if err := body.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	user, err := cfg.q.FindUserByEmail(r.Context(), body.Email)
	if err != nil {
		slog.Error("Failed to find user: ", err.Error(), err)
		respondWithError(w, http.StatusUnauthorized, "Credenciais inválidas", "Usuário ou senha incorretos.")
		return
	}

	if !service.CheckPasswordHash(body.Password, user.Password) {
		respondWithError(w, http.StatusUnauthorized, "Credenciais inválidas", "Usuário ou senha incorretos.")
		return
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token, err := cfg.jwt.Sign(claims)
	if err != nil {
		slog.Error("Failed to sign token: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (cfg *apiConfig) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Failed to decode request body: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	if err := body.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Dados iválidos.", err.Error())
		return
	}

	user, err := cfg.q.FindUserByEmail(r.Context(), body.Email)
	if err != nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Se o email existir, um link de recuperação será enviado."})
		return
	}

	resetToken := uuid.New().String()
	err = cfg.redis.Set(r.Context(), "reset_token:"+resetToken, user.Email, 7*time.Minute)
	if err != nil {
		slog.Error("Failed to save reset token in Redis: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	resetLink := "https://example.com/reset-password?token=" + resetToken

	err = cfg.m.SendEmail(user.Email, "Redefinição de Senha", "Clique no link para redefinir sua senha: "+resetLink)
	if err != nil {
		slog.Error("Failed to send reset email: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Se o email existir, um link de recuperação será enviado."})
}

func (cfg *apiConfig) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error("Failed to decode request body: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	if err := body.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Dados inválidos", err.Error())
		return
	}

	email, err := cfg.redis.Get(r.Context(), "reset_token:"+body.Token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Token inválido ou expirado", "Solicite um novo link de recuperação.")
		return
	}

	hashedPwd, err := service.HashPassword(body.NewPassword)
	if err != nil {
		slog.Error("Failed to hash password: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	err = cfg.q.UpdateUserPassword(r.Context(), pg.UpdateUserPasswordParams{
		Password: hashedPwd,
		Email:    email,
	})
	if err != nil {
		slog.Error("Failed to update password: ", err.Error(), err)
		respondWithInternalServerError(w)
		return
	}

	if err := cfg.redis.Del(r.Context(), "reset_token:"+body.Token); err != nil {
		slog.Error("Failed to delete reset token from Redis: ", err.Error(), err)
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Senha redefinida com sucesso."})
}
