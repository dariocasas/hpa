package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var startedAt = time.Now()
var randIniSecs float64 = rand.Float64() * 20

func main() {

	http.HandleFunc("/", hello)
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/configmap", configmap)
	http.HandleFunc("/secret", secret)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func healthz(w http.ResponseWriter, _ *http.Request) {

	duration := time.Since(startedAt)

	okSecs, err := strconv.ParseFloat(os.Getenv("OK_SECS"), 64)
	if err != nil {
		okSecs = 60
	}

	iniSecs, err := strconv.ParseFloat(os.Getenv("INI_SECS"), 64)
	if err != nil {
		iniSecs = 10
	}

	printInfo := func() {
		fmt.Fprintf(w, "randIniSecs: %v\n", randIniSecs)
		fmt.Fprintf(w, "iniSecs: %v\n", iniSecs)
		fmt.Fprintf(w, "ok_secs: %v\n", okSecs)
	}

	if duration.Seconds() < iniSecs+randIniSecs {
		w.WriteHeader(http.StatusInternalServerError)
		printInfo()
		fmt.Fprintf(w, "Initializing. Duration: %v", duration.Seconds())
	} else if duration.Seconds() > okSecs+iniSecs+randIniSecs {
		w.WriteHeader(http.StatusInternalServerError)
		printInfo()
		fmt.Fprintf(w, "Error. Duration: %v", duration.Seconds())
	} else {
		printInfo()
		fmt.Fprintf(w, "ok, running since %v", duration.Seconds())
	}
}

func secret(w http.ResponseWriter, _ *http.Request) {

	user := os.Getenv("USER")
	pass := os.Getenv("PASS")

	fmt.Fprintf(w, "user: %s, pass: %s", user, pass)

}

func configmap(w http.ResponseWriter, _ *http.Request) {
	data, err := os.ReadFile("/app/config/members.txt")
	if err != nil {
		fmt.Fprintf(w, "error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Members: %s.\n", data)

	data, err = os.ReadFile("/app/config/ui-properties.conf")
	if err != nil {
		fmt.Fprintf(w, "error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ui-properties: %s.\n", data)
}

func hello(w http.ResponseWriter, _ *http.Request) {

	name := os.Getenv("NAME")
	age := os.Getenv("AGE")

	fmt.Fprintf(w, "Hello! I'm %s, I'm %s", name, age)

}
