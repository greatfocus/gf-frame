package server

import (
	"net/http"
	"strings"
	"time"

	gfjwt "github.com/greatfocus/gf-jwt"
	"github.com/greatfocus/gf-sframe/config"
)

// Token struct
type Token struct {
	Role        string
	Permissions []string
	UserID      int64
}

// JWT struct
type JWT struct {
	Secret     string
	Authorized bool
	Minutes    int64
	algorithm  gfjwt.Algorithm
}

// Init method prepare module
func (j *JWT) Init(config *config.Config) {
	j.Secret = config.Server.JWT.Secret
	j.Authorized = config.Server.JWT.Authorized
	j.Minutes = config.Server.JWT.Minutes
	j.algorithm = gfjwt.HmacSha256(j.Secret)
}

// CreateToken generates jwt for API login
func (j *JWT) CreateToken(userID int64, role string, permissions []string) (string, error) {
	claims := gfjwt.NewClaim()
	claims.Set("authorized", j.Authorized)
	claims.Set("userID", userID)
	claims.Set("role", role)
	claims.Set("permissions", permissions)
	claims.Set("exp", time.Now().Add(time.Minute*time.Duration(j.Minutes)).Unix()) //JWT expires after 1 hour
	token, err := j.algorithm.Encode(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

// TokenValid checks for jwt validity
func (j *JWT) TokenValid(r *http.Request) error {
	token := j.extractToken(r)
	err := j.algorithm.Validate(token)
	if err != nil {
		return err
	}
	return nil
}

// extractToken get jwt from header
func (j *JWT) extractToken(r *http.Request) string {
	keys := r.URL.Query()
	jwt := keys.Get("jwt")
	if jwt != "" {
		return jwt
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// GetToken get jwt id from header
func (j *JWT) GetToken(r *http.Request) (Token, error) {
	tokenString := j.extractToken(r)
	claims, err := j.algorithm.Decode(tokenString)
	if err != nil {
		return Token{}, err
	}

	userID, _ := claims.Get("userID")
	role, _ := claims.Get("Role")
	permissions, _ := claims.Get("permissions")
	var token = Token{
		UserID:      userID.(int64),
		Role:        role.(string),
		Permissions: permissions.([]string),
	}

	return token, nil
}
