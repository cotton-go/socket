package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Ase 结构体
type Ase struct {
	key []byte
}

// NewAse 创建一个新的 Ase 实例
func NewAse(key string) (*Ase, error) {
	if key == "" {
		newKey, _ := generateRandomKey(32)
		key = string(newKey)
	}

	// 你可能需要处理 key 的长度和格式，这里简单地使用 key 的字节数组
	return &Ase{key: []byte(key)}, nil
}

// GenerateRandomKey 生成指定长度的随机字节
func generateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt 使用 AES 加密给定的数据
func (sc Ase) encrypt(plaintext string) (any, error) {
	block, err := aes.NewCipher(sc.key)
	if err != nil {
		return "", err
	}

	// 使用随机生成的 IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用 AES 解密给定的数据
func (sc Ase) decrypt(ciphertext string) (any, error) {
	block, err := aes.NewCipher(sc.key)
	if err != nil {
		return "", err
	}

	// 解码 base64
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 提取 IV
	iv := decoded[:aes.BlockSize]
	decoded = decoded[aes.BlockSize:]

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decoded, decoded)

	return string(decoded), nil
}

func (sc Ase) Encode(value any) (any, error) {
	var data string

	switch v := value.(type) {
	case string:
		data = v
	case []byte:
		data = string(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		data = fmt.Sprintf("%t", v)
	default:
		return "", fmt.Errorf("Unsupported data type: %T", value)
	}

	return sc.encrypt(data)
}

func (sc Ase) Decode(value any) (any, error) {
	var data string
	return sc.decrypt(data)
}
