package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Element struct {
	Log string `json:"log"`
}

func convert(src, dst, format string) error {
	var destinationFile string

	if format == "json" {
		err := printFile(src, dst)
		if err != nil {
			fmt.Printf("File convert failed: %q\n", err)
			return err
		}
		return nil
	} else if format == ".json" {
		file := src[strings.LastIndex(src, "/")+1:]
		nameFile := strings.Split(file, ".")

		destinationFile = dst + "/" + nameFile[0] + format
		err := printFile(src, destinationFile)
		if err != nil {
			fmt.Printf("File convert failed: %q\n", err)
			return err
		}
		return nil
	} else if format != "" {
		file := src[strings.LastIndex(src, "/")+1:]
		nameFile := strings.Split(file, ".")

		destinationFile = dst + "/" + nameFile[0] + format

		fmt.Printf("Converting %s to %s\n", src, destinationFile)
		err := Copy(src, destinationFile)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
			return err
		}
		return nil
	}
	destinationFile = dst

	fmt.Printf("Converting %s to %s\n", src, destinationFile)
	err := Copy(src, destinationFile)
	if err != nil {
		fmt.Printf("File copying failed: %q\n", err)
		return err
	}

	return nil
}

func printFile(src, dst string) error {
	var _, err = os.Stat(dst)

	if os.IsNotExist(err) {
		file, err := os.Create(dst)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer file.Close()

		var dataSlice = make([]Element, 0)
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			dataSlice = append(dataSlice, Element{scanner.Text()})
			// io.WriteString(os.Stdout, scanner.Text())
			// io.WriteString(os.Stdout, "\n")
		}

		bts, err := json.Marshal(dataSlice)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(file, string(bts))
		// fmt.Fprintf(file, "[%s]: ", filename)
		// fmt.Printf("%s", bts)
	} else {
		fmt.Println("File already exists!", dst)
		return err
	}

	fmt.Println("File created successfully", dst)
	return nil
}

func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func help(args, pwd string) {
	if args == "-h" {
		file, err := os.Open("help.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print(string(b))
		os.Exit(2)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please input your flag")
		return
	}

	// setFlag
	file := flag.NewFlagSet(os.Args[1], flag.ExitOnError)

	convertFile := file.String("t", "text", "type convert file")
	createFile := file.String("o", "", "save file to other place")
	// HelpCommand = file.Bool("h", true, "show step how to use tools")

	//---parse the command line into the defined flags---
	source := os.Args[1]

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	help(os.Args[1], pwd)

	if len(os.Args) < 3 {
		err = convert(source, pwd, ".txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return

	} else {
		file.Parse(os.Args[2:])
	}

	switch *convertFile {
	case "json":
		if *createFile != "" {
			err = convert(source, *createFile, "json")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		err = convert(source, pwd, ".json")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return

	case "text":
		if *createFile != "" {
			err = convert(source, *createFile, "")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		err = convert(source, pwd, ".txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return

	default:
		fmt.Printf("Converting %s to %s\n", source, *createFile)
		err := Copy(source, *convertFile)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
			os.Exit(1)
		}
		return
	}
}
