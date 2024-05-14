package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var (
	keyName       string
	keySecret     string
	requestMethod = "GET"
	requestHost   = "api.coinbase.com"
	requestPath   = "/api/v3/brokerage/accounts"
	serviceName   = "retail_rest_api_proxy"
)

func init() {
	var ok bool

	keyName, ok = os.LookupEnv("COINBASE_ORGANIZATION")
	if !ok {
		fmt.Println("COINBASE_ORGANIZATION environment variable not set.")
	}

	keySecret = os.Getenv("COINBASE_PRIVATE_KEY")
	// Required when grabbing key as envronment variable because of formatting
	keySecret = strings.ReplaceAll(keySecret, "\\n", "\n") // Replace "\\n" with "\n"

	if keySecret == "" {
		fmt.Println("COINBASE_PRIVATE_KEY environment variable not set.")
	}
}

type APIKeyClaims struct {
	*jwt.Claims
	URI string `json:"uri"`
}

func generateJtw() (string, error) {
	uri := fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)
	block, _ := pem.Decode([]byte(keySecret))
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing EC private key: %w", err)
	}

	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: key},
		(&jose.SignerOptions{NonceSource: nonceSource{}}).WithType("JWT").WithHeader("kid", keyName),
	)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	cl := &APIKeyClaims{
		Claims: &jwt.Claims{
			Subject:   keyName,
			Issuer:    "coinbase-cloud",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Expiry:    jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			Audience:  jwt.Audience{serviceName},
		},
		URI: uri,
	}
	jwtString, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	return jwtString, nil
}

var max = big.NewInt(math.MaxInt64)

type nonceSource struct{}

func (n nonceSource) Nonce() (string, error) {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return r.String(), nil
}

func printJtw() {
	jwt, err := generateJtw()

	if err != nil {
		fmt.Println("error building jwt: ", err)
	}
	fmt.Println("export JWT=" + jwt)
}
