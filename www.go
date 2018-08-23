package main

import (
  "fmt"
  "errors"
  "net/http"
  "log"
  "io"
  "strings"
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
  relativePath := r.URL.Path
  if relativePath[len(relativePath) - 1] == '/' {
    relativePath = relativePath[:len(relativePath) - 1]
  }
  filePath := "public" + relativePath

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
      sendDir(w, filePath, relativePath)
      return
    }
  } else {
    sendFile(w, filePath)
    return
  }

  sendStatus(w, http.StatusInternalServerError)
}

func sendFile(w http.ResponseWriter, path string) {
  fileInfo, err := os.Lstat(path)
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  buf := make([]byte, fileInfo.Size())
  file, err := os.Open(path)
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  _, err = file.Read(buf)
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  w.Write(buf)
}

func sendDir(w http.ResponseWriter, filePath string, relativePath string) {
  file, err := os.Open(filePath)
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  files, err := file.Readdir(0)
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  var b strings.Builder
  for i := 0; i < len(files); i++ {
    name := files[i].Name()
    b.WriteString("<a href=\"" + relativePath + "/" + name + "\">" + name + "</a>")
  }
  io.WriteString(w, b.String())
}

func putHandler(w http.ResponseWriter, r *http.Request) {
  filePath := "public" + r.URL.Path

  fileInfo, err := os.Lstat(filePath)

  if err == nil && fileInfo.IsDir() {
    sendStatus(w, http.StatusBadRequest)
    return
  }

  created := err != nil

  if created {
    _, err = createFile(filePath)
    if err != nil {
      sendStatus(w, http.StatusNotAcceptable)
      return
    }
  }

  file, err := os.OpenFile(filePath, os.O_WRONLY, os.FileMode(0644))
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  buf := make([]byte, r.ContentLength)
  n, err := io.ReadFull(r.Body, buf)
  if err != nil || n <= 0 {
    sendStatus(w, http.StatusBadRequest)
    return
  }

  err = file.Truncate(int64(n))
  if err != nil {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  n, err = file.Write(buf)
  if err != nil || n <= 0 {
    sendStatus(w, http.StatusInternalServerError)
    return
  }

  file.Close()

  if created {
    sendStatus(w, http.StatusCreated)
  } else {
    sendStatus(w, http.StatusOK)
  }
}

func createFile(path string) ([]string, error) {
  var created []string
  dirs := strings.Split(path, "/")
  for i := 0; i < len(dirs) - 1; i++ {
    fp := strings.Join(dirs[:i+1], "/")
    fi, err := os.Lstat(fp)
    if err != nil {
      err := os.Mkdir(fp, os.FileMode(0755))
      fmt.Println(fp)
      if err != nil {
        return created, err
      }
      created = append(created, fp)
    } else {
      if !fi.IsDir() {
        return created, errors.New("file in path")
      }
    }
  }

  f, err := os.OpenFile(path, os.O_CREATE, 0644)
  if err != nil {
    return created, err
  }
  created = append(created, path)
  f.Close()
  return created, nil
}

func sendStatus(w http.ResponseWriter, status int) {
  w.WriteHeader(status)
  io.WriteString(w, http.StatusText(status))
}
