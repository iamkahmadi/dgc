package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

// SignData signs the given data hash using the provided private key
func SignData(privKey *ecdsa.PrivateKey, dataHash string) ([]byte, error) {
	hashBytes := []byte(dataHash)

	// Generate signature (r, s values)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hashBytes)
	if err != nil {
		return nil, err
	}

	// Convert r, s values to bytes and concatenate them
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// PublicKeyToHex converts an ECDSA public key to a hexadecimal string
func PublicKeyToHex(pubKey *ecdsa.PublicKey) string {
	if pubKey == nil {
		return ""
	}

	// Create a byte slice to hold the X and Y coordinates
	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	return hex.EncodeToString(pubKeyBytes)
}

// HexToPublicKey converts a hexadecimal string to an ECDSA public key
func HexToPublicKey(hexString string) (*ecdsa.PublicKey, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	if len(bytes) != 64 { // Ensure we have 64 bytes (32 bytes for X and 32 bytes for Y)
		return nil, errors.New("invalid public key length")
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(), // Change this if using a different curve
		X:     new(big.Int).SetBytes(bytes[:32]),
		Y:     new(big.Int).SetBytes(bytes[32:]),
	}

	return pubKey, nil
}
