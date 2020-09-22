package misc

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

type Utils struct {
}

func ToUint32(param string) uint32 {

	if s, err := strconv.ParseUint(param, 10, 32); err == nil {
		return uint32(s)
	}
	return 0
}

func ToUint8(param string) uint8 {

	if s, err := strconv.ParseUint(param, 10, 8); err == nil {
		return uint8(s)
	}
	return 0
}

func ContainsInArray(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Enum2Array(enumstring string) []string {

	enumstring = strings.ReplaceAll(enumstring, "{", "")
	enumstring = strings.ReplaceAll(enumstring, "}", "")

	return strings.Split(enumstring, ",")
}

/**
 * Метод возвращает NULL если строка пустая для SQL
 */
func StringSql(param string) *string {

	nullable := &param
	if param == "" {
		nullable = nil
	}
	return nullable
}

/**
 * Метод возвращает NULL если int64 пустой
 */
func Int64Sql(param int64) *int64 {

	nullable := &param
	if param == 0 {
		nullable = nil
	}
	return nullable
}

/**
 * Метод превращает IP в Integer
 */
func Ip2uint32(ip string) uint32 {

	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}
