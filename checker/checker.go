package checker

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
)

var hashes = []Hash{
	Md5,
	Sha1,
	Sha256,
	Sha512}

func CheckDir() {
	entries, _ := os.ReadDir(".")

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fmt.Printf("Checking: %s\n", entry.Name())
		for _, hashType := range hashes {
			getFileEntry(hashType, entry)
		}
	}
	fmt.Println("")
}

func CheckFile(filename string) {
	fmt.Printf("Checking: %s\n", filename)

	for _, hashType := range hashes {
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

		readChecksum(hashType, f)
	}
}

func getFileEntry(hash Hash, entry os.DirEntry) {
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

	readChecksum(hash, f)
}

func readChecksum(hashType Hash, f *os.File) {
	var h hash.Hash

	switch hashType {
	case Md5:
		h = md5.New()
	case Sha1:
		h = sha1.New()
	case Sha256:
		h = sha256.New()
	case Sha512:
		h = sha512.New()
	default:
		log.Fatalf("Unknown hash type: %s", hashType)
	}

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("- %-10s %x\n", hashType, h.Sum(nil))
}
