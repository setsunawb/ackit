package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type input struct {
	
}

type output struct {
	
}

func solve(in input) output {
	
}

func in(fs []string, rv reflect.Value, nums []int) ([]string, error) {
	switch rv.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(fs[0])
		if err != nil {
			return []string{}, err
		}
		rv.SetInt(int64(i))
		return fs[1:], nil
	case reflect.String:
		rv.SetString(fs[0])
		return fs[1:], nil
	case reflect.Slice:
		var err error
		num := nums[0]
		sl := reflect.MakeSlice(rv.Type(), num, num)
		for i := 0; i < num; i++ {
			elem := (reflect.New(rv.Type().Elem())).Elem()
			fs, err = in(fs, elem, nums[1:])
			if err != nil {
				return []string{}, err
			}
			sl.Index(i).Set(elem)
		}
		rv.Set(sl)
		return fs, nil
	default:
		return []string{}, fmt.Errorf("in: not supported kind [%v]", rv.Kind())
	}
}

func nums(rv reflect.Value, ft reflect.StructField, lps []int, p int) (nums []int) {
	numTag := ft.Tag.Get("num")
	if numTag == "" {
		if ft.Type.Kind() == reflect.Slice {
			for _, lp := range lps {
				if p < lp {
					return []int{lp - p}
				}
			}
			log.Panicf("nums: illegal args lps [%v], p [%v]", lps, p)
		} else {
			nums = []int{1}
		}
	} else {
		for _, numStr := range strings.Split(numTag, ",") {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				numField := rv.FieldByName(numStr)
				num = int(numField.Int())
			}
			if num <= 0 {
				log.Panicf("nums: got an incorrect number of fields. num [%v]", num)
			}
			nums = append(nums, num)
		}
	}
	return
}

func deserialize(s []string) (i input) {
	rt := reflect.TypeOf(input{})
	rv := reflect.Indirect(reflect.ValueOf(&i))
	nf := rt.NumField()
	var (
		fs  []string
		lps []int
		p   = 0
		err error
	)
	for _, l := range s {
		lfs := strings.Fields(l)
		fs = append(fs, lfs...)
		lps = append(lps, len(fs))
	}
	for l := 0; l < nf; l++ {
		ft := rt.Field(l)
		fv := rv.Field(l)
		nums := nums(rv, ft, lps, p)
		for _, n := range nums {
			p += n
		}
		fs, err = in(fs, fv, nums)
		if err != nil {
			log.Panicf("deserialize: %v", err)
		}
	}
	return
}

func out(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Int:
		return fmt.Sprint(v.Int()), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Slice:
		var sb strings.Builder
		for i := 0; i < v.Len(); i++ {
			o, err := out(v)
			if err != nil {
				return "", err
			}
			sb.WriteString(o)
			if v.Elem().Kind() == reflect.Slice {
				sb.WriteByte('\n')
			} else {
				sb.WriteByte(' ')
			}
		}
		return sb.String(), nil
	default:
		return "", fmt.Errorf("out: cannot out reflect.value [%v]", v)
	}
}

func (o output) serialize() string {
	rt := reflect.TypeOf(o)
	rv := reflect.ValueOf(o)
	nf := rt.NumField()
	var sb strings.Builder
	for i := 0; i < nf; i++ {
		ft := rt.Field(i)
		fv := rv.Field(i)
		o, err := out(fv)
		if err != nil {
			log.Panicf("serialize: %v", err)
		}
		sb.WriteString(o)
		if i == nf-1 {
			sb.WriteByte('\n')
		} else {
			if ft.Tag.Get("EOL") == "true" {
				sb.WriteByte('\n')
			} else {
				sb.WriteByte(' ')
			}
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func interact(is []string) string {
	i := deserialize(is)
	o := solve(i)
	return o.serialize()
}

func readLine(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

func main() {
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	var lines []string

	for {
		b, err := readLine(r)

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Panicln(err)
		}
		lines = append(lines, string(b))

		if string(b) == "" {
			break
		}
	}

	_, err := w.WriteString(interact(lines))
	if err != nil {
		log.Panicln(err)
	}
	w.Flush()
}
