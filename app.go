package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/chrisprobst/songtags"
)

var counter int
var mutex sync.Mutex

func recognizeHandler(w http.ResponseWriter, r *http.Request) {
	// Grab input file
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	// Make room for file
	mutex.Lock()
	fp := fmt.Sprint("/tmp/recognize.", counter)
	out, err := os.Create(fp)
	counter++
	mutex.Unlock()
	if err != nil {
		fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		return
	}
	defer out.Close()

	// Write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	fmt.Fprintf(w, "File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)
	out.Close()

	// Calc songtags
	res := songtags.ForFile(fp)
	fmt.Fprintf(w, "Result: ")
	fmt.Fprintf(w, res)
}

func main() {
	http.HandleFunc("/recognize", recognizeHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":7133", nil)
}
