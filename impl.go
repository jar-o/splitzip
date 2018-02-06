package splitzip

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	p "path"
	"sort"
)

type FileSize struct {
	Size int64
	Path string
}

// Interface to sort files largest -> smallest
type ByLargest []FileSize
func (f ByLargest) Len() int { return len(f) }
func (f ByLargest) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f ByLargest) Less(i, j int) bool { return f[i].Size > f[j].Size }

/*
	'savings' a guesstimate about the expected space saving. E.g. you may expect
	an average 10% space savings for the corpus you're processing. In that case,
	set it to 0.10. If you're unsure, or if you have a wide variety of file types,
	just leave it at 0.
*/
func groupBySize(files []FileSize, maxSize int64, savings float64) ([][]string, error) {
	var ret [][]string
	maxSize = maxSize + int64(float64(maxSize)*savings)

	sort.Sort(ByLargest(files))
	for i, e := range files {
			if e.Size == -1 { continue }
			// TODO a reasonable percentage over required size, maybe?
			if e.Size > maxSize {
				return ret, errors.New("Single file is too large to meet requirements.")
			}
			files[i].Size = -1
			
			group := []string{e.Path}
			d := maxSize - e.Size
			
			j := i + 1
			s := e.Size
			for ; j < len(files); j++ {
				if files[j].Size == -1 { continue }
					if files[j].Size <= d {
					s += files[j].Size
					d = maxSize - s
					group = append(group, files[j].Path)
					files[j].Size = -1
				}
				if s > maxSize { break }
			}
			ret = append(ret, group)
	}
	return ret, nil
}

func loadFiles(path string) ([]FileSize, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	path, _ = filepath.Abs(path)

	var ret []FileSize
	for _, file := range files {
		if file.IsDir() { continue }
		f := FileSize{
			Size: file.Size(),
			Path: p.Join(path, file.Name()),
		}
		ret = append(ret, f)
	}
	return ret, nil
}

func splitZip(groupedFiles [][]string, outpath string, prefix string) error {
	outpath, _ = filepath.Abs(outpath)

	for i, g := range groupedFiles {
		zipfile, err := os.Create(p.Join(outpath, prefix + fmt.Sprintf("%d", i) + ".zip"))
		if err != nil {
			return err
		}
		zw := zip.NewWriter(zipfile)

		for _, f := range g {
			zf, err := os.Open(f)
			if err != nil {
				return err
			}

			// Get the file information
			info, err := zf.Stat()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// Change to deflate to gain better compression
			// see http://golang.org/pkg/archive/zip/#pkg-constants
			header.Method = zip.Deflate

			writer, err := zw.CreateHeader(header)
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, zf)
			if err != nil {
				return err
			}
			zf.Close()
		}

		zw.Close()
		zipfile.Close()

	}
	return nil
}

func Zip(inpath string, outpath string, prefix string, maxSize int64, savings float64) error {
	f, err := loadFiles(inpath)
	if err != nil { return err }

	g, err := groupBySize(f, maxSize, savings)
	if err != nil { return err }

	return splitZip(g, outpath, prefix)
}
