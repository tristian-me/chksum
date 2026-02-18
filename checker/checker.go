package checker

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"os"
)

func CheckDir() {
	items, _ := os.ReadDir(".")

	for _, item := range items {
		if item.IsDir() {
			continue
		}
		getFileEntry(item)
	}
}

func CheckFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", filename, err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Unable to open %s", filename)
		}
	}(f)

	fmt.Printf("Checking: %s\n", filename)
	readChecksum(filename, f)
}

func getFileEntry(entry os.DirEntry) {
	f, err := os.Open(entry.Name())
	if f == nil || entry.IsDir() {
		return
	}

	if err != nil {
		fmt.Printf("Unable to open %s", entry.Name())
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Unable to close %s", entry.Name())
		}
	}(f)

	readChecksum(entry.Name(), f)
}

func readChecksum(filename string, f *os.File) {
	fmt.Printf("Checking: %s\n", filename)

	hasherMd5 := md5.New()
	hasherSha1 := sha1.New()
	hasherSha256 := sha256.New()
	hasherSha512 := sha512.New()

	if _, err := io.Copy(hasherMd5, f); err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(hasherSha1, f); err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(hasherSha256, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("- MD5:    %x\n", hasherMd5.Sum(nil))
	fmt.Printf("- SHA1:   %x\n", hasherSha1.Sum(nil))
	fmt.Printf("- SHA256: %x\n", hasherSha256.Sum(nil))
	fmt.Printf("- SHA256: %x\n", hasherSha512.Sum(nil))
	fmt.Println("")
}
