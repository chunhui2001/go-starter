package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	gerror "github.com/go-errors/errors"
	"github.com/ubiq/go-ubiq/common/hexutil"
)

// 2006-01-02T15:04:05.999Z
var TimeStampFormat = "2006-01-02T15:04:05.000Z07:00"

func RootDir2() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func RootDir3() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func RootDir() string {
	dir, _ := os.Getwd()
	return dir
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

func DateTimeUTCString() string {
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

func randint64() (int64, error) {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))
	if err != nil {
		return 0, err
	}
	return val.Int64(), nil
}

func ShortId() string {

	// encoded, _ := basex.Encode(BigIntRandom().String())

	val, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))

	if err != nil {
		panic(err)
	}

	b := make([]byte, 8)

	binary.LittleEndian.PutUint64(b, uint64(val.Int64()))
	encoded := base64.StdEncoding.EncodeToString([]byte(b))

	var replacer = strings.NewReplacer(
		"+", "-",
		"/", "_",
	)

	return replacer.Replace(encoded[:11])

}

func ToJsonString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func AsMap(buf []byte) map[string]interface{} {
	var m map[string]interface{}
	if err := json.Unmarshal(buf, &m); err != nil {
		panic(err)
	}
	return m
}

func ToString(s any) string {
	if reflect.TypeOf(s).String() == "string" {
		return fmt.Sprintf("%s", s)
	}
	if reflect.TypeOf(s).String() == "bool" {
		return fmt.Sprintf("%t", s)
	}
	return fmt.Sprintf("%d", s)
}

func ErrorToString(err interface{}) string {
	return gerror.Wrap(err, 2).ErrorStack()
}

func MapOf(kv ...any) map[string]interface{} {

	if kv == nil {
		return nil
	}

	if len(kv)%2 != 0 {
		panic(errors.New("Invalid map size: currentSize=" + ToString(len(kv))))
	}

	m := make(map[string]interface{})

	for i := 0; i < len(kv); i++ {
		k := ToString(kv[i])
		m[k] = kv[i+1]
		i++
	}

	return m

}

func StringToBytes(str string) []byte {
	return []byte(str)
}

func TempDir() string {
	return os.TempDir()
}

// https://github.com/git-time-metric/gtm/blob/master/util/string.go
func PadLeft(s string, padStr string, maxLen int) string {
	var padCountInt = 1 + ((maxLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - maxLen):]
}

func TrimRight(s string) string {
	return strings.TrimSuffix(s, ",")
}

func Split(s string, sep string) []string {
	return strings.Split(s, sep)
}

func Matches(s string, regx string) [][]string {
	re := regexp.MustCompile(regx)
	return re.FindAllStringSubmatch(s, -1)
}

// params := getParams(`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})`, `2015-05-27`)
// fmt.Println(params)
// ### and the output will be:
// map[Year:2015 Month:05 Day:27]
func MatchesGroup(regEx, str string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(str)

	paramsMap = make(map[string]string)

	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	return paramsMap

}

func IfNull(obj any, defaultValue interface{}) interface{} {
	if obj == nil {
		return defaultValue
	}
	return obj
}

func IfElse(b bool, obj any, defaultValue any) any {
	if b {
		return defaultValue
	}
	return obj
}
