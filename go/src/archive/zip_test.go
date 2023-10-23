package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tryOpen 打开一个文件，如果文件不存在，则尝试创建
func tryOpen(path string, flag int) (*os.File, error) {
	fabs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(fabs, flag, os.ModePerm)
	if os.IsNotExist(err) {
		// NOTE[Dokiy] 2023/10/23: 创建的文件夹需要是0766才能够在文件夹中创建文件
		err = os.MkdirAll(filepath.Dir(fabs), os.ModePerm)
		if err != nil {
			return nil, err
		}
		return os.OpenFile(fabs, flag, os.ModePerm)
	}
	return f, err
}
func TestZip_Write(t *testing.T) {
	// 创建一个缓冲区用来保存压缩文件内容
	buf := new(bytes.Buffer)
	// 创建一个压缩文档
	w := zip.NewWriter(buf)
	// 将文件加入压缩文档
	var files = []struct {
		Name, Body string
	}{
		{"4399/4399_1.txt", "https://www.4399.com/"},
		{"4399/4399_2.txt", "https://www.4399.com/"},
		{"baidu.txt", "https://www.baidu.com/"},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			t.Fatal(err)
		}
	}

	// 关闭压缩文档
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	// 将压缩文档内容写入文件
	f, err := tryOpen("testdata_local/TestZip_Write.zip", os.O_CREATE|os.O_WRONLY)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := buf.WriteTo(f); err != nil {
		t.Fatal(err)
	}
}

func archivePrepareFile() error {
	var err error
	f2, err := tryOpen("testdata_local/TestZip_Archive/sub_dir/2", os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return err
	}
	if _, err := f2.Write([]byte("2")); err != nil {
		return err
	}
	defer f2.Close()

	f1, err := tryOpen("testdata_local/TestZip_Archive/1", os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return err
	}
	if _, err := f1.Write([]byte("1")); err != nil {
		return err
	}
	defer f1.Close()

	return nil
}
func TestZip_Archive(t *testing.T) {
	if err := archivePrepareFile(); err != nil {
		t.Fatal(err)
	}

	var dirname = "testdata_local/TestZip_Archive"
	var filename = "testdata_local/TestZip_Archive.zip"
	// 创建zip文件
	zipFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766)
	if err != nil {
		t.Fatal(err)
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// 遍历Zip
	_ = filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var header *zip.FileHeader
		if header, err = zip.FileInfoHeader(info); err != nil {
			return err
		}

		header.Name = strings.TrimLeft(strings.ReplaceAll(path, dirname, ""), "/")
		if info.IsDir() {
			header.Name += "/"
		} else {
			// 如果是文件则压缩
			header.Method = zip.Deflate
		}

		var writer io.Writer
		if writer, err = archive.CreateHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.OpenFile(path, os.O_RDONLY, 0666)
			if err != nil {
				return err
			}

			defer file.Close()
			_, _ = io.Copy(writer, file)
		}

		return err
	})
}

func readPrepareFile() error {
	// 创建一个缓冲区用来保存压缩文件内容
	buf := new(bytes.Buffer)
	// 创建一个压缩文档
	w := zip.NewWriter(buf)
	// 将文件加入压缩文档
	var files = []struct {
		Name, Body string
	}{
		{"4399/4399_1.txt", "https://www.4399.com/"},
		{"4399/4399_2.txt", "https://www.4399.com/"},
		{"baidu.txt", "https://www.baidu.com/"},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return err
		}
	}

	// 关闭压缩文档
	if err := w.Close(); err != nil {
		return err
	}
	// 将压缩文档内容写入文件
	f, err := tryOpen("testdata_local/TestZip_Read.zip", os.O_CREATE|os.O_WRONLY)
	if err != nil {
		return err
	}

	if _, err := buf.WriteTo(f); err != nil {
		return err
	}

	return nil
}
func TestZip_Read(t *testing.T) {
	if err := readPrepareFile(); err != nil {
		t.Fatal(err)
	}

	// 打开一个zip格式文件
	r, err := zip.OpenReader("testdata_local/TestZip_Read.zip")
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer r.Close()

	// 迭代压缩文件中的文件，打印出文件中的内容
	for _, f := range r.File {
		fmt.Printf("文件名(%s) 内容:", f.Name)
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.CopyN(os.Stdout, rc, int64(f.UncompressedSize64))
		if err != nil {
			t.Fatal(err)
		}
		rc.Close()
		fmt.Println()
	}
}
