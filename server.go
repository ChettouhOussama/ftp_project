package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/goftp/server"
)

// MyDriver implements the server.Driver interface
type MyDriver struct {
	RootPath string
}

// ChangeDir implements server.Driver.
func (d *MyDriver) ChangeDir(string) error {
	panic("unimplemented")
}

// DeleteDir implements server.Driver.
func (d *MyDriver) DeleteDir(string) error {
	panic("unimplemented")
}

// DeleteFile implements server.Driver.
func (d *MyDriver) DeleteFile(string) error {
	panic("unimplemented")
}

// Init implements server.Driver.
func (d *MyDriver) Init(*server.Conn) {
	panic("unimplemented")
}

// MakeDir implements server.Driver.
func (d *MyDriver) MakeDir(string) error {
	panic("unimplemented")
}

// PutFile implements server.Driver.
func (d *MyDriver) PutFile(string, io.Reader, bool) (int64, error) {
	panic("unimplemented")
}

// Rename implements server.Driver.
func (d *MyDriver) Rename(string, string) error {
	panic("unimplemented")
}

// MyFileInfo implements the server.FileInfo interface
type MyFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

// Group implements server.FileInfo.
func (f *MyFileInfo) Group() string {
	panic("unimplemented")
}

// Owner implements server.FileInfo.
func (f *MyFileInfo) Owner() string {
	panic("unimplemented")
}

// Sys implements server.FileInfo.
func (f *MyFileInfo) Sys() any {
	panic("unimplemented")
}

func (f *MyFileInfo) Name() string       { return f.name }
func (f *MyFileInfo) Size() int64        { return f.size }
func (f *MyFileInfo) Mode() os.FileMode  { return f.mode }
func (f *MyFileInfo) ModTime() time.Time { return f.modTime }
func (f *MyFileInfo) IsDir() bool        { return f.isDir }

// NewDriver creates a new driver instance
func (d *MyDriver) NewDriver() (server.Driver, error) {
	return d, nil
}

// Stat retrieves file or directory information
func (d *MyDriver) Stat(path string) (server.FileInfo, error) {
	absPath := filepath.Join(d.RootPath, path)
	info, err := os.Stat(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	return &MyFileInfo{
		name:    info.Name(),
		size:    info.Size(),
		mode:    info.Mode(),
		modTime: info.ModTime(),
		isDir:   info.IsDir(),
	}, nil
}

// ListDir lists the contents of a directory
func (d *MyDriver) ListDir(path string, callback func(server.FileInfo) error) error {
	absPath := filepath.Join(d.RootPath, path)
	entries, err := os.ReadDir(absPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		fileInfo := &MyFileInfo{
			name:    info.Name(),
			size:    info.Size(),
			mode:    info.Mode(),
			modTime: info.ModTime(),
			isDir:   entry.IsDir(),
		}

		if err := callback(fileInfo); err != nil {
			return err
		}
	}
	return nil
}

// GetFile retrieves a file from the server
func (d *MyDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	absPath := filepath.Join(d.RootPath, path)
	file, err := os.Open(absPath)
	if err != nil {
		return 0, nil, err
	}

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		file.Close()
		return 0, nil, err
	}
	return offset, file, nil
}

func main() {
	rootDir := "./ftpdir"

	if err := os.MkdirAll(rootDir, 0755); err != nil {
		log.Fatalf("Failed to create root directory: %v", err)
	}

	opts := &server.ServerOpts{
		Factory: &MyDriver{RootPath: rootDir},
		Auth:    &server.SimpleAuth{Name: "user", Password: "pass"},
		Port:    2121,
	}

	ftpServer := server.NewServer(opts)

	log.Println("FTP server running on port 2121...")
	if err := ftpServer.ListenAndServe(); err != nil {
		log.Fatalf("FTP server stopped: %v", err)
	}
}
