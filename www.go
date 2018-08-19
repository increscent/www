package main

import (
  "net/http"
  "log"
  "io"
  "strings"
  "io/ioutil"
  "regexp"
)

func main() {
  http.HandleFunc("/", handler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
  switch method := strings.ToUpper(r.Method); method {
  case "GET":
    getHandler(w, r)
  case "PUT":
    putHandler(w, r)
  default:
    sendStatus(w, http.StatusMethodNotAllowed)
  }
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  filePath := r.URL.Path

  if r.URL.Path == "/" {
    filePath = "/index.html"
  }

  fileContents, err := getFileContents(filePath)
  if err != nil {
    sendStatus(w, http.StatusNotFound)
    return
  }

  w.Write(fileContents)
}

func getFileContents(relativePath string) ([]byte, error) {
  contents, err := ioutil.ReadFile("public" + relativePath)
  if err == nil {
    return contents, nil
  } else {
    contents, err = ioutil.ReadFile("public" + relativePath + ".html")
    return contents, err
  }
}

func putHandler(w http.ResponseWriter, r *http.Request) {
  filePath := r.URL.Path

  if filePath[len(filePath) - 1] == '/' {
    sendStatus(w, http.StatusBadRequest)
    return
  }

  if !regexp.MatchString("\\.html", filePath) {
    if _, err := getFileContents(filePath + ".html"); err == nil {
      filePath += ".html"
    }
  }

  buf := make([]byte, r.ContentLength)
  _, err := io.ReadFull(r.Body, buf)
  if err != nil || n <= 0 {
    sendStatus(w, http.StatusBadRequest)
  }

  sendStatus(w, http.StatusOK)
}

func sendStatus(w http.ResponseWriter, status int) {
  w.WriteHeader(status)
  io.WriteString(w, http.StatusText(status))
}
