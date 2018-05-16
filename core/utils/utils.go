/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GenDate() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05 PM")
}

func GenTimestamp() string {
	t := time.Now()
	nanos := t.UnixNano()
	millis := nanos / 1000000 //ms len=13
	return strconv.FormatInt(millis, 10)
}

func GenSpecialTimestamp(fullTimeStr string) (string, error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println(err)
	}

	the_time, err := time.ParseInLocation("2006-01-02 15:04:05", fullTimeStr, local)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	unix_time := the_time.UnixNano() / 1000000
	return strconv.FormatInt(unix_time, 10), nil
}

func GenSpecialTimestampAfterSeconds(timestamp string, seconds int64) (string, error) {
	timestampInt64, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return timestamp, err
	}
	unix_time := timestampInt64 + seconds*1000
	return strconv.FormatInt(unix_time, 10), nil
}

func GenSpecialTimestampAfterSecondsWithFullTimeStr(fullTimeStr string, seconds int64) (string, error) {
	local, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println(err)
	}

	the_time, err := time.ParseInLocation("2006-01-02 15:04:05", fullTimeStr, local)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	unix_time := the_time.UnixNano()/1000000 + seconds*1000
	return strconv.FormatInt(unix_time, 10), nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	var mapObj map[string]interface{}
	objBytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return mapObj, err
	}
	json.Unmarshal(objBytes, &mapObj)
	return mapObj, err
}

func MapToStruct(mapObj map[string]interface{}) (interface{}, error) {
	var obj interface{}
	mapObjBytes, err := json.Marshal(mapObj)
	if err != nil {
		fmt.Println(err)
		return obj, err
	}
	json.Unmarshal(mapObjBytes, &obj)
	return obj, err
}

/*
The json package always orders keys when marshalling. Specifically:

Maps have their keys sorted lexicographically.
Structs keys are marshalled in the order defined in the struct

*/
/*------------------------------ struct serialize must use this -----------------------------*/
/*------------------------------ Hash and Sign use this -----------------------------*/
func StructSerialize(obj interface{}, escapeHTML ...bool) string {
	objMap, err := StructToMap(obj)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if len(escapeHTML) >= 1 {
		return Serialize(objMap, escapeHTML[0])
	}
	return Serialize(objMap)
}

//only for selfTest, format json output
func StructSerializePretty(obj interface{}, escapeHTML ...bool) string {
	objMap, err := StructToMap(obj)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if len(escapeHTML) >= 1 {
		return SerializePretty(objMap, escapeHTML[0])
	}
	return SerializePretty(objMap)
}

/*------------- Structs keys are marshalled in the order defined in the struct ------------------*/
func Serialize(obj interface{}, escapeHTML ...bool) string {
	setEscapeHTML := false
	if len(escapeHTML) >= 1 {
		setEscapeHTML = escapeHTML[0]
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	// disabled the HTMLEscape for &, <, and > to \u0026, \u003c, and \u003e in json string
	enc.SetEscapeHTML(setEscapeHTML)
	err := enc.Encode(obj)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return strings.TrimSpace(buf.String())
	//return strings.Replace(strings.TrimSpace(buf.String()), "\n", "", -1)
}

//only for selfTest, format json output
func SerializePretty(obj interface{}, escapeHTML ...bool) string {
	setEscapeHTML := false
	if len(escapeHTML) >= 1 {
		setEscapeHTML = escapeHTML[0]
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	// disabled the HTMLEscape for &, <, and > to \u0026, \u003c, and \u003e in json string
	enc.SetEscapeHTML(setEscapeHTML)
	enc.SetIndent("", "\t")
	err := enc.Encode(obj)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return strings.TrimSpace(buf.String())
	//return strings.Replace(strings.TrimSpace(buf.String()), "\n", "", -1)
}

func Deserialize(jsonStr string) interface{} {
	var dat interface{}
	err := json.Unmarshal([]byte(jsonStr), &dat)
	if err != nil {
		fmt.Println(err)
	}
	return dat
}

func TypeToString(name interface{}) string {

	value, ok := name.(string)
	if !ok {
		fmt.Println("Type conversion error")
	}
	return value
}

/*
 * try...catch
 * \param [in]
 * \return
 */
func Try(execFunc func(), afterPanic func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			afterPanic(err)
		}
	}()
	execFunc()
}

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
