package mime

import (
	"bufio"
	"embed"
	"mime"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed etc/mime.types
var fsys embed.FS

type FileTypeDetector interface {
	DetectFileType(path string) string
}

var fileTypeDetector = NewFileTypeDetector()

func DetectFileType(path string) string {
	return fileTypeDetector.DetectFileType(path)
}

type fileTypeDetectorImpl struct {
	sync.RWMutex
	extensions map[string]string
}

func NewFileTypeDetector() FileTypeDetector {
	detector := &fileTypeDetectorImpl{
		extensions: make(map[string]string),
	}
	detector.loadMimeFile("etc/mime.types")
	return detector
}

func (d *fileTypeDetectorImpl) AddTypeExtentions(mimeType string, exts ...string) error {
	d.Lock()
	defer d.Unlock()
	t, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return err
	}
	for _, ext := range exts {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		ext = strings.ToLower(ext)

		d.extensions[ext] = t
	}
	return nil
}

func (d *fileTypeDetectorImpl) loadMimeFile(filename string) {

	f, err := fsys.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) <= 1 || fields[0][0] == '#' {
			continue
		}
		mimeType := fields[0]
		for _, ext := range fields[1:] {
			if ext[0] == '#' {
				break
			}
			err := d.AddTypeExtentions(mimeType, "."+ext)
			if err != nil {
				panic(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func (d *fileTypeDetectorImpl) DetectFileType(path string) string {
	ext := filepath.Ext(path)
	if ext != "" {
		ext = strings.ToLower(ext)
		d.RLock()
		defer d.RUnlock()
		return d.extensions[ext]
	} else {
		return ""
	}
}
