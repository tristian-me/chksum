package checker

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
)

var hashes = []Hash{
	Md5,
	Sha1,
	Sha256,
	Sha512}

var hashFactories = map[Hash]func() hash.Hash{
	Md5:    md5.New,
	Sha1:   sha1.New,
	Sha256: sha256.New,
	Sha512: sha512.New,
}

func CheckDir() error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if err := CheckFile(entry.Name()); err != nil {
			return err
		}
	}
	fmt.Println("")
	return nil
}

func CheckFile(filename string) error {
	fmt.Printf("Checking: %s\n", filename)

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %s", filename)
	}

	defer func() {
		_ = f.Close()
	}()

	sums, err := computeChecksums(f, hashes)
	if err != nil {
		return fmt.Errorf("error computing checksums: %s", filename)
	}

	for _, ht := range hashes {
		fmt.Printf("- %-10s %x\n", ht, sums[ht])
	}

	return nil
}

func computeChecksums(r io.Reader, types []Hash) (map[Hash][]byte, error) {
	hashers := make([]hash.Hash, 0, len(types))
	writers := make([]io.Writer, 0, len(types))
	byType := make(map[Hash]hash.Hash, len(types))

	for _, t := range types {
		newHash, ok := hashFactories[t]
		if !ok {
			return nil, fmt.Errorf("unknown hash type: %s", t)
		}
		h := newHash()
		hashers = append(hashers, h)
		writers = append(writers, h)
		byType[t] = h
	}

	if _, err := io.Copy(io.MultiWriter(writers...), r); err != nil {
		return nil, err
	}

	out := make(map[Hash][]byte, len(types))
	for _, t := range types {
		out[t] = byType[t].Sum(nil)
	}
	return out, nil
}
