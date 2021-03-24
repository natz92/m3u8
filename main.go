package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/natz92/m3u8/dl"
	"gopkg.in/yaml.v2"
)

var (
	yamlFile string
	output   string
	chanSize int
)

type FileInfo struct {
	Url        string `yaml:"url"`
	Id         string `yaml:"id"`
	Name       string `yaml:"name"`
	Command    string `yaml:"command"`
	Downloaded bool   `yaml:"downloaded"`
}

type FileInfos struct {
	Infos []FileInfo `yaml:"files"`
}

func init() {
	flag.StringVar(&yamlFile, "i", "", "M3U8 yaml info, required")
	flag.IntVar(&chanSize, "c", 25, "Maximum number of occurrences")
	flag.StringVar(&output, "o", "", "Output folder, required")
}

func main() {
	f := FileInfos{}

	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			log.Println("[error]", r)
			os.Exit(-1)
		}
	}()
	if yamlFile == "" {
		panicParameter("i")
	}
	if chanSize <= 0 {
		panic("parameter 'c' must be greater than 0")
	}

	loadYaml, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(loadYaml, &f)
	if err != nil {
		log.Printf("Unmarshal: %s", err)
	}

	for i, v := range f.Infos {
		f.Infos[i].Command = fmt.Sprintf("ffmpeg.exe -i \"C:\\Users\\tkddn\\Documents\\m3u8\\output\\%s.ts\" -map 0 -c copy \"%s.mp4\"", v.Name, v.Name)
		if v.Downloaded {
			continue
		}
		log.Printf("%d, %s: Start.", i, v.Id)

		downloader, err := dl.NewTask("output", v.Name, v.Url)
		if err != nil {
			log.Fatalln(err)
		}
		if err := downloader.Start(chanSize); err != nil {
			log.Fatalln(err)
		}

		f.Infos[i].Downloaded = true
		log.Printf("%d, %s: Done.", i, v.Id)
	}

	yaml.FutureLineWrap()
	d, err := yaml.Marshal(&f)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile(yamlFile, d, 0644)

	log.Println("\n\n All Done!")
}

func panicParameter(name string) {
	panic("parameter '" + name + "' is required")
}
