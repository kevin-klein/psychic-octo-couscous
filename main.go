package main

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func worker(jobs <-chan string, done chan<- bool) {
	for path := range jobs {
		file, err := os.Open(path)
		if err != nil {
			done <- false
			log.Println(err)
			return
		}
		defer file.Close()
		image, err := jpeg.Decode(file)
		if err != nil {
			done <- false
			log.Println(err)
			return
		}

		m := resize.Resize(1280, 0, image, resize.NearestNeighbor)

		out, err := os.Create("./verkleinert/" + path)
		if err != nil {
			done <- false
			log.Println(err)
			return
		}
		defer out.Close()

		jpeg.Encode(out, m, nil)

		done <- true
	}
}

func main() {
	jobs := make(chan string, 100)
	results := make(chan bool, 100)

	for w := 0; w < runtime.NumCPU(); w++ {
		go worker(jobs, results)
	}

	files, _ := filepath.Glob("./*.jpg")
	os.Mkdir("./verkleinert", os.ModePerm)
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	for i := 0; i < len(files); i++ {
		<-results
	}

}
