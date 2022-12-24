package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userID string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		issuer: "kuraohikari",
		secretKey: getSecretKey(),
	}
}

func getSecretKey() string {
	key := os.Getenv("JWT_SECRET")
	return key
}

func (j *jwtService) GenerateToken(UserID string) string {
	claims := &jwtCustomClaim{
		UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer: j.issuer,
		},
		// jwt.StandardClaims{
		// 	ExpiresAt: time.Now().AddDate(1,0,0).Unix(),
		// 	Issuer: j.issuer,
		// 	IssuedAt: time.Now().Unix(),
		// },
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error){
	return jwt.Parse(token, func(t_ *jwt.Token)(interface{}, error){
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			// return nil, errors.New("that's not even a token")
			return nil, fmt.Errorf("unexpected signing method %v",  t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	// tokenRes , err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(j.secretKey), nil
	// })
	// if tokenRes.Valid {
	// 	return tokenRes, nil
	// }else if errors.Is(err, jwt.ErrTokenMalformed) {
	// 	return nil, errors.New("that's not even a token")
	// } else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
	// 	// Token is either expired or not active yet
	// 	return nil, errors.New("timing is everything")
	// 	// fmt.Println("Timing is everything")
	// } else {
	// 	return nil,  errors.New("couldn't handle this token")
	// 	// fmt.Println("Couldn't handle this token:", err)
	// }
	// return nil, fmt.Errorf("That's not even a token")
}