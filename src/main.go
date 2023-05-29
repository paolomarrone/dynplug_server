package main

import (
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "syscall"
)

const pipename = "dynplug_magicpipe"
var pipepath string

func main() {
    // Estabilish the pipe
    pipepath = path.Join(os.TempDir(), "/", pipename)
    err := syscall.Mkfifo(pipepath, 0666)
    if err != nil {
        if err.Error() == "file exists" {
            log.Println("named pipe already exists: ok")
        } else {
            log.Fatal("Mkfifo error:", err)
        }
    }

    server := http.Server{
        Handler: http.HandlerFunc(handler),
        Addr: ":10001",
    }

    log.Println("Starting dynplug_server")
    log.Println(server.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/uploadfile":
        handleFileInForm(w, r)
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("Not Found"))
    }
}

func handleFileInForm(w http.ResponseWriter, r *http.Request) {
    
    /*
    f, fh, err := r.FormFile("file")
    if err != nil {
        log.Println("formfile error", err)
        unsuccess(w)
        return
    }

    if fh.Size <= 0 {
        unsuccess(w)
        return
    }
    */

    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    outFile, err := os.CreateTemp("", "dynplug")
    if err != nil {
        log.Println("Error creating tmp file on disk", err)
        unsuccess(w)
        return
    }

    log.Println("Writing " + outFile.Name());

    //written, err := io.Copy(outFile, f)
    outFile.Write(data)
    if err != nil {
        log.Println("copy error", err)
        unsuccess(w)
        return
    }

    p, err := os.OpenFile(pipepath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0777)
    if err != nil {
        log.Fatalf("Error opening pipe: %v", err)
    }
    p.WriteString(outFile.Name())
    p.Close();

    success(w)
    log.Println("Written")
}

func success(w http.ResponseWriter) {
    w.WriteHeader(200)
    w.Write([]byte("Got your file\n"))
}

func unsuccess(w http.ResponseWriter) {
    w.WriteHeader(500)
    w.Write([]byte("Something went wrong\n"))
}
