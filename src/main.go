/*
 * Dynplug
 *
 * Copyright (C) 2022 Orastron Srl unipersonale
 *
 * Copyright is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3 of the License.
 *
 * Copyright is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Copyright.  If not, see <http://www.gnu.org/licenses/>.
 *
 * File authors: Paolo Marrone
 */

package main

import (
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "syscall"
)

const address = ":10001"
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

    http.HandleFunc("/", handler)

    log.Println("Starting dynplug_server")
    log.Println(http.ListenAndServeTLS(address, "../keys/localhost.crt", "../keys/localhost.key", nil))
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

    data, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    outFile, err := os.CreateTemp("", "dynplug")
    if err != nil {
        errmsg := "Error creating tmp file on disk: " + err.Error()
        log.Println(errmsg)
        unsuccess(w, errmsg)
        return
    }

    log.Println("Writing " + outFile.Name());

    outFile.Write(data)
    if err != nil {
        errmsg := "Copy error: " + err.Error()
        log.Println(errmsg)
        unsuccess(w, errmsg)
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

func unsuccess(w http.ResponseWriter, msg string) {
    w.WriteHeader(500)
    w.Write([]byte("Error: " + msg))
}
