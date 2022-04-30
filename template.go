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

type (
	input struct {
		// Enter input fields
	}

	output struct {
		// Enter output fields
	}
)

func solve(in input) output {
	// Implement your solution
	return output{}
}

func sizeSlice(rt reflect.Type, lps []int, p int) (size []int) {
	fsc, lsc := 0, 0
	for i, lp := range lps {
		if p < lp {
			fsc = lp - p
			lsc = len(lps) - i - 1
			break
		}
	}
	if fsc <= 0 || lsc <= 0 {
		log.Panicf("sizeSlice: illegal args lps [%v], p [%v]", lps, p)
	}
	et := rt.Elem()
	if et.Kind() == reflect.Slice {
		size = append([]int{lsc}, sizeSlice(et, lps, p)...)
	} else {
		size = []int{fsc}
	}
	return
}

func size(rv reflect.Value, ft reflect.StructField, lps []int, p int) (size []int) {
	sizeTag := ft.Tag.Get("size")
	if sizeTag == "" {
		if ft.Type.Kind() == reflect.Slice {
			size = sizeSlice(ft.Type, lps, p)
		} else {
			size = []int{1}
		}
	} else {
		for _, sizeStr := range strings.Split(sizeTag, ",") {
			n, err := strconv.Atoi(sizeStr)
			if err != nil {
				sf := rv.FieldByName(sizeStr)
				n = int(sf.Int())
			}
			if n <= 0 {
				log.Panicf("size: got an incorrect number of fields. n [%v]", n)
			}
			size = append(size, n)
		}
	}
	return
}

func in(fs []string, rv reflect.Value, size []int) ([]string, error) {
	switch rv.Kind() {
	case reflect.Int:
		i, err := strconv.ParseInt(fs[0], 0, 0)
		if err != nil {
			return []string{}, err
		}
		rv.SetInt(i)
		return fs[1:], nil
	case reflect.Uint:
		u, err := strconv.ParseUint(fs[0], 0, 0)
		if err != nil {
			return []string{}, err
		}
		rv.SetUint(u)
		return fs[1:], nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(fs[0], 32)
		if err != nil {
			return []string{}, err
		}
		rv.SetFloat(f)
		return fs[1:], nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(fs[0], 64)
		if err != nil {
			return []string{}, err
		}
		rv.SetFloat(f)
		return fs[1:], nil
	case reflect.String:
		rv.SetString(fs[0])
		return fs[1:], nil
	case reflect.Slice:
		var err error
		n := size[0]
		sl := reflect.MakeSlice(rv.Type(), n, n)
		for i := 0; i < n; i++ {
			elem := (reflect.New(rv.Type().Elem())).Elem()
			fs, err = in(fs, elem, size[1:])
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

func deserialize(s []string) (i input) {
	rt := reflect.TypeOf(input{})
	rv := reflect.Indirect(reflect.ValueOf(&i))
	n := rt.NumField()
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
	for l := 0; l < n; l++ {
		ft := rt.Field(l)
		fv := rv.Field(l)
		size := size(rv, ft, lps, p)
		for _, n := range size {
			p += n
		}
		fs, err = in(fs, fv, size)
		if err != nil {
			log.Panicf("deserialize: %v", err)
		}
	}
	return
}

func out(rv reflect.Value) (string, error) {
	switch rv.Kind() {
	case reflect.Int:
		return fmt.Sprint(rv.Int()), nil
	case reflect.Uint:
		return fmt.Sprint(rv.Uint()), nil
	case reflect.Float32:
		return fmt.Sprint(float32(rv.Float())), nil
	case reflect.Float64:
		return fmt.Sprint(rv.Float()), nil
	case reflect.String:
		return rv.String(), nil
	case reflect.Slice:
		var sb strings.Builder
		l := rv.Len()
		for i := 0; i < l; i++ {
			o, err := out(rv.Index(i))
			if err != nil {
				return "", err
			}
			sb.WriteString(o)
			if i < l-1 {
				if rv.Type().Elem().Kind() == reflect.Slice {
					sb.WriteByte('\n')
				} else {
					sb.WriteByte(' ')
				}
			}
		}
		return sb.String(), nil
	default:
		return "", fmt.Errorf("out: not supported kind [%v], value [%v]", rv.Kind(), rv)
	}
}

func (o output) serialize() string {
	rt := reflect.TypeOf(o)
	rv := reflect.ValueOf(o)
	n := rt.NumField()
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ft := rt.Field(i)
		fv := rv.Field(i)
		o, err := out(fv)
		if err != nil {
			log.Panicf("serialize: %v", err)
		}
		sb.WriteString(o)
		if i == n-1 {
			sb.WriteByte('\n')
		} else {
			if ft.Tag.Get("EOL") == "false" {
				sb.WriteByte(' ')
			} else {
				sb.WriteByte('\n')
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
