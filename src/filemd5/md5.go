package filemd5

import (
	"crypto/md5"
	"encoding/base64"
	"io"
	"os"
)

func FileMD5(file *os.File, start int64, end int64) string {
	_md5 := md5.New()
	file.Seek(start, 0)

	io.CopyN(_md5, file, end-start+1)
	md5Value := _md5.Sum(nil)
	//log.Printf(base64.RawURLEncoding.EncodeToString(md5Value))
	return base64.RawURLEncoding.EncodeToString(md5Value)
}

func FileMD5whole(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	md5Value := _md5.Sum(nil)
	//log.Printf(base64.RawURLEncoding.EncodeToString(md5Value))
	return base64.RawURLEncoding.EncodeToString(md5Value)
}

//
//func main(){
//	path := "/servg/er/data/text_2.txt"
//
//
//}
