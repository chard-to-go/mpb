package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vbauerster/mpb"
)

func main() {
	url := "https://homebrew.bintray.com/bottles/libtiff-4.0.7.sierra.bottle.tar.gz"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server return non-200 status: %s\n", resp.Status)
		return
	}

	size := resp.ContentLength

	// create dest
	destName := filepath.Base(url)
	dest, err := os.Create(destName)
	if err != nil {
		fmt.Printf("Can't create %s: %v\n", destName, err)
		return
	}
	defer dest.Close()

	p := mpb.New().SetWidth(64)
	// if you omit following line, download will complete fine, but rendering bar
	// may not complete, thus better always use even in single thread.
	p.Wg.Add(1)

	bar := p.AddBar(int(size)).PrependCounters(mpb.UnitBytes, 19).AppendETA()

	// create proxy reader
	reader := bar.ProxyReader(resp.Body)

	// and copy from reader, ignoring errors
	io.Copy(dest, reader)

	p.WaitAndStop()
	fmt.Println("Finished")
}
