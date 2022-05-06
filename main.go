package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, dirname string, isFileNeed bool) error {
	if isFileNeed {
		var intend []bool
		return dirWithFiles(out, dirname, intend, 0)
	} else {
		var intend []bool
		return dirWithoutFiles(out, dirname, intend, 0)
	}
}

func dirWithoutFiles(out io.Writer, dirname string, intend []bool, place int) error {
	dirsWithF, err := os.ReadDir(dirname)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if place >= len(intend) {
		intend = append(intend, true)
	} else {
		intend[place] = true
	}
	var dirs []os.DirEntry

	for _, sortDir := range dirsWithF {
		if sortDir.IsDir() {
			dirs = append(dirs, sortDir)
		}
	}

	for i, dir := range dirs {
		pushSpaces(out, intend, place)
		if i == len(dirs)-1 {
			intend[place] = false
			fmt.Fprint(out, "└───")
		} else {
			fmt.Fprint(out, "├───")
		}
		fmt.Fprintln(out, dir.Name())
		err = dirWithoutFiles(out, dirname+`/`+dir.Name(), intend, place+1)
		if err != nil {
			fmt.Println("error in dirWithoutFiles ", err)
			os.Exit(1)
		}
	}
	return nil
}

func dirWithFiles(out io.Writer, dirname string, intend []bool, place int) error {
	dirs, err := os.ReadDir(dirname)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if place >= len(intend) {
		intend = append(intend, true)
	} else {
		intend[place] = true
	}

	for i, dir := range dirs {
		pushSpaces(out, intend, place)
		if i == len(dirs)-1 {
			intend[place] = false
			fmt.Fprint(out, "└───")
		} else {
			fmt.Fprint(out, "├───")
		}
		if dir.IsDir() {
			fmt.Fprintln(out, dir.Name())
			err = dirWithFiles(out, dirname+`/`+dir.Name(), intend, place+1)
			if err != nil {
				fmt.Println("error in dirWithoutFiles ", err)
				os.Exit(1)
			}
		} else {
			file, err := os.Open(dirname + `/` + dir.Name())
			if err != nil {
				fmt.Println("here1 ", err)
				os.Exit(1)
			}
			stat, err := file.Stat()
			if err != nil {
				fmt.Println("here2 ", err)
				os.Exit(1)
			}
			//name := path.Base(file.Name())
			fmt.Fprint(out, dir.Name(), " ")
			if stat.Size() == 0 {
				fmt.Fprintln(out, "(empty)")
			} else {
				fmt.Fprintf(out, "(%db)\n", stat.Size())
			}
		}
	}
	return nil
}

func pushSpaces(out io.Writer, intend []bool, place int) {
	for i := 0; i < place; i++ {
		if intend[i] {
			fmt.Fprint(out, "│\t")
		} else {
			fmt.Fprint(out, "\t")
		}
	}
}
