package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	host = flag.String("host", "0.0.0.0", "Define what TCP host to bind to")
	port = flag.Int("port", 9000, "Define what TCP port to bind to")
	root = flag.String("root", ".", "Define the root filesystem path")
)

// Get root directory
func getRootDir() string {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	root := string(*root)

	if len(root) == 0 || root == "." {
		root = ""
	}

	if len(root) > 0 && string(root[0]) == "/" {
		return root
	}

	return strings.TrimSuffix(dir, "/") + "/" + root
}

// Get string with terminal color.
func getColorString(firstColor int, lastColor int, message string) string {
	return fmt.Sprintf("\u001b[%d%dm%s\u001b[0m", firstColor, lastColor, message)
}

// Print serve line
func printServeLine() {
	fmt.Printf("%s %s %s %s", getColorString(9, 0, "serving"), getColorString(3, 6, getRootDir()), getColorString(9, 0, "on"), fmt.Sprintf("%s:%d", *host, *port))
	println()
}

func logging(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path

		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}

		// Print served file
		fmt.Printf("%s %s", getColorString(9, 0, r.Method), getColorString(3, 6, path.Base(upath)))
		println()

		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	printServeLine()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), logging(http.FileServer(http.Dir(*root))))

	if err != nil {
		log.Fatal(err)
	}
}
