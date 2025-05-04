package user

import (
	""
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type SingleFilter struct {
	ID           *uuid.UUID
	Username     *string
	SolanaWallet *string
}

type Filter struct {
	IDs           []uuid.UUID
	Usernames     []string
	SolanaWallets []string
}

type Service interface {
	Register(ctx context.Context, newUser NewUser) error
	FindOne(ctx context.Context, filter SingleFilter) (User, error)
	Find(ctx context.Context, filter Filter) ([]User, error)
}

type AuthService interface {
	GenerateTokenPair(ctx context.Context, creds Credentials) (TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (TokenPair, error)
}

type service struct {
	logger      *slog.Logger
	userStorage Storage
}

func NewService(logger *slog.Logger, userStorage Storage) *service {
	return &service{
		logger:      logger.WithGroup("core-user-service"),
		userStorage: userStorage,
	}
}

// service implements Service and AuthService

type authService struct {
	logger      *slog.Logger
	userStorage Storage
	jwtSecret   []byte
	accessTTL   time.Duration
	refreshTTL  time.Duration
}

// NewAuthService creates a new AuthService with given secret and storage
func NewAuthService(logger *slog.Logger, storage Storage, jwtSecret []byte) AuthService {
	return &authService{
		logger:      logger.WithGroup("auth-service"),
		userStorage: storage,
		jwtSecret:   jwtSecret,
		accessTTL:   15 * time.Minute,
		refreshTTL:  7 * 24 * time.Hour,
	}
}

func (a *authService) GenerateTokenPair(ctx context.Context, creds Credentials) (TokenPair, error) {
	// authenticate user
	user, err := a.userStorage.FindOneByUsername(ctx, creds.Username)
	if err != nil {
		return TokenPair{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		return TokenPair{}, errors.New("invalid credentials")
	}
	// generate tokens
	accessToken, err := a.makeToken(user.ID, a.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, err := a.makeToken(user.ID, a.refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *authService) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	// parse token
	claims := &jwt.RegisteredClaims{}
	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})
	if err != nil || !tkn.Valid {
		return TokenPair{}, errors.New("invalid refresh token")
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return TokenPair{}, err
	}
	// issue new tokens
	accessToken, err := a.makeToken(id, a.accessTTL)
	if err != nil {
		return TokenPair{}, err
	}
	newRefresh, err := a.makeToken(id, a.refreshTTL)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: newRefresh}, nil
}

func (a *authService) makeToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	return claims.SignedString(a.jwtSecret)
}

// Service methods (Register, FindOne, Find)

func (s *service) Register(ctx context.Context, newUser NewUser) error {
	id := uuid.New()
	user := User{
		ID:                    id,
		Username:              newUser.Username,
		SolanaWalletPublicKey: newUser.SolanaWalletPublicKey,
		PasswordHash:          newUser.PasswordHash,
		Games:                 []string{},
	}
	err := s.userStorage.Save(ctx, user)
	if err != nil {
		s.logger.Error("register failed", "error", err)
	}
	return err
}

func (s *service) FindOne(ctx context.Context, filter SingleFilter) (User, error) {
	return s.userStorage.FindOne(ctx, filter)
}

func (s *service) Find(ctx context.Context, filter Filter) ([]User, error) {
	return s.userStorage.Find(ctx, filter)
}
