// package auth

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/go-redis/redis/v8"
// 	"github.com/golang-jwt/jwt/v5"
// 	"golang.org/x/crypto/bcrypt"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"

// 	"insta/auth/config"
// 	"insta/auth/database"
// )

// type Auth struct {
// 	config *config.Config
// 	db     *database.Database
// }

// func NewAuth(config *config.Config, db *database.Database) *Auth {
// 	return &Auth{config: config, db: db}
// }

// func (a *Auth) SendVerificationEmail(ctx context.Context, email, code string) error {
// 	// implementation
// }

// func (a *Auth) GenerateJWT(email string, ttl int) (string, error) {
// 	// implementation
// }

// func (a *Auth) VerifyJWT(tokenString string) (string, error) {
// 	// implementation
// }

// func (a *Auth) GenerateSignature() (string, error) {
// 	// implementation
// }