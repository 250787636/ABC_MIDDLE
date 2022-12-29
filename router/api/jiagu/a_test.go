package jiagu

import (
	"example.com/m/pkg/log"
	"fmt"
	"github.com/google/uuid"
	"github.com/levigross/grequests"
	"os"
	"path"
	"testing"
)

func TestC(t *testing.T) {
	resp, err := grequests.Get("https://www.python.org/ftp/python/pc/32python.zip",nil)
	if err !=nil {
		fmt.Println(err.Error())
	}
	if err != nil {
		log.Info(err)
		fmt.Println(err)
	}
	fileSourceName := path.Base(resp.RawResponse.Request.URL.Path)
	uid := uuid.New()
	randomPath := fmt.Sprintf("/normal/%s%s", uid, path.Ext(fileSourceName))
	defer resp.Close()
	os.MkdirAll("/normal",0766)
	err = resp.DownloadToFile(randomPath)
	if err !=nil {
		fmt.Println(err)
	}
	f, _ := os.Open(randomPath)
	go func() {
		_ = os.Remove(randomPath)
	}()


	fmt.Println(f)
}
