package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

var addr = flag.String("addr", ":8080", "address to bind")

type Req struct {
	Program string `json:"program"`
	Example string `json:"example"`
}

type Resp struct {
	Err        error             `json:"err"`
	Match      bool              `json:"match"`
	Submatches map[string]string `json:"submatches"`
}

func main() {
	flag.Parse()

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			fmt.Fprintf(w, "Incorrect HTTP method")
			return
		}
		checkError := mkCheckError(w)
		decoder := json.NewDecoder(r.Body)
		values := Req{}
		err := decoder.Decode(&values)
		if checkError(err) {
			log.Printf("could not decode body: %s", err)
			return
		}

		resp := Resp{}
		exp, err := regexp.Compile(values.Program)
		if err != nil {
			resp.Err = err
			b, _ := json.Marshal(resp)
			w.Write(b)
			return
		}
		resp.Match = exp.MatchString(values.Example)
		if resp.Match {
			resp.Submatches = ParseValues(exp, values.Example)
		}
		b, _ := json.Marshal(resp)
		w.Write(b)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func mkCheckError(w http.ResponseWriter) func(error) bool {
	return func(err error) bool {
		if err != nil {
			log.Printf("meme failed: %s", err)
			w.WriteHeader(500)
			e, _ := json.Marshal(err)
			w.Write(e)
			return true
		}
		return false
	}
}

func ParseValues(r *regexp.Regexp, body string) map[string]string {
	out := map[string]string{}
	subs := r.FindStringSubmatch(body)
	if len(subs) == 0 {
		return out
	}
	for i, n := range r.SubexpNames() {
		out[n] = subs[i]
	}
	return out
}
