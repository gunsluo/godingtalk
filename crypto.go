package dingtalk

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	mathrand "math/rand"
	"sort"
	"strings"
	"time"
)

type PushCryptoSuit struct {
	token    string
	aesKey   string
	suiteKey string
	bKey     []byte
	block    cipher.Block
}

func NewPushCryptoSuit(token, aesKey, suiteKey string) (*PushCryptoSuit, error) {
	if len(aesKey) != aesEncodeKeyLen {
		return nil, errors.New("invalid aes key")
	}
	bkey, err := base64.StdEncoding.DecodeString(aesKey + "=")
	if err != nil {
		return nil, errors.New("failed to base64 aes key")
	}
	block, err := aes.NewCipher(bkey)
	if err != nil {
		return nil, errors.New("failed to new cipher")
	}
	suit := &PushCryptoSuit{
		token:    token,
		aesKey:   aesKey,
		suiteKey: suiteKey,
		bKey:     bkey,
		block:    block,
	}
	return suit, nil
}

func (s *PushCryptoSuit) Decrypt(signature, timestamp, nonce, secretMsg string) (string, error) {
	if !s.verificationSignature(s.token, timestamp, nonce, secretMsg, signature) {
		return "", errors.New("ERROR: 签名不匹配")
	}
	decode, err := base64.StdEncoding.DecodeString(secretMsg)
	if err != nil {
		return "", err
	}
	if len(decode) < aes.BlockSize {
		return "", errors.New("ERROR: 密文太短")
	}
	blockMode := cipher.NewCBCDecrypter(s.block, s.bKey[:s.block.BlockSize()])
	plantText := make([]byte, len(decode))
	blockMode.CryptBlocks(plantText, decode)
	plantText = pkCS7UnPadding(plantText)
	size := binary.BigEndian.Uint32(plantText[16:20])
	plantText = plantText[20:]
	corpID := plantText[size:]
	if string(corpID) != s.suiteKey {
		return "", errors.New("ERROR: CorpID匹配不正确")
	}
	return string(plantText[:size]), nil
}

func (s *PushCryptoSuit) Encrypt(msg, timestamp, nonce string) (string, string, error) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(msg)))
	msg = randomString(16) + string(size) + msg + s.suiteKey
	plantText := pkCS7Padding([]byte(msg), s.block.BlockSize())
	if len(plantText)%aes.BlockSize != 0 {
		return "", "", errors.New("ERROR: 消息体size不为16的倍数")
	}
	blockMode := cipher.NewCBCEncrypter(s.block, s.bKey[:s.block.BlockSize()])
	chipherText := make([]byte, len(plantText))
	blockMode.CryptBlocks(chipherText, plantText)
	outMsg := base64.StdEncoding.EncodeToString(chipherText)
	signature := s.createSignature(s.token, timestamp, nonce, string(outMsg))
	return string(outMsg), signature, nil
}

// 验证数据签名
func (s *PushCryptoSuit) verificationSignature(token, timestamp, nonce, msg, sigture string) bool {
	return s.createSignature(token, timestamp, nonce, msg) == sigture
}

// 数据签名
func (s *PushCryptoSuit) createSignature(token, timestamp, nonce, msg string) string {
	params := make([]string, 0)
	params = append(params, token)
	params = append(params, timestamp)
	params = append(params, nonce)
	params = append(params, msg)
	sort.Strings(params)
	return sha1Sign(strings.Join(params, ""))
}

// 解密补位
func pkCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

// 加密补位
func pkCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func sha1Sign(s string) string {
	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(s))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	return fmt.Sprintf("%x", bs)
}

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// 随机字符串
func randomString(n int, alphabets ...byte) string {
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		mathrand.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[mathrand.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[mathrand.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return string(bytes)
}
