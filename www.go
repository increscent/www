package main

import (
  "fmt"
  "net/http"
  "log"
  "io"
  "strings"
  "io/ioutil"
  "os"
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
  filePath := "public" + r.URL.Path
  if filePath[len(filePath) - 1] == '/' {
    filePath = filePath[:len(filePath) - 1]
  }

  fmt.Println(filePath)

  fileInfo, err := os.Lstat(filePath)
  if err != nil {
    sendStatus(w, http.StatusNotFound)
    return
  }

  if fileInfo.IsDir() {
    if _, err = os.Lstat(filePath + "/index.html"); err == nil {
      sendFile(w, filePath + "/index.html")
      return
    } else {
      //show directory
    }
  } else {
    sendFile(w, filePath)
    return
  }

  sendStatus(w, http.StatusInternalServerError)
}

func sendFile(w http.ResponseWriter, path string) {
  contents, err := ioutil.ReadFile(path)
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  w.Write(contents)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
  filePath := "public" + r.URL.Path

  if fileInfo, err := os.Lstat(filePath); err == nil && fileInfo.IsDir() {
    sendStatus(w, http.StatusBadRequest)
    return
  }


  buf := make([]byte, r.ContentLength)
  n, err := io.ReadFull(r.Body, buf)
  if err != nil || n <= 0 {
    sendStatus(w, http.StatusBadRequest)
    return
  }

  sendStatus(w, http.StatusOK)
}

func sendStatus(w http.ResponseWriter, status int) {
  w.WriteHeader(status)
  io.WriteString(w, http.StatusText(status))
}
