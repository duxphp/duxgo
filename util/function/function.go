package function

import (
	"bytes"
	"crypto/md5"
	"github.com/duxphp/duxgo/core"
	"encoding/hex"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"
)

// HashEncode 密文加密
func HashEncode(content []byte) string {
	hash, err := bcrypt.GenerateFromPassword(content, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

// HashVerify 验证密文
func HashVerify(hashedPwd string, password []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		return false
	}
	return true
}

// PageLimit 分页计算
func PageLimit(page int, total int, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		return 0, 0
	}
	totalPage := total / limit
	if total%limit != 0 {
		totalPage++
	}
	offset := (page - 1) * limit
	return offset, totalPage
}

// UcFirst 首字母转大写
func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// LcFirst 首字母小写
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// MapPluck 提取键值
func MapPluck(data []map[string]any, value string, key ...string) map[any]any {
	newData := map[any]any{}
	for index, item := range data {
		var name any
		if len(key) > 0 {
			name = item[key[0]]
		} else {
			name = index
		}
		newData[name] = item[value]

	}
	return newData

}

// Url 普通编译Url
func Url(urlString string, params map[string]any, absolutes ...bool) string {
	var uri url.URL
	q := uri.Query()
	for k, v := range params {
		q.Add(k, cast.ToString(v))
	}
	urlBuild := urlString + "?" + q.Encode()
	var absolute bool
	if len(absolutes) > 0 {
		absolute = absolutes[0]
	}
	if absolute {
		urlBuild = core.Config["app"].GetString("app.baseUrl") + urlBuild
	}
	return urlBuild
}

// BuildUrl 编译Url
func BuildUrl(urlString string, params map[string]any, absolute bool, expand ...any) string {
	// 前缀js变量
	var prefix string
	if len(expand) > 0 {
		prefix = cast.ToString(expand[0]) + "."
	}
	// 默认字段
	var fields []string
	if len(expand) > 1 {
		fields = cast.ToStringSlice(expand[1])
	}

	var uri url.URL
	q := uri.Query()

	paramsFix := map[string]string{}
	paramsModel := map[string]string{}

	for key, value := range params {
		val := cast.ToString(value)
		val = strings.TrimSpace(val)
		if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
			val = strings.Trim(val, "{")
			val = strings.Trim(val, "}")
			paramsFix[key] = val
			q.Add(key, "xx"+key+"xx")
		} else {
			_, ok := lo.Find[string](fields, func(i string) bool {
				return i == value
			})
			if ok {
				paramsModel[key] = val
				q.Add(key, "xx"+key+"xx")
			} else {
				q.Add(key, cast.ToString(val))
			}
		}
	}
	urlBuild := urlString + "?" + q.Encode()
	if absolute {
		urlBuild = core.Config["app"].GetString("app.baseUrl") + urlBuild
	}
	for k, v := range paramsFix {
		urlBuild = strings.Replace(urlBuild, "xx"+k+"xx", "${ "+v+" || ''}", -1)
	}
	for k, v := range paramsModel {
		urlBuild = strings.Replace(urlBuild, "xx"+k+"xx", "${ "+prefix+v+" || ''}", -1)
	}
	return "`" + urlBuild + "`"
}

// FormatFileSize 格式化文件大小
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// RandString 随机字符
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// Md5 生成32位MD5
func Md5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// FileMd5 文件MD5
func FileMd5(data []byte) string {
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

// IsExist 判断目录文件存在
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// CreateDir 创建目录
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0777)
	if err != nil {
		core.Logger.Error().Err(err).Msg(dirName)
		return false
	}
	return true
}

// GetUuid 获取uuid
func GetUuid() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// Round 四舍五入保留小数
func Round(val float64, precision int) float64 {
	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}

// InTimeSpan 范围时间查询
func InTimeSpan(start, end, check time.Time, includeStart, includeEnd bool) bool {
	_start := start
	_end := end
	_check := check
	if end.Before(start) {
		_end = end.Add(24 * time.Hour)
		if check.Before(start) {
			_check = check.Add(24 * time.Hour)
		}
	}
	if includeStart {
		_start = _start.Add(-1 * time.Nanosecond)
	}
	if includeEnd {
		_end = _end.Add(1 * time.Nanosecond)
	}
	return _check.After(_start) && _check.Before(_end)
}

// CtxBody 提取body
func CtxBody(ctx echo.Context) []byte {
	s, err := ioutil.ReadAll(ctx.Request().Body)
	ctx.Request().Body.Close()
	ctx.Request().Body = ioutil.NopCloser(bytes.NewReader(s))
	if err != nil {
		return []byte("")
	}
	return s
}
