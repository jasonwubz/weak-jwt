package handler

import (
	"crypto"
	"crypto/hmac"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// Implements the CUSTOM HMAC-SHA256 signing methods signing methods
// Expects key type of []byte for both signing and validation
type SigningMethodHMAC struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for HS256 and company
var (
	SigningMethodCS256  *SigningMethodHMAC
	ErrSignatureInvalid = errors.New("signature is invalid")
)

func init() {
	// HS256
	SigningMethodCS256 = &SigningMethodHMAC{"CS256", crypto.SHA256}
	jwt.RegisterSigningMethod(SigningMethodCS256.Alg(), func() jwt.SigningMethod {
		return SigningMethodCS256
	})

}

func (m *SigningMethodHMAC) Alg() string {
	return m.Name
}

// Verify the signature of HSXXX tokens.  Returns nil if the signature is valid.
func (m *SigningMethodHMAC) Verify(signingString, signature string, key interface{}) error {
	// Verify the key is the right type
	keyBytes, ok := key.([]byte)
	if !ok {
		return jwt.ErrInvalidKeyType
	}

	// Decode signature, for comparison
	sig, err := jwt.DecodeSegment(signature)
	if err != nil {
		return err
	}

	// Can we use the specified hashing method?
	if !m.Hash.Available() {
		return jwt.ErrHashUnavailable
	}

	// This signing method is symmetric, so we validate the signature
	// by reproducing the signature from the signing string and key, then
	// comparing that against the provided signature.
	hasher := hmac.New(m.Hash.New, keyBytes)
	hasher.Write([]byte(signingString))
	hashedStr := hasher.Sum(nil)
	// fmt.Printf("%x %d\n", hashedStr, len(hashedStr))
	// fmt.Printf("%x %d\n", sig, len(sig))
	x := fmt.Sprintf("%x", hashedStr)
	y := fmt.Sprintf("%x", sig)
	if unsafeCompare(x, y) == 0 {
		return ErrSignatureInvalid
	}

	//fmt.Println("no validation error")

	// No validation errors.  Signature is good.
	return nil
}

func unsafeCompare(x string, y string) (result int) {
	if len(x) != len(y) {
		return 0
	}
	for idx, i := range x {
		if string(i) != y[idx:idx+1] {
			return 0
		}
	}
	return 1
}

// Implements the Sign method from SigningMethod for this signing method.
// Key must be []byte
func (m *SigningMethodHMAC) Sign(signingString string, key interface{}) (string, error) {
	if keyBytes, ok := key.([]byte); ok {
		if !m.Hash.Available() {
			return "", jwt.ErrHashUnavailable
		}

		hasher := hmac.New(m.Hash.New, keyBytes)
		hasher.Write([]byte(signingString))

		return jwt.EncodeSegment(hasher.Sum(nil)), nil
	}

	return "", jwt.ErrInvalidKeyType
}
