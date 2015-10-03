package json

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	body_r *regexp.Regexp
	unit_r *regexp.Regexp
	err    error
)

func init() {
	body_r, err = regexp.Compile(`{[ ,:'"0-9a-zA-Z]+}`)
	checkErr(err)
	unit_r, err = regexp.Compile(`^[0-9a-zA-Z]+$`)
	checkErr(err)
}

type unit struct {
	value []byte
	err   error
}
type Json struct {
	origin []byte
	body   []byte
	keys   []unit
	values []unit
	err    error
}

func NewJson(o []byte) *Json {
	return &Json{origin: o}
}
func (j *Json) OutputMap() (map[string]string, error) {
	j.getBody()
	j.split()
	j.trimSpace()
	j.trim('\'')
	j.trim('"')
	if !j.checkUnit() {
		err := j.err
		j.close()
		return nil, err
	}
	var result map[string]string = make(map[string]string)
	for i, _ := range j.keys {
		result[string(j.keys[i].value)] = string(j.values[i].value)
	}
	return result, nil
}
func (j *Json) close() {
	j = nil
}
func (j *Json) getBody() {
	body_tmp := body_r.Find(j.origin)
	if body_tmp != nil {
		j.body = body_tmp[1 : len(body_tmp)-1]
	}
}
func (j *Json) split() {
	for _, v := range split(j.body, ',') {
		val := split(v, ':')
		if len(val) != 2 {
			j.err = errors.New("wrong format with :")
		}
		j.keys = append(j.keys, unit{value: val[0]})
		j.values = append(j.values, unit{value: val[1]})
	}
}
func split(value []byte, flag byte) [][]byte {
	var result [][]byte
	var index int
	for i, v := range value {
		if v == flag {
			result = append(result, value[index:i])
			index = i + 1
		}
	}
	return append(result, value[index:])
}
func (j *Json) trim(flag byte) {
	for i, _ := range j.keys {
		(&j.keys[i]).tirm(flag)
	}
	for i, _ := range j.values {
		(&j.values[i]).tirm(flag)
	}
}
func (j *Json) trimSpace() {
	for i, _ := range j.keys {
		(&j.keys[i]).trimSpace()
	}
	for i, _ := range j.values {
		(&j.values[i]).trimSpace()
	}
}
func (u *unit) tirm(flag byte) {
	length := len(u.value)
	if u.value[0] == flag && u.value[length-1] == flag {
		u.value = u.value[1 : length-1]
	}
}
func (u *unit) trimSpace() {
	for i, v := range u.value {
		if v != ' ' {
			u.value = u.value[i:]
			break
		}
	}
	length := len(u.value)
	for i, _ := range u.value {
		if u.value[length-i-1] != ' ' {
			u.value = u.value[:length-i]
			break
		}
	}
}

func (j *Json) checkUnit() bool {
	for _, v := range append(j.keys, j.values...) {
		if !v.checkUnit() {
			j.err = v.err
			return false
		}
	}
	return true
}
func (u *unit) checkUnit() bool {
	if unit_r.Find(u.value) == nil {
		fmt.Println(string(u.value))
		u.err = errors.New("invalid symbol")
		return false
	}
	return true
}
func checkErr(err error) bool {
	if err != nil {
		dealtError(err)
		return true
	}
	return false
}
func dealtError(err error) {
	fmt.Println(err)
}
