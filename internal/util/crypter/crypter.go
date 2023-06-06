package crypter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

// New initalizes crypter service
func New() *Service {
	return &Service{}
}

// Service holds crypter methods
type Service struct{}

// HashPassword hashes the password using bcrypt
func (*Service) HashPassword(password string) string {
	return HashPassword(password)
}

// CompareHashAndPassword matches hash with password. Returns true if hash and password match.
func (*Service) CompareHashAndPassword(hash, password string) bool {
	return CompareHashAndPassword(hash, password)
}

// UID returns unique string ID
func (*Service) UID() string {
	return UID()
}

// RoundFloat rounds float64 to 2 decimal places
func (s *Service) RoundFloat(f float64) float64 {
	return toFixedFloat(f, 2)
}

// Float64ToByte converts float64 to byte
func (s *Service) Float64ToByte(f float64) []byte {
	return float64ToByte(f)
}

// NanoID return unique string nano ID
func (s *Service) NanoID() (string, error) {
	// timeT := time.Now().UnixNano()
	// existedID := make(map[string]int64, 0)
	// retried := 0
	// for retried <= MaxRetry {
	// newID, err := generateNanoID()
	// if err != nil {
	// 	return "", err
	// }

	// * check the new id is existed or not
	// if _, ok := existedID[newID]; ok {
	// 	continue
	// }
	// if !ok {
	// 	return newID, nil
	// }

	// * add the new id to existed code
	// existedID[newID] = timeT

	// }

	return generateNanoID()
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

// CompareHashAndPassword matches hash with password. Returns true if hash and password match.
func CompareHashAndPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// UID returns unique string ID
func UID() string {
	return ksuid.New().String()
}

// generateNanoID the new code based on nanoid
func generateNanoID() (string, error) {
	t := time.Now()
	generate, err := nanoid.CustomASCII(ASCII, CodeLen)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", t.Format(DateLayout), generate()), nil
}

func toFixedFloat(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}
