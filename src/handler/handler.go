package handler

import (
	"encoding/json"
	"filemd5"
	"io"
	"log"
	"meta"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFileHandler(w http.ResponseWriter, req *http.Request) {
	fullpath := req.URL.Path
	directory := filepath.Dir(fullpath)
	filename := filepath.Base(fullpath)
	truemd5 := req.Header.Get("Content-MD5")
	filesize := req.Header.Get("Content-Size")
	if filename == "" || directory == "" || truemd5 == "" || filesize == "" || fullpath == "" {
		http.Error(w, "400 param is wrong.", http.StatusBadRequest)
	}
	fileMeta := meta.FileMeta{
		FileMD5:   truemd5,
		FileName:  filename,
		Directory: directory,
		FileSize:  filesize,
		FullPath:  fullpath,
	}
	thelist, _ := meta.GetFileMetaList()
	for _, v := range thelist {
		if v.FileMD5 == truemd5 && v.FullPath != fullpath {
			//meta.UpdateFileMeta(fileMeta, fullpath)
			meta.UploadMeta(fileMeta)
			meta.AddFileLife(truemd5)
			w.Write([]byte("Already exist"))
			log.Println("Already exist")
			return
		}
		return
	}
	file, _, err := req.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	newFile, err := os.Create("./tmp/1" + truemd5)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	defer file.Close()
	io.Copy(newFile, file)
	_, err = newFile.Seek(0, 0)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	if truemd5 == filemd5.FileMD5whole(newFile) {
		meta.UploadMeta(fileMeta)
		meta.AddFileLife(truemd5)
		log.Println("upload succeed!")
		return
	}
	http.Error(w, "500 file broken.", http.StatusInternalServerError)
	return
}

func DownloadFileHandler(w http.ResponseWriter, req *http.Request) {
	fullpath := req.URL.Path
	directory := filepath.Dir(fullpath)
	filename := filepath.Base(fullpath)
	if filename == "" || directory == "" || fullpath == "" {
		http.Error(w, "400 param is wrong.", http.StatusBadRequest)
		return
	}
	fmeta, _ := meta.GetFileMeta(fullpath)
	if fmeta == nil {
		http.Error(w, "404 file does not exist.", http.StatusNotFound)
		return
	}
	truemd5 := fmeta.FileMD5
	filesize := fmeta.FileSize
	src, err := os.Open("./tmp/" + truemd5)
	if err != nil {
		log.Println(err)
		w.Write([]byte("something wrong"))
		return
	}
	defer src.Close()
	w.Header().Set("Content-Size", filesize)
	_, err = io.Copy(w, src)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	log.Println("Download succeed")
	return
}

func ListAllHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Only GET method supported!", http.StatusNotFound)
		return
	}
	alldata, err := meta.GetFileMetaList()
	if alldata == nil {
		http.Error(w, "404 file does not exist.", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(alldata)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}

func GetMetaHandler(w http.ResponseWriter, req *http.Request) {
	fullpath := req.URL.Path
	log.Println(fullpath)
	directory := filepath.Dir(fullpath)
	filename := filepath.Base(fullpath)
	if filename == "" || directory == "" || fullpath == "" {
		log.Println("wrong")
		http.Error(w, "400 param is wrong.", http.StatusBadRequest)
	}
	fmeta, err := meta.GetFileMeta(fullpath)
	if fmeta == nil {
		http.Error(w, "404 file does not exist.", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-MD5", fmeta.FileMD5)
	w.Header().Set("Directory", fmeta.Directory)
	w.Header().Set("Filename", fmeta.FileName)
	w.Header().Set("Content-Size", fmeta.FileSize)
	return
}

func DeleteFileHandler(w http.ResponseWriter, req *http.Request) {
	fullpath := req.URL.Path
	if fullpath == "" {
		http.Error(w, "400 param is wrong.", http.StatusBadRequest)
		return
	}
	fmeta, err := meta.GetFileMeta(fullpath)
	if fmeta == nil {
		http.Error(w, "404 file does not exist.", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	meta.MinusFileLife(fmeta.FileMD5)
	meta.DeleteFileMeta(fullpath)
	log.Println("deleted!")
	return
}

func UpdatePathHandler(w http.ResponseWriter, req *http.Request) {
	newdirect := KeyParser(req, "newdirect")
	fullpath := req.URL.Path
	directory := filepath.Dir(newdirect)
	filename := filepath.Base(newdirect)
	fmeta, err := meta.GetFileMeta(fullpath)
	if fmeta == nil {
		http.Error(w, "404 file does not exist.", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "500 server error.", http.StatusInternalServerError)
		return
	}
	filesize := fmeta.FileSize
	fileMD5 := fmeta.FileMD5
	if fullpath == "" || newdirect == "" || directory == "" || filename == "" || filesize == "" || fileMD5 == "" {
		http.Error(w, "400 param is wrong.", http.StatusBadRequest)
		return
	}
	fileMeta := meta.FileMeta{
		FileMD5:   fileMD5,
		FileName:  filename,
		Directory: directory,
		FileSize:  filesize,
		FullPath:  newdirect,
	}
	meta.DeleteFileMeta(fullpath)
	meta.UploadMeta(fileMeta)
	log.Println("updated!")
	return
}

func KeyParser(req *http.Request, keyname string) string {

	keynames, ok := req.URL.Query()[keyname]
	if !ok || len(keynames[0]) < 1 {
		log.Println("Url Param " + keyname + "is missing")
		return ""
	}
	key := keynames[0]
	return key
}

//func DownloadFilePartHandler(w http.ResponseWriter, req *http.Request) {
//
//	var start, end int64
//	fmt.Sscanf(req.Header.Get("Range"), "bytes=%d-%d", &start, &end)
//	md5 := KeyParser(req, "md5")
//	if _, ok := meta.FileMetas[md5]; ok {
//
//		filename := meta.FileMetas[md5].FileName
//		filepath := meta.FileMetas[md5].Location
//		filesize := meta.FileMetas[md5].FileSize
//		if start < 0 || start > filesize || end < 0 || end > filesize {
//			log.Println("range out of order")
//			return
//		}
//		if end == 0 {
//			end = filesize - 1
//		}
//		log.Println(filename)
//		src, err := os.Open("./tmp/" + md5)
//		checkerr(err)
//
//		_, _ = src.Seek(start, 0)
//		defer src.Close()
//		w.Header().Set("Content-MD5", md5)
//		w.Header().Set("filepath", filepath)
//		w.Header().Set("filename", filename)
//		w.Header().Set("filesize", strconv.FormatInt(filesize, 10))
//		_, err = io.CopyN(w, src, end-start)
//		checkerr(err)
//		log.Println("Download part succeed")
//		return
//	} else {
//
//		log.Println("File doesn't exist!")
//		return
//	}
//	return
//}
