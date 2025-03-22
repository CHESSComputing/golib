package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// code is based on https://github.com/AlanBar13/pass-generator
const voc string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers string = "0123456789"
const symbols string = "!@#$%&*+_-="

// RandomString generates random string
func RandomString() string {
	var str string
	chars := voc + numbers
	for i := 0; i < 16; i++ {
		str += string([]rune(chars)[rand.Intn(len(chars))])
	}
	return str
}

// RecordSize returns actual record size of given interface object
func RecordSize(v interface{}) (int64, error) {
	data, err := json.Marshal(v)
	if err == nil {
		return int64(binary.Size(data)), nil
	}
	return 0, err
}

// Stack returns full runtime stack
func Stack() string {
	trace := make([]byte, 2048)
	count := runtime.Stack(trace, false)
	return fmt.Sprintf("\nStack of %d bytes: %s\n", count, trace)
}

// ErrPropagate helper function which can be used in defer ErrPropagate()
func ErrPropagate(api string) {
	if err := recover(); err != nil {
		log.Println("ERROR", api, "error", err, Stack())
		panic(fmt.Sprintf("%s:%s", api, err))
	}
}

// ErrPropagate2Channel helper function which can be used in goroutines as
// ch := make(chan interface{})
//
//	go func() {
//	   defer ErrPropagate2Channel(api, ch)
//	   someFunction()
//	}()
func ErrPropagate2Channel(api string, ch chan interface{}) {
	if err := recover(); err != nil {
		log.Println("ERROR", api, "error", err, Stack())
		ch <- fmt.Sprintf("%s:%s", api, err)
	}
}

// GoDeferFunc runs any given function in defered go routine
func GoDeferFunc(api string, f func()) {
	ch := make(chan interface{})
	go func() {
		defer ErrPropagate2Channel(api, ch)
		f()
		ch <- "ok" // send to channel that we can read it later in case of success of f()
	}()
	err := <-ch
	if err != nil && err != "ok" {
		panic(err)
	}
}

// BasePath function provides end-point path for given api string
func BasePath(base, api string) string {
	if base != "" {
		if strings.HasPrefix(api, "/") {
			api = strings.Replace(api, "/", "", 1)
		}
		if strings.HasPrefix(base, "/") {
			return fmt.Sprintf("%s/%s", base, api)
		}
		return fmt.Sprintf("/%s/%s", base, api)
	}
	return api
}

// GetHash generates SHA256 hash for given data blob
func GetHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// helper function to extract file name
func FileName(fname string) string {
	arr := strings.Split(fname, "/")
	f := arr[len(arr)-1]
	arr = strings.Split(f, ".")
	return arr[0]
}

// FindFiles find files in given path
func FindFiles(root string) []string {
	var files []string
	if root == "" {
		return files
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("WARNING: unable to access %s/%s, error %v", root, path, err)
		}
		//         log.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Printf("FindFiles root %v, error %v\n", root, err)
	}
	return files
}

// TimeFormat helper function to convert Unix time into human readable form
func TimeFormat(ts interface{}) string {
	var err error
	var t int64
	switch v := ts.(type) {
	case int:
		t = int64(v)
	case int32:
		t = int64(v)
	case int64:
		t = v
	case float64:
		t = int64(v)
	case string:
		t, err = strconv.ParseInt(v, 0, 64)
		if err != nil {
			return fmt.Sprintf("%v", ts)
		}
	default:
		return fmt.Sprintf("%v", ts)
	}
	layout := "2006-01-02 15:04:05"
	return time.Unix(t, 0).UTC().Format(layout)
}

// SizeFormat helper function to convert size into human readable form
func SizeFormat(val interface{}) string {
	var size float64
	var err error
	switch v := val.(type) {
	case int:
		size = float64(v)
	case int32:
		size = float64(v)
	case int64:
		size = float64(v)
	case float64:
		size = v
	case string:
		size, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
	default:
		return fmt.Sprintf("%v", val)
	}
	base := 1000. // CMS convert is to use power of 10
	xlist := []string{"", "KB", "MB", "GB", "TB", "PB"}
	for _, vvv := range xlist {
		if size < base {
			return fmt.Sprintf("%v (%3.1f%s)", val, size, vvv)
		}
		size = size / base
	}
	return fmt.Sprintf("%v (%3.1f%s)", val, size, xlist[len(xlist)])
}

// GetEnv fetches value from user environement
func GetEnv(key string) string {
	for _, item := range os.Environ() {
		value := strings.Split(item, "=")
		if value[0] == key {
			return value[1]
		}
	}
	return ""
}

// FullPath returns full path of given file name wrt to current location
func FullPath(fname string) string {
	if !strings.HasPrefix(fname, "/") {
		// we got relative path (e.g. server_test.json)
		if wdir, err := os.Getwd(); err == nil {
			fname = filepath.Join(wdir, fname)
		}
	}
	return fname
}

// Domain return domain string
func Domain() string {
	domain := "localhost"
	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: unable to get hostname, error:", err)
	}
	if !strings.Contains(hostname, ".") {
		hostname = "localhost"
	} else {
		arr := strings.Split(hostname, ".")
		domain = strings.Join(arr[len(arr)-2:], ".")
	}
	log.Println("Domain", domain)
	return domain
}

// PaddedKey returns padded key up to maxLen
func PaddedKey(key string, maxLen int) string {
	if len(key) < maxLen {
		pad := maxLen - len(key)
		for i := 0; i < pad; i++ {
			key += " "
		}
	}
	return key
}

// Publish2DOIService function publishes record into FOXDEN DOI service
func Publish2DOIService(record map[string]any) (string, string, error) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	var schema, url string
	var err error
	if val, ok := record["schema"]; ok {
		schema = val.(string)
	} else {
		err = errors.New("unable to look-up schema in FOXDEN record")
	}
	if srvConfig.Config.Services.DOIServiceURL == "" {
		return schema, url, errors.New("FOXDEN configuration does not provide DOIServiceURL")
	}
	url = srvConfig.Config.Services.DOIServiceURL
	return schema, url, err
}
