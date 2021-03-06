package img_server

import "os"

// IsDir is_dir()
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	fm := fd.Mode()
	return fm.IsDir(), nil
}

// IsFile is_file()
func IsFile(filename string) bool {
	f, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}
