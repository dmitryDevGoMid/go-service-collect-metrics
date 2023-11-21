package hashSha256

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
)

type HashSha256 interface {
	GetSha256ByData(data []byte) ([]byte, error)
	CheckHashSHA256Data(body []byte, hashData string) bool
	CheckHashSHA256Key() bool
}

type hashSha struct {
	cfg *config.Config
}

func NewSha256(cfg *config.Config) HashSha256 {
	return &hashSha{cfg: cfg}
}

func (s *hashSha) GetSha256ByData(data []byte) ([]byte, error) {
	err := errors.New("empty cofig key for create sha256")

	if s.cfg.HashSHA256.Key != "" {
		dataString := string(data) + s.cfg.HashSHA256.Key

		h := sha256.New()
		h.Write([]byte(dataString))

		bs := h.Sum(nil)

		return bs, nil
	}

	return nil, err
}

func (s *hashSha) CheckHashSHA256Key() bool {
	return s.cfg.HashSHA256.Key != ""
}

func (s *hashSha) CheckHashSHA256Data(body []byte, hashData string) bool {
	if s.CheckHashSHA256Key() {

		calc, err := s.GetSha256ByData(body)

		if err != nil {
			fmt.Println("Error calc hash: ", err)
		}
		calcHash := fmt.Sprintf("%x", calc)

		//fmt.Println("req:", hashData)
		//fmt.Println("calc:", calcHash)

		if hashData == calcHash {
			return true
		}
	}

	return false
}
