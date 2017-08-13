package model

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

func TestConv(t *testing.T) {
	dirs := []string{"1", "2", "3", "4"}
	pretty.Println("....")
	for _, dir := range dirs {
		pretty.Println(dir)
		dd, err := ioutil.ReadDir("update" + "/" + dir)
		if err != nil {
			t.Error(err)
		}

		tid, err := strconv.Atoi(dir)
		if err != nil {
			t.Error(err)
		}

		tids := []uint64{uint64(tid)}

		for _, d := range dd {
			filename := "update/" + dir + "/" + d.Name()
			pretty.Println(filename)
			if path.Ext(filename) == ".doc" {
				res, err := tUpload(filename)
				if err != nil {
					t.Error(err)
				}

				if err := CreateBook(strings.TrimSuffix(d.Name(), path.Ext(d.Name())), []byte(res), tids); err != nil {
					t.Error(err)
				}
			}
		}
	}

}

func tUpload(filename string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", err
	}

	//这里必须Close，否则不会向bodyBuf写入boundary分隔符
	err = bodyWriter.Close()
	if err != nil {
		return "", err
	}
	contentType := bodyWriter.FormDataContentType()

	req, err := http.NewRequest("POST", "http://localhost:8080/conv", bodyBuf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
