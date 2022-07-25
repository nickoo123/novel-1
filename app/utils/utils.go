// Copyright 2017 Vckai Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/scrypt"
)

var (
	ErrHttpNotFound = errors.New("请求未发现")
	ErrHttpError    = errors.New("请求错误")
	RedisPool       *redis.Pool
)

// ParseBool returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func ParseBool(val interface{}) (value bool, err error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%s", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1 {
				return true, nil
			} else if v == 0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

// convert any type to string
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

type argInt []int

func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

// convert any numeric value to int64
func ToInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}

	return
}

// 返回参数的类型
func Type(v interface{}) string {
	t := reflect.TypeOf(v)
	k := t.Kind()
	return k.String()
}

// 判断是否在数组中
func InArray(in interface{}, list interface{}) bool {
	ret := false
	if in == nil {
		in = ""
	}

	// 判断list是否slice
	l := reflect.TypeOf(list).String()
	t := Type(in)
	if "[]"+t != l {
		return false
	}

	switch t {
	case "string":
		tv := reflect.ValueOf(in).String()
		for _, l := range list.([]string) {
			v := reflect.ValueOf(l)
			if tv == v.String() {
				ret = true
				break
			}
		}

	case "int":
		tv := reflect.ValueOf(in).Int()
		for _, l := range list.([]int) {
			v := reflect.ValueOf(l)
			if tv == v.Int() {
				ret = true
				break
			}
		}
	}

	return ret
}

// gbk convert utf-8
func GBK2UTF(text string) string {
	enc := mahonia.NewDecoder("GB18030")

	text = enc.ConvertString(text)

	return strings.Replace(text, "聽", "&nbsp;", -1)
}

// 通过scrypt生成密码
func NewPass(passwd, salt string) (string, error) {
	dk, err := scrypt.Key([]byte(passwd), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(dk), nil
}

/**
 * unicode transfer utf8
 */
func unicode2utf8(source string) string {
	var res = []string{""}
	sUnicode := strings.Split(source, "\\u")
	var context = ""
	for _, v := range sUnicode {
		var additional = ""
		if len(v) < 1 {
			continue
		}
		if len(v) > 4 {
			rs := []rune(v)
			v = string(rs[:4])
			additional = string(rs[4:])
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			context += v
		}
		context += fmt.Sprintf("%c", temp)
		context += additional
	}
	res = append(res, context)
	return strings.Join(res, "")
}

func InitCache() {
	var conn = beego.AppConfig.String("redisHost") + ":" + beego.AppConfig.String("redisPort")
	var password = beego.AppConfig.String("redisPass")
	var dbNum, _ = strconv.Atoi(beego.AppConfig.String("redisDB"))
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", conn)
		if err != nil {
			return nil, err
		}

		if password != "" {
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
	InitPermissionRedisData()
}

//初始化缓存数据
func InitPermissionRedisData() {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	clearPermissionData(redisConn, redisPrefix)
}

//清除相关redis数据
func clearPermissionData(redisConn redis.Conn, redisPrefix string) {
	redisKeys, err := redis.Strings(redisConn.Do("KEYS", redisPrefix+"*"))
	if err != nil {
		return
	}
	redisConn.Do("MULTI")
	for _, v := range redisKeys {
		redisConn.Do("DEL", v)
	}
	redisConn.Do("EXEC")
}

//设置key value
func SetRedisKeyValue(key, value string) {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()
	redisConn.Do("SET", redisPrefix+key, value)
	redis.Int(redisConn.Do("EXPIRE", redisPrefix+key, 3*60*60))
}

//删除key值
func DelRedisKeys(keys []string) {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()

	redisConn.Do("MULTI")
	for _, v := range keys {
		redisConn.Do("DEL", redisPrefix+v)
	}
	redisConn.Do("EXEC")
}

// 用key值获取数据
func GetRedisKeys(key string) interface{} {
	var redisPrefix = beego.AppConfig.String("redisPrefix")
	redisConn := RedisPool.Get()
	defer redisConn.Close()
	res, err := redis.String(redisConn.Do("Get", redisPrefix+key))
	if err != nil {
		fmt.Println("err:", err)
		return nil
	}
	return res
}
