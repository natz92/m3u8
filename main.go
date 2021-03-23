package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/natz92/m3u8/dl"
	"gopkg.in/yaml.v2"
)

var (
	url            string
	yamlFile       string
	output         string
	outputFileName string
	chanSize       int
)

type FileInfo struct {
	Url        string `yaml:"url"`
	Id         string `yaml:"id"`
	Name       string `yaml:"name"`
	Downloaded bool   `yaml:downloaded`
}

type FileInfos struct {
	Infos []FileInfo `yaml:"files"`
}

func init() {
	flag.StringVar(&url, "u", "", "M3U8 URL, required")
	flag.StringVar(&yamlFile, "i", "", "M3U8 yaml info, required")
	flag.IntVar(&chanSize, "c", 25, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
	flag.StringVar(&outputFileName, "f", "", "Output File Name")
}

func main() {
	f := FileInfos{}

	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	if yamlFile == "" {
		panicParameter("i")
	}
	if outputFileName == "" {
		outputFileName = "merged"
	}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}

	yamlFile, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("--- s\n:%s", yamlFile)
	err = yaml.Unmarshal(yamlFile, &f)
	if err != nil {
		fmt.Printf("Unmarshal: %s", err)
	}
	// fmt.Printf("--- t:\n%v\n\n", f)

	for i, v := range f.Infos {
		if v.Downloaded == true {
			continue
		}
		fmt.Printf("%d, %s: Start.", i, v.Id)

		downloader, err := dl.NewTask("output", v.Name, v.Url)
		if err != nil {
			panic(err)
		}
		if err := downloader.Start(chanSize); err != nil {
			panic(err)
		}

		fmt.Printf("%d, %s: Done.", i, v.Id)
	}

	fmt.Println("\n\n All Done!")
}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}
