package utils

import (
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// Stat — Gives information about a file
func Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

// PathInfo — Returns information about a file path
// -1: all; 1: dirname; 2: basename; 4: extension; 8: filename
// Usage:
// PathInfo("/home/go/path/src/php2go/php2go.go", 1|2|4|8)
func PathInfo(path string, options int) map[string]string {
	if options == -1 {
		options = 1 | 2 | 4 | 8
	}

	info := make(map[string]string)
	if (options & 1) == 1 {
		info["dirname"] = filepath.Dir(path)
	}

	if (options & 2) == 2 {
		info["basename"] = filepath.Base(path)
	}

	if ((options & 4) == 4) || ((options & 8) == 8) {
		basename := ""
		if (options & 2) == 2 {
			basename, _ = info["basename"]
		} else {
			basename = filepath.Base(path)
		}

		p := strings.LastIndex(basename, ".")
		filename, extension := "", ""
		if p > 0 {
			filename, extension = basename[:p], basename[p+1:]
		} else if p == -1 {
			filename = basename
		} else if p == 0 {
			extension = basename[p+1:]
		}

		if (options & 4) == 4 {
			info["extension"] = extension
		}

		if (options & 8) == 8 {
			info["filename"] = filename
		}
	}

	return info
}

// FileExists — Checks whether a file or directory exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsFile — Tells whether the filename is a regular file
func IsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// IsDir — Tells whether the filename is a directory
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	fm := fd.Mode()
	return fm.IsDir(), nil
}

// FileSize — Gets file size
func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}

	return info.Size(), nil
}

// FilePutContents — Write data to a file
func FilePutContents(filename string, data string, mode os.FileMode) error {
	return ioutil.WriteFile(filename, []byte(data), mode)
}

// FileGetContents — Reads entire file into a string
func FileGetContents(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	return string(data), err
}

// Unlink — Deletes a file
func Unlink(filename string) error {
	return os.Remove(filename)
}

// Delete — See Unlink()
func Delete(filename string) error {
	return os.Remove(filename)
}

// Copy — Copies file
func Copy(source, dest string) (bool, error) {
	fd1, err := os.Open(source)
	if err != nil {
		return false, err
	}
	defer fd1.Close()

	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer fd2.Close()

	_, e := io.Copy(fd2, fd1)
	if e != nil {
		return false, e
	}

	return true, nil
}

// IsReadable — Tells whether a file exists and is readable
func IsReadable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return false
	}

	return true
}

// IsWriteable — Tells whether the filename is writable
func IsWriteable(filename string) bool {
	_, err := syscall.Open(filename, syscall.O_WRONLY, 0)
	if err != nil {
		return false
	}

	return true
}

// Rename — Renames a file or directory
func Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// Touch — Sets access and modification time of file
func Touch(filename string) (bool, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}
	fd.Close()

	return true, nil
}

// Mkdir — Makes directory
func Mkdir(filename string, mode os.FileMode) error {
	return os.Mkdir(filename, mode)
}

// GetCwd — Gets the current working directory
func GetCwd() (string, error) {
	dir, err := os.Getwd()
	return dir, err
}

// RealPath — Returns canonicalized absolute pathname
func RealPath(path string) (string, error) {
	return filepath.Abs(path)
}

// BaseName — Returns trailing name component of path
func BaseName(path string) string {
	return filepath.Base(path)
}

// Chmod — Changes file mode
func Chmod(filename string, mode os.FileMode) bool {
	return os.Chmod(filename, mode) == nil
}

// Chown — Changes file owner
func Chown(filename string, uid, gid int) bool {
	return os.Chown(filename, uid, gid) == nil
}

// Fclose — Closes an open file pointer
func Fclose(handle *os.File) error {
	return handle.Close()
}

// FileMtime — Gets file modification time
func FileMtime(filename string) (int64, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	fileinfo, err := fd.Stat()
	if err != nil {
		return 0, err
	}

	return fileinfo.ModTime().Unix(), nil
}

// FgetCsv — Gets line from file pointer and parse for CSV fields
func FgetCsv(handle *os.File, length int, delimiter rune) ([][]string, error) {
	reader := csv.NewReader(handle)
	reader.Comma = delimiter
	// TODO length limit
	return reader.ReadAll()
}

// Glob — Find pathnames matching a pattern
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
