package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	path      = "./"
	iconPath  = "icon"
	writeFile = "README.md"
	beginFlag = "### Icons Table"
)

func main() {
	files, err := getIcons(iconPath)
	if err != nil {
		fmt.Printf("err %v", err)
		os.Exit(0)
	}
	if _, err := addIconTable(path+writeFile, files); err != nil {
		fmt.Printf("err %v", err)
	}
}

func getIcons(iconPath string) ([]string, error) {
	var (
		files   []string
		dirPath []string
		dir     []os.FileInfo
		err     error
	)
	if dir, err = ioutil.ReadDir(path+iconPath); err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirPath = append(dirPath, iconPath+PthSep+fi.Name())
			tmpFiles, err := getIcons(iconPath + PthSep + fi.Name())
			if err == nil {
				files = append(files, tmpFiles...)
			}
		} else {
			files = append(files, iconPath+PthSep+fi.Name())
		}
	}
	return files, nil
}

func addIconTable(writeFile string, files []string) (bool, error) {
	f, err := os.OpenFile(writeFile, os.O_RDWR, 0755)
	if err != nil {
		return false, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	var size int64 = 0
	for {
		content, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		size += int64(len(content)) + 2
		if string(content) == beginFlag {
			break
		}
	}
	length := size
	length, err = writeAt(f,"\r\n", length)
	if err != nil {
		return false, err
	}
	length, err = writeAt(f, "| uri | picture | uri | picture | \r\n| - | - | - | - |\r\n", length)
	if err != nil {
		return false, err
	}
	for i := 0 ; i < len(files) - 1; {
		length, err = writeAt(f, "| "+files[i]+" | <img width=\"30\" height=\"30\" src=\""+path+files[i]+"\" />", length)
		// length, err = writeAt(f, "| "+files[i]+" | <img src=\""+files[i]+"\" style=\"height:30px\" />", length)
		if err != nil {
			return false, err
		}
		length, err = writeAt(f, " | "+files[i+1]+" | <img width=\"30\" height=\"30\" src=\""+path+files[i+1]+"\" /> |\r\n", length)
		// length, err = writeAt(f, " | "+files[i+1]+" | <img src=\""+files[i+1]+"\" style=\"height:30px\" /> |\r\n", length)
		if err != nil {
			return false, err
		}
		i += 2
	}
	return true, nil
}

func writeAt(f *os.File, content string, length int64) (len int64, err error){
	l, err := f.WriteAt([]byte(content), length)
	if err != nil {
		return 0, err
	}
	return length + int64(l), nil
}
