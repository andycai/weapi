package core

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/andycai/weapi/library/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

var (
	zoneUTC   = time.UTC
	zone      = time.FixedZone("CST", 3600)
	validator = utils.NewValidator()
	lang      = "en"
)

type Ctx = fiber.Ctx

// IP get remote IP
func IP(c *Ctx) string {
	return c.IP()
}

func HashPassword(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h)
}

func CheckPassword(password, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(plain))
	return err == nil
}

func Render(c *Ctx, name string, bind interface{}, layouts ...string) error {
	return c.Render(fmt.Sprintf("%s", name), bind, layouts...)
}

func Validate(i interface{}) error {
	return validator.Validate(i)
}

//#region I18n

func SetLang(l string) {
	lang = l
}

func Lang() string {
	return lang
}

//#endregion

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
var numberRunes = []rune("0123456789")

func randRunes(n int, source []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

func RandText(n int) string {
	return randRunes(n, letterRunes)
}

func RandNumberText(n int) string {
	return randRunes(n, numberRunes)
}

func SafeCall(f func() error, failHandle func(error)) error {
	defer func() {
		if err := recover(); err != nil {
			if failHandle != nil {
				eo, ok := err.(error)
				if !ok {
					es, ok := err.(string)
					if ok {
						eo = errors.New(es)
					} else {
						eo = errors.New("unknown error type")
					}
				}
				failHandle(eo)
			} else {
				// Error(err)
			}
		}
	}()
	return f()
}

func StructAsMap(form any, fields []string) (vals map[string]any) {
	vals = make(map[string]any)
	v := reflect.ValueOf(form)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return vals
	}
	for i := 0; i < len(fields); i++ {
		k := v.FieldByName(fields[i])
		if !k.IsValid() || k.IsZero() {
			continue
		}
		if k.Kind() == reflect.Ptr {
			if !k.IsNil() {
				vals[fields[i]] = k.Elem().Interface()
			}
		} else {
			vals[fields[i]] = k.Interface()
		}
	}
	return vals
}
