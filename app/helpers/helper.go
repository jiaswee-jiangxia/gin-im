package helpers

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"goskeleton/app/global/variable"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	MobileNo string `json:"mobile_no"`
	jwt.StandardClaims
}

// GetUnixTimestamp 转换所有time.Time格式去Unix Timestamp
func GetUnixTimestamp(t time.Time) int64 {
	return t.Unix()
}

func ParseToken(str string) (*Claims, string) {
	var claimsObj Claims
	tokenString := str
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return variable.PublicKey, nil
	})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return nil, "token_expired"
		}
		return nil, "invalid_token"
	}
	text, _ := json.Marshal(claims)
	err = json.Unmarshal(text, &claimsObj)
	if err != nil {
		return nil, "decode_token_failed"
	}
	return &claimsObj, ""
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func EncryptAES(key []byte, plaintext string) string {
	// create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return err.Error()
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, []byte(plaintext))
	// return hex string
	return hex.EncodeToString(out)
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
