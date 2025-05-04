package user

import "github.com/google/uuid"

type User struct {
	ID                    uuid.UUID
	Username              string
	SolanaWalletPublicKey string
	PasswordHash          string
	Games                 []string // list of games player's interested in
}

type NewUser struct {
	Username              string
	SolanaWalletPublicKey string
	PasswordHash          string
}

type Credentials struct {
	Username string
	Password string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
