package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dineshappavoo/basex"
	"github.com/ubiq/go-ubiq/common/hexutil"
)

var TimeStampFormat = "2006-01-02T15:04:05.000Z07:00"

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func DateTime() string {
	return time.Now().Format(TimeStampFormat)
}

func BigIntRandom() *big.Int {
	// Max value, a 130-bits integer, i.e 2^130 - 1
	var max *big.Int = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(130), nil)
	// Generate cryptographically strong pseudo-random between [0, max)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return n
}

func BigIntHexString(num *big.Int) string {
	return fmt.Sprintf("0x%x", num)
}

func BigIntFromString(num string) *big.Int {
	i := new(big.Int)
	_, err := fmt.Sscan(num, i)
	if err != nil {
		panic(err)
	}
	return i
}

func BigIntFromHexString(num string) *big.Int {
	a, err := hexutil.DecodeBig(num)
	if err != nil {
		panic(err)
	}
	return a
}

func ShortId() string {
	encoded, err := basex.Encode(BigIntRandom().String())
	if err != nil {
		panic(err)
	}
	return encoded
}

func ToJsonString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func MapOf() map[string]interface{} {
	return make(map[string]interface{})
}
