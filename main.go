package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Missing arguments! please provide a pom.xml and output dir")
		return
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Can't read file:", os.Args[1])
		panic(err)
	}



	// data is the file content, you can use it
	var pom = parsePom(data)

	absPath, _ := filepath.Abs(os.Args[2])
	os.MkdirAll(absPath, os.ModePerm)

	// load dependencies
	for i := range pom.Dependencies.Dependency {
		var dep = pom.Dependencies.Dependency[i]
		if dep.Scope != "provided" && dep.Scope != "test" && dep.Scope != "pom" {
			fmt.Println("Extracting " + dep.ArtifactId + " from local maven cache")
			var p = dep.localJarPath()
			if p != "" {
				copy(p, absPath + "/" + dep.ArtifactId + "-" + dep.Version + ".jar")
			}
		}
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
