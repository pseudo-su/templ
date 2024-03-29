// Code generated by "esc -o datadog_assets.go -pkg main stencil-plugin-datadog.base.yaml"; DO NOT EDIT.

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

	"/stencil-plugin-datadog.base.yaml": {
		name:    "stencil-plugin-datadog.base.yaml",
		local:   "stencil-plugin-datadog.base.yaml",
		size:    9627,
		modtime: 1562855558,
		compressed: `
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
sqCbJQAA
`,
	},
}

var _escDirs = map[string][]os.FileInfo{}
