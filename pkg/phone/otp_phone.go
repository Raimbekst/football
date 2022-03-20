package phone

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

type SecretCodeGenerator interface {
	GetRandNum() (string, error)
}
type SecretGenerator struct{}

func NewSecretGenerator() SecretGenerator {
	return SecretGenerator{}
}

func (g *SecretGenerator) GetRandNum() (string, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(8999))

	if err != nil {
		return "", fmt.Errorf("phone.GetRandNum: %w", err)
	}
	return strconv.FormatInt(nBig.Int64()+1000, 10), nil
}
