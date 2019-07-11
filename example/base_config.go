// Code generated by "esc -o base_config.go -pkg main base-datadog.config.yaml"; DO NOT EDIT.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return fis[0:limit], nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/base-datadog.config.yaml": {
		name:    "base-datadog.config.yaml",
		local:   "base-datadog.config.yaml",
		size:    9610,
		modtime: 1562852286,
		compressed: `
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
2/jO7px/BwAA//+C1wkbiiUAAA==
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
