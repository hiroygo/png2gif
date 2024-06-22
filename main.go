package main

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/png"
	"os"
	"path/filepath"
)

func filterFiles(dir, ext string) ([]string, error) {
	var pathes []string
	err := filepath.Walk(dir, func(visitPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(f.Name()) == ext {
			pathes = append(pathes, visitPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return pathes, nil
}

func replaceExt(path, newExt string) string {
	fileName := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
	fileName += newExt
	return filepath.Join(filepath.Dir(path), fileName)
}

func readImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func writeGif(img image.Image, writePath string) (rerr error) {
	f, err := os.Create(writePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			rerr = err
		}
	}()

	if err := gif.Encode(f, img, nil); err != nil {
		return err
	}

	return
}

func main() {
	pathes, err := filterFiles("./", ".png")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, path := range pathes {
		img, err := readImage(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read %q: %v\n", path, err)
			continue
		}

		writePath := replaceExt(path, ".gif")
		if err := writeGif(img, writePath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %q: %v\n", writePath, err)
			continue
		}
	}
}
