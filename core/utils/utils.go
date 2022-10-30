package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"mime/multipart"
	"net"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	gerror "github.com/go-errors/errors"
	_ "github.com/kardianos/osext"
	"github.com/shopspring/decimal"
	"github.com/ubiq/go-ubiq/common/hexutil"
	"golang.org/x/exp/slices"
)

// 2006-01-02T15:04:05.999Z
var TimeStampFormat = "2006-01-02T15:04:05.000Z07:00"

func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func OutboundIP() net.IP {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func Ip2Int(ip net.IP) int64 {
	i := big.NewInt(0)
	i.SetBytes(ip)
	return i.Int64()
}

func ParseIp(ipv4 string) net.IP {
	return net.ParseIP(ipv4).To4()
}

func RootDir() string {

	var appRoot string = os.Getenv("APP_ROOT")

	if appRoot != "" {
		return appRoot
	}

	dir, _ := os.Getwd()
	return dir

	// folderPath, _ := osext.ExecutableFolder()
	// return folderPath
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

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// https://gist.github.com/josephspurrier/90e957f1277964f26852
func GetFileMd5(file multipart.File) (md5Str string) {

	h := md5.New()

	if _, err := file.Seek(0, 0); err != nil {
		panic(err)
	}

	if _, err := io.Copy(h, file); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x\n", h.Sum(nil))

}

func DateTimeParse(s string) time.Time {
	if t, err := time.Parse(TimeStampFormat, s); err == nil {
		return t
	} else {
		panic(err)
	}
}

func DateTimeUTCString() string {
	return time.Now().Format(TimeStampFormat)
}

func ToDateTimeUTCString(tm time.Time) string {
	return tm.Format(TimeStampFormat)
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

func ToBase64String(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func FromBase64String(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func ToJsonString(v interface{}) string {

	if v == nil {
		return ""
	}

	if reflect.TypeOf(v).String() == "string" {
		return v.(string)
	}

	b, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func ToJsonBytes(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func AsMap(buf []byte) map[string]interface{} {
	var m map[string]interface{}
	if err := json.Unmarshal(buf, &m); err != nil {
		panic(err)
	}
	return m
}

func ToMap(v interface{}) map[string]interface{} {
	return AsMap(ToJsonBytes(v))
}

func StrToInt(str string) int {
	intVar, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return intVar
}

func StrToBool(str string) bool {
	boolVal, err := strconv.ParseBool(str)
	if err != nil {
		panic(err)
	}
	return boolVal
}

func ArrayContains(array []string, val string) bool {
	return slices.Contains(array, val)
}

func ToString(s any) string {

	switch s.(type) {
	case float64, float32:
		return fmt.Sprintf("%f", s)
	case string:
		return fmt.Sprintf("%s", s)
	case bool:
		return fmt.Sprintf("%t", s)
	case byte:
		return fmt.Sprintf("%x", s)
	case []uint8:
		return string(s.([]byte))
	default:
		return fmt.Sprintf("%d", s)
	}
}

func ErrorToString(err interface{}) string {
	return gerror.Wrap(err, 2).ErrorStack()
}

func OfMap(kv ...string) map[string]string {

	if kv == nil {
		return make(map[string]string)
	}

	if len(kv)%2 != 0 {
		panic(errors.New("Invalid map size: currentSize=" + ToString(len(kv))))
	}

	m := make(map[string]string)

	for i := 0; i < len(kv); i++ {
		k := ToString(kv[i])
		m[k] = kv[i+1]
		i++
	}

	return m

}

func MapOf(kv ...any) map[string]interface{} {

	if kv == nil {
		return make(map[string]interface{})
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

func MapsOf(kv ...any) *map[string]interface{} {

	m := make(map[string]interface{})

	if kv == nil {
		return &m
	}

	if len(kv)%2 != 0 {
		panic(errors.New("Invalid map size: currentSize=" + ToString(len(kv))))
	}

	for i := 0; i < len(kv); i++ {
		k := ToString(kv[i])
		m[k] = kv[i+1]
		i++
	}

	return &m

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

func Lower(s string) string {
	return strings.ToLower(s)
}

func Join(delim string, s ...any) string {
	return strings.Trim(strings.Replace(fmt.Sprint(s...), " ", delim, -1), "[]")
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
	if obj == (*any)(nil) {
		return defaultValue
	}
	if obj == (*string)(nil) {
		return defaultValue
	}
	return obj
}

func IfElse(b bool, obj any, defaultValue any) any {
	if b {
		return obj
	}
	return defaultValue
}

func GoroutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Base64UUID() string {

	b := make([]byte, 16)

	_, err := rand.Read(b)

	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		panic(err)
	}

	b[6] &= 0x0f /* clear the 4 most significant bits for the version  */
	b[6] |= 0x40 /* set the version to 0100 / 0x40 */

	/* Set the variant:
	 * The high field of th clock sequence multiplexed with the variant.
	 * We set only the MSB of the variant*/
	b[8] &= 0x3f /* clear the 2 most significant bits */
	b[8] |= 0x80 /* set the variant (MSB is set)*/

	return base64.RawURLEncoding.EncodeToString(b)

}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func HumanFileSizeWithInt(size int) string {
	return HumanFileSize(float64(size))
}

func HumanFileSize(size float64) string {

	if size < 0 {
		return "0"
	}

	var suffixes [5]string

	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)

	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)

	getSuffix := suffixes[int(math.Floor(base))]

	return strconv.FormatFloat(getSize, 'f', -1, 64) + "" + string(getSuffix)

}

func HumanFileSizeUint(size uint64) string {

	if size < 0 {
		return "0"
	}

	var suffixes [5]string

	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(float64(size)) / math.Log(1024)

	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)

	getSuffix := suffixes[int(math.Floor(base))]

	return strconv.FormatFloat(getSize, 'f', -1, 64) + "" + string(getSuffix)

}

func SortedKeysInt(maps ...map[int]interface{}) (map[int]interface{}, []int) {

	var keys []int
	resultMap := make(map[int]interface{})

	for _, currMap := range maps {
		for k, v := range currMap {
			resultMap[k] = v
			keys = append(keys, k)
		}
	}

	sort.Ints(keys)

	return resultMap, keys

}

func ReverseMapOfStringSlice(ss []*map[string]interface{}) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func DecimalFromString(str string) decimal.Decimal {
	price, err := decimal.NewFromString(str)
	if err != nil {
		panic(err)
	}
	return price
}

// Division with specified precision
func DecimalPow(d1 string, d2 int64) decimal.Decimal {

	d11, err := decimal.NewFromString(d1)

	if err != nil {
		panic(err)
	}

	d22 := decimal.NewFromInt(d2)

	return d11.Pow(d22)

}

func DecimalDiv(d1 string, d2 string, precision int32) decimal.Decimal {

	d11, err := decimal.NewFromString(d1)

	if err != nil {
		panic(err)
	}

	d22, err := decimal.NewFromString(d2)

	if err != nil {
		panic(err)
	}

	return d11.DivRound(d22, precision)

}

func DecimalMul(d1 string, d2 string) decimal.Decimal {

	d11, err := decimal.NewFromString(d1)

	if err != nil {
		panic(err)
	}

	d22, err := decimal.NewFromString(d2)

	if err != nil {
		panic(err)
	}

	return d11.Mul(d22)

}
