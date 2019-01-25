// base template -

package main

import (
    "io"
    "net/http"
    "os"
)

func main() {

    fileUrl := "https://golangcode.com/images/avatar.jpg"

    err := DownloadFile("avatar.jpg", fileUrl)
    if err != nil {
        panic(err)
    }

}

func DownloadFile(filepath string, url string) error {

    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}