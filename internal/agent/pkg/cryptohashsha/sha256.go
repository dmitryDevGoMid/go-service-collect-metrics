package cryptohashsha

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
)

type HashSha256 interface {
	GetSha256ByData(data []byte) ([]byte, error)
}

type hashSha256 struct {
	cfg *config.Config
}

func NewSha256(cfg *config.Config) HashSha256 {
	return &hashSha256{cfg: cfg}
}

func (s *hashSha256) GetSha256ByData(data []byte) ([]byte, error) {
	err := errors.New("empty cofig key for create sha256")

	if s.cfg.SHA256.Key != "" {
		//fmt.Println(s.cfg.SHA256.Key)
		//dataString := string(data) + s.cfg.SHA256.Key

		h := hmac.New(sha256.New, []byte(s.cfg.SHA256.Key))
		h.Write([]byte(data))
		bs := h.Sum(nil)

		/*h := sha256.New()
		h.Write([]byte(dataString))

		bs := h.Sum(nil)*/

		return bs, nil
	}

	return nil, err
}
