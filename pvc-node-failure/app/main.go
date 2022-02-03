package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type app struct {
	store string
	lock  sync.Mutex
}

func main() {

	filename := flag.String("file", "count.txt", "location of file to store counter")
	flag.Parse()

	a := app{
		store: *filename,
		lock:  sync.Mutex{},
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		count, err := a.count()
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(fmt.Sprintf("%d\n", count)))
	})

	http.ListenAndServe(":8080", r)
}

func (a *app) count() (int, error) {

	a.lock.Lock()
	defer a.lock.Unlock()

	count := 0

	fh, err := os.OpenFile(a.store, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return count, err
	}
	defer fh.Close()

	r := bufio.NewReader(fh)
	line, _ := r.ReadString('\n') // we swallow the error because we overwrite file contents

	line = strings.TrimSuffix(line, "\n")
	count, err = strconv.Atoi(line)
	if err != nil {
		count = 0
	}

	count++
	fh.Truncate(0)
	fh.Seek(0, 0)
	_, err = fh.WriteString(fmt.Sprintf("%d\n", count))
	fh.Close()
	return count, err
}
