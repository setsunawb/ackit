package ackit

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

const (
	inDir       = "testdata/in"
	outDir      = "testdata/out"
	logFileName = "output.log"
)

type FileLogger struct {
	*os.File
}

func NewFileLogger(name string) (*FileLogger, error) {
	logfile, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &FileLogger{logfile}, nil
}

func (l *FileLogger) Logf(format string, a ...interface{}) {
	l.WriteString(fmt.Sprintf(format, a...))
}

type testdata struct {
	name string
	in   string
	out  string
}

func (d *testdata) string() string {
	return fmt.Sprintf("name: %v\n[in]\n%v\n[out]\n%v\n", d.name, d.in, d.out)
}

func getTestdata() ([]testdata, error) {
	ifs, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Panicln(err)
	}
	n := len(ifs)
	data := make([]testdata, n)

	for i := 0; i < n; i++ {
		name := ifs[i].Name()
		in, err := ioutil.ReadFile(path.Join(inDir, name))
		if err != nil {
			return nil, err
		}

		out, err := ioutil.ReadFile(path.Join(outDir, name))
		if err != nil {
			return nil, err
		}
		data[i] = testdata{name: name, in: string(in), out: string(out)}
	}
	return data, nil
}

func Test(f func(), t *testing.T) {
	t.Helper()
	stdin, stdout := os.Stdin, os.Stdout
	defer func() {
		os.Stdin, os.Stdout = stdin, stdout
	}()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		r.Close()
		w.Close()
	}()

	os.Stdin = r
	os.Stdout = w
	if err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)

	logger, err := NewFileLogger(logFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	data, err := getTestdata()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range data {
		var (
			bytes = []byte(tt.in)
			acm   = 0
			i     = 0
			err   error
		)
		logger.Logf("%v", tt.string())

		ch := make(chan struct{})
		go func() {
			start := time.Now()
			f()
			elapse := time.Since(start)
			logger.Logf("elapse: %vms\n", elapse.Milliseconds())
			ch <- struct{}{}
		}()

		for {
			i, err = writer.WriteString(tt.in + "\n")
			if err != nil {
				log.Panicln(err)
			}
			writer.Flush()

			acm += i
			if acm >= len(bytes) {
				break
			}
		}
		<-ch

		t.Run(tt.name, func(t *testing.T) {

			var sb strings.Builder
			for scanner.Scan() {
				if scanner.Text() == "" {
					break
				}
				sb.WriteString(scanner.Text() + "\n")
			}
			ans := sb.String()
			logger.Logf("[ans]\n%v\n\n\n", ans)
			if ans != string(tt.out) {
				t.Errorf("got [%#v] , want [%#v]", ans, tt.out)
			}
		})
	}
}
