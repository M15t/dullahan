package crypter

import (
	"fmt"
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
