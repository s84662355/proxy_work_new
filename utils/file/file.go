package file

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func DeleteFile(srcFileName string) error {
	return os.Remove(srcFileName)
}

func DeleteDir(deleteDir string) error {
	return os.RemoveAll(deleteDir)
}

func GetFileSize(filePath string) (uint64, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	// 获取文件大小
	return uint64(fileInfo.Size()), nil
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func CopyFile(srcFileName string, dstFileName string) error {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstFileName)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func handleSymlink(srcPath, dstPath string) error {
	// 实现符号链接的处理逻辑
	originalTarget, err := os.Readlink(srcPath)
	if err != nil {
		fmt.Printf("Error reading symlink '%s': %s\n", srcPath, err)
		return err
	}

	// 可以选择重新创建符号链接或复制符号链接指向的实际文件/目录
	fmt.Printf("Original symlink target: '%s'\n", originalTarget)
	// 创建相同的符号链接在目标位置
	err = os.Symlink(originalTarget, dstPath)
	if err != nil {
		fmt.Printf("Error creating symlink '%s' -> '%s': %s\n", dstPath, originalTarget, err)
		return err
	}
	return nil
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 确保目标路径不是目录
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	if srcInfo.IsDir() {
		return fmt.Errorf("source is a directory")
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return dstFile.Sync()
}

// CopyDir 递归复制目录及其内容
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		fileInfo, err := os.Lstat(srcPath)
		if err != nil {
			fmt.Println("get path lstat fail -----", err)
			return err
		}

		if fileInfo.Mode()&os.ModeSymlink != 0 {
			fmt.Println("find the link------------", srcPath)
			if err := handleSymlink(srcPath, dstPath); err != nil {
				return err
			}
		} else if fileInfo.IsDir() {
			// fmt.Printf("is dir and copy from '%s' to '%s' \n", srcPath, dstPath)
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				// fmt.Printf("Copying from '%s' to '%s' [IsDir: %v, Mode: %v]\n", srcPath, dstPath, fileInfo.IsDir(), fileInfo.Mode())
				return err
			}
		}
	}
	return nil
}

func GetFileType(name string) (string, string, string, error) {
	filet, err := os.Open(name)
	if err != nil {
		return "", "", "", err
	}
	defer filet.Close()

	// 读取文件的前 512 个字节
	buffer := make([]byte, 512)
	n, err := filet.Read(buffer)
	if err != nil {
		return "", "", "", err
	}
	// 调用 http.DetectContentType 方法判断文件类型
	// 实际上，如果字节数超过 512，该函数也只会使用前 512 个字节
	contentType := http.DetectContentType(buffer[:n])
	return filepath.Base(name), filepath.Ext(name), contentType, nil
}

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetFileCreateTime(filename string) int64 {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	return fileinfo.ModTime().Unix()
}

func GetFileData(f string) []byte {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return nil
	}
	return content
}

func GetFileJsonData[res any](f string, r res) error {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, r); err != nil {
		return err
	}
	return nil
}

func WriteJsonFile(f string, j interface{}) error {
	jsonfile, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer jsonfile.Close()

	if data, err := json.Marshal(j); err != nil {
		return err
	} else {
		if _, err := jsonfile.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func IsDirExists(path string) bool {
	_, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// FileExistsInSubdir 检查文件是否存在于指定目录或其任何子目录中
func FileExistsInSubdir(rootDir, filenameToFind string) (bool, error) {
	var found bool
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // 传递错误
		}
		if !d.IsDir() && filepath.Base(path) == filenameToFind {
			found = true
			return fmt.Errorf("file found") // 找到文件，提前终止遍历
		}
		return nil
	})

	if err != nil && !found {
		return false, err // 遍历中出现错误
	}

	return found, nil // 返回找到的状态和nil错误
}
