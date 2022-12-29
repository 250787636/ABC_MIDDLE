package minios

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
)

type Client struct {
	client *minio.Client
}

func NewMinioClient(endpoint, accessKeyID, secretAccesskey string, useSSL bool) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccesskey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Client{client: minioClient}, nil
}

func (c *Client) Raw() *minio.Client {
	return c.client
}

func (c *Client) GetObject(bucketName, objectName string) (*minio.Object, error) {
	return c.client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
}

func (c *Client) GetFileBytes(bucketName, objectName string) ([]byte, error) {
	object, err := c.GetObject(bucketName, objectName)
	if err != nil {
		return nil, errors.WithMessagef(err, "获取minio文件失败,bucket=%s,object=%s", bucketName, objectName)
	}
	defer func() { _ = object.Close() }()
	byteArr, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return byteArr, nil
}

func (c *Client) FGetObject(bucketName, objectName, filePath string) error {
	return c.client.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
}

func (c *Client) FPutObject(bucketName, filePath string) error {
	return c.FPutObjectWithName(bucketName, filepath.Base(filePath), filePath)
}

func (c *Client) FPutObjectWithName(bucketName, objectName, filePath string) error {
	_, err := c.client.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// object文件不存在则进行上传
func (c *Client) FPutObjectNX(buckentName, objectName, filePath string) error {
	_, err := c.client.StatObject(context.Background(), buckentName, objectName, minio.StatObjectOptions{})
	if err == nil {
		return nil
	}
	if err = c.FPutObject(buckentName, filePath); err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveObject(bucketName, objectName string) error {
	return c.client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
}

// http文件上传
func (c *Client) ServeUpload(buckentName, objectName string, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	_, err = c.client.PutObject(context.Background(), buckentName, objectName, file, fileHeader.Size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// http文件下载
func (c *Client) ServeContent(w http.ResponseWriter, req *http.Request, bucketName, objectName string) error {
	return c.ServeContentWithFilename(w, req, bucketName, objectName, "")
}

// http文件下载 指定文件名
func (c *Client) ServeContentWithFilename(w http.ResponseWriter, req *http.Request, buckentName, objectName, filename string) error {
	object, err := c.GetObject(buckentName, objectName)
	if err != nil {
		return err
	}
	objectStat, err := object.Stat()
	if err != nil {
		return err
	}
	if filename == "" {
		filename = url.PathEscape(filepath.Base(objectName))
	} else {
		filename = url.PathEscape(filename)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, filename))
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	http.ServeContent(w, req, filename, objectStat.LastModified, object)
	return nil
}
