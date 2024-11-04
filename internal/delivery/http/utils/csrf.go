package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "io"
    "time"
    "errors"

    "github.com/google/uuid"
	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"go.uber.org/zap"
)

type CryptToken struct {
    Secret []byte
}

type TokenData struct {
    SessionID uuid.UUID
    UserID    uuid.UUID
    Exp       time.Time
}

var (
	ErrSecretKeyNotSet = errors.New("CSRF secret is not set")
	ErrSecretKeyInvalid = errors.New("secret key must be 32 bytes for AES-256")
	ErrMarshalTokenData = errors.New("failed to marshal token data")
	ErrDecodeToken = errors.New("failed to decode base64 token")
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	ErrDecryptionFailed = errors.New("decryption failed")
	ErrInvalidTokenDataFormat = errors.New("invalid token data format")
	ErrTokenExpired = errors.New("token expired")
)

func NewAesCryptHashToken(logger *zap.Logger) (*CryptToken, error) {
    secret, err := getSecretKey()
    if err != nil {
        logger.Error("failed to load secret: %v", zap.Error(err))
        return nil, ErrSecretKeyNotSet
    }

    if len(secret) != 32 {
        logger.Error("secret key must be 32 bytes for AES-256")
        return nil, ErrSecretKeyInvalid
    }

    return &CryptToken{Secret: secret}, nil
}

func getSecretKey() ([]byte, error) {
    secret := config.GetCSRFSecret()
    if len(secret) == 0 {
        return nil, ErrSecretKeyNotSet
    }
    
    key := []byte(secret)
    if len(key) != 32 {
        return nil, ErrSecretKeyInvalid
    }
    
    return key, nil
}

func (tk *CryptToken) Create(sessionId uuid.UUID, userId uuid.UUID, tokenExpTime int64) (string, error) {
    block, err := aes.NewCipher(tk.Secret)
    if err != nil {
        return "", err
    }

    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, aesgcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    td := &TokenData{SessionID: sessionId, UserID: userId, Exp: time.Unix(tokenExpTime, 0)}
    data, err := json.Marshal(td)
    if err != nil {
        return "", ErrMarshalTokenData
    }

    ciphertext := aesgcm.Seal(nil, nonce, data, nil)
    res := append(nonce, ciphertext...)
    return base64.URLEncoding.EncodeToString(res), nil
}

func (tk *CryptToken) Check(sessionId uuid.UUID, userId uuid.UUID, inputToken string) (bool, error) {
    block, err := aes.NewCipher(tk.Secret)
    if err != nil {
        return false, err
    }
    aesgcm, err := cipher.NewGCM(block)
    if err != nil {
        return false, err
    }

    ciphertext, err := base64.URLEncoding.DecodeString(inputToken)
    if err != nil {
        return false, ErrDecodeToken
    }

    nonceSize := aesgcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return false, ErrCiphertextTooShort
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return false, ErrDecryptionFailed
    }

    td := TokenData{}
    err = json.Unmarshal(plaintext, &td)
    if err != nil {
        return false, ErrInvalidTokenDataFormat
    }

    if td.Exp.Before(time.Now()) {
        return false, ErrTokenExpired
    }

    return sessionId == td.SessionID && userId == td.UserID, nil
}
