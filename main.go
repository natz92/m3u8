package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/natz92/m3u8/dl"
)

var (
	url            string
	output         string
	outputFileName string
	chanSize       int
)

func init() {
	flag.StringVar(&url, "u", "", "M3U8 URL, required")
	flag.IntVar(&chanSize, "c", 25, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
	flag.StringVar(&outputFileName, "f", "", "Output File Name")
}

func main() {
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	if url == "" {
		panicParameter("u")
	}
	if output == "" {
		panicParameter("o")
	}
	if outputFileName == "" {
		outputFileName = "merged"
	}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}
	downloader, err := dl.NewTask(output, outputFileName, url)
	if err != nil {
		panic(err)
	}
	if err := downloader.Start(chanSize); err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}
