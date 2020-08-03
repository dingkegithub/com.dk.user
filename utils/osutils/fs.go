package osutils
import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func Mkdir(p string, reverse bool) error {
	if reverse {
		return os.Mkdir(p, 0766)
	} else {
		return os.MkdirAll(p, 0766)
	}
}

func Touch(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()
	return nil
}