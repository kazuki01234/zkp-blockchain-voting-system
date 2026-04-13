package crypto

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

func VerifySignature(pubKeyHex, message, sigHex string) bool {
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return false
	}

	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false
	}

	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return false
	}

	sig, err := ecdsa.ParseDERSignature(sigBytes)
	if err != nil {
		return false
	}

	hash := sha256.Sum256([]byte(message))

	return sig.Verify(hash[:], pubKey)
}

func HashPublicKey(pubKey string) string {
	hash := sha256.Sum256([]byte(pubKey))
	return hex.EncodeToString(hash[:])
}