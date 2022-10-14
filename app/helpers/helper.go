package helpers

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"goskeleton/app/global/variable"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Claims struct {
	Username string `json:"username"`
	MobileNo string `json:"mobile_no"`
	jwt.StandardClaims
}

const DINE_IN = 0
const DELIVERY = 1
const TAKE_AWAY = 2
const OTHERS = 3

func ParseToken(str string) (*Claims, string) {
	var claimsObj Claims
	tokenString := str
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return variable.PublicKey, nil
	})
	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return nil, "token_expired"
		}
		return nil, "invalid_token"
	}
	text, _ := json.Marshal(claims)
	err = json.Unmarshal(text, &claimsObj)
	if err != nil {
		return nil, "decode_token_failed"
	}
	return &claimsObj, ""
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func ShardHash(str string) string {
	hash := sha256.Sum256([]byte(str))
	hashStr := hex.EncodeToString(hash[:])
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	s := reg.ReplaceAllString(hashStr, "")
	number, err := strconv.Atoi(s[len(s)-2:])
	if err != nil {
		return ""
	}
	number += 1
	f := strconv.Itoa(number)
	return f
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func IntContains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ArrayKeyExists(key interface{}, m map[interface{}]interface{}) bool {
	_, ok := m[key]
	return ok
}

func MbStrlen(str string) int {
	return utf8.RuneCountInString(str)
}

func MbStrLenValid(str string, min int, max int) bool {
	length := MbStrlen(str)

	return (length >= min) && (length <= max)
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func MicroTime() float64 {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	micSeconds := float64(now.Nanosecond()) / 1000000000
	return float64(now.Unix()) + micSeconds
}

func Ip2long(ipstr string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return
}

func Stripos(s string, substr string) int {
	return strings.Index(s, substr)
}

func ParseDecimal(n interface{}, decimal float64) (float64, error) {
	var num float64
	switch i := n.(type) {
	case float64:
		num = i
	case float32:
		num = float64(i)
	case int64:
		num = float64(i)
	case string:
		if s, err := strconv.ParseFloat(i, 64); err == nil {
			num = s
		} else {
			return math.NaN(), errors.New("getFloat: unknown value is of incompatible type")
		}
	default:
		return math.NaN(), errors.New("getFloat: unknown value is of incompatible type")
	}
	parentNo := math.Pow(10, decimal)
	return math.Floor(num*parentNo) / parentNo, nil
}

func GetFileContentType(file []byte) (string, error) {
	// the function that actually does the trick
	contentType := http.DetectContentType(file)

	return contentType, nil
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandStringOtpRunes(n int) string {
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func SerialNo2Id(serialNo string) (int, error) {
	leftover := strings.Trim(serialNo, "MY")

	id, err := strconv.Atoi(leftover)
	if err != nil {
		return 0, err
	}
	return id - 10000, nil
}

func GenDocNo(n int, prefix string) string {
	now := time.Now()
	docNo := RandStringRunes(n) + strconv.Itoa(int(now.Unix()+102800))
	return prefix + docNo
}

func GenStoreNo(prefix string) string {
	now := time.Now()
	docNo := strconv.Itoa(int(now.Unix() + 102800))
	return prefix + docNo
}

func Distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func CalculateDelFee(lat1 float64, lng1 float64, lat2 float64, lng2 float64, areaType string) (string, float64) {
	base := 0.8
	rate := 0.08
	add := 1.1
	if areaType == "" || areaType == "city" {
		base = 1.5
		rate = 0.1
		add = 1.8
	}
	distance := Distance(lat1, lng1, lat2, lng2, "K")
	distanceAdd := distance + add
	fee := (distanceAdd * rate) + base

	tDistance := fmt.Sprintf("%.1f", distanceAdd)
	feeRound, _ := ParseDecimal(fee, 2)
	return tDistance + "KM", feeRound
}

func DefaultPageLimit(page int, limit int) (int, int) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	return page, limit
}

func GenerateSpaceUrl(uploadType string, memId string, fileTypeVal string) string {
	var url string
	app := os.Getenv("APP_NAME")
	url = "https://data.carimakan.biz/prd/"
	if app == "CM_STG" {
		url = "https://data.carimakan.biz/stag/"
	}
	switch uploadType {
	case "merchant.profile":
		url = url + "merchant/" + memId + "/profile/"
		break
	case "merchant.store_cover":
		url = url + "store/" + fileTypeVal + "/cover/"
		break
	case "merchant.store_images":
		url = url + "store/" + fileTypeVal + "/images/"
		break
	}
	return url
}
