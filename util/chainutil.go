package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

// ChainUtil provides cryptographic utility functions
type ChainUtil struct{}

// GenKeyPair generates an elliptic curve key pair (private and public keys)
func GenKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// ID generates a unique identifier using UUID
func ID() string {
	return uuid.New().String()
}

// Hash creates a SHA256 hash from the given data
func ChainUtilHash(data interface{}) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", data)))
	return hex.EncodeToString(hash.Sum(nil))
}

// VerifySignature verifies if the given signature is valid for the provided data hash
func VerifySignature(publicKey *ecdsa.PublicKey, signature []byte, dataHash string) bool {
	hashBytes := []byte(dataHash)
	r, s := new(big.Int), new(big.Int)
	signatureLen := len(signature)
	r.SetBytes(signature[:(signatureLen / 2)])
	s.SetBytes(signature[(signatureLen / 2):])

	return ecdsa.Verify(publicKey, hashBytes, r, s)
}
