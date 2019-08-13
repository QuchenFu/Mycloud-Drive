package main

import (
	"filemd5"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	
	
	uploadclient()
	//listallclient()
	//downloadclient()
	//downloadpartclient()
	//updatepathclient()
	//deletefileclient()

}

func listallclient() {
	_, err := http.Get("http://localhost:8080/listall?")
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func getmetaclient() {
	fullpath := "/serverdata/text_2.txt"
	r, err := http.NewRequest("HEAD", "http://localhost:8080"+fullpath, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{}
	client.Do(r)
	return

}

func updatepathclient() {
	fullpath := "/serverdata/text_2.txt"
	newdirect := "/new/text_2.txt"
	r, err := http.NewRequest("PUT", "http://localhost:8080"+fullpath+"?newdirect="+newdirect, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{}
	client.Do(r)
	return
}

func deletefileclient() {
	fullpath := "/serverdata/text_2.txt"
	r, err := http.NewRequest("DELETE", "http://localhost:8080"+fullpath, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &http.Client{}
	client.Do(r)
	return
}

func uploadclient() {
	messages := make(chan string, 3)
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		file, err := os.Open("./clientdata/text_2.txt")
		if err != nil {
			log.Println(err)
			return
		}
		filestat, _ := file.Stat()
		messages <- strconv.FormatInt(filestat.Size(), 10)
		messages <- filemd5.FileMD5whole(file)
		part, err := m.CreateFormFile("theform", "./text_2.txt")
		if err != nil {
			log.Println(err)
			return
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			log.Println(err)
			return
		}
		return
	}()
	filename := "text_2.txt"
	fullpath := "/serverdata/" + filename
	filesize := <-messages
	md5 := <-messages
	req, err := http.NewRequest("POST", "http://localhost:8080"+fullpath, r)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Content-Type", m.FormDataContentType())
	req.Header.Add("Content-MD5", md5)
	req.Header.Add("Content-Size", filesize)
	client := &http.Client{}
	client.Do(req)
	return
}

func downloadclient() {
	fullpath := "/new/text_2.txt"
	resp, err := http.Get("http://localhost:8080/" + fullpath)
	filename := filepath.Base(fullpath)
	if err != nil {
		log.Println(err)
		return
	}
	newFile, err := os.Create("./clientdata/down-" + filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer newFile.Close()
	defer resp.Body.Close()
	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

//func downloadpartclient() {
//	start := 2
//	end := 3
//	md5 := "RjVq_lX6POqcvnOtRCytRw"
//	client := &http.Client{}
//	req, _ := http.NewRequest("GET", "http://localhost:8080/downloadpart?&md5="+md5, nil)
//	range_header := "bytes=" + strconv.Itoa(start) + "-" + strconv.Itoa(end)
//	req.Header.Add("Range", range_header)
//	resp, _ := client.Do(req)
//	filename := resp.Header.Get("filename")
//	newFile, err := os.Create("./clientdata/part-" + filename)
//	checkerr(err)
//	defer newFile.Close()
//	defer resp.Body.Close()
//	_, err = io.Copy(newFile, resp.Body)
//	checkerr(err)
//	return
//}

//func concurrentdownloadclient() {
//	var wg sync.WaitGroup
//	md5 := "73ydC-Iyij5pC-NMKbpqMQ"
//	resp, _ := http.Get("http://localhost:8080/meta?&md5="+md5)
//	//filename :=resp.Header.Get("filename")
//	filesize,_ :=strconv.ParseInt(resp.Header.Get("filesize"), 10, 64)
//	//log.Println(resp.Header.Get("filesize"))
//	limit := int64(10)
//	chunk := filesize / limit
//	remain := filesize % limit
//	for i := int64(0); i < limit; i++ {
//		min := chunk * i       // Min range
//		max := chunk * (i + 1) // Max range
//
//		if i == limit-1 {
//			max += remain // Add the remaining bytes in the last request
//		}
//		go func(min int64, max int64, i int64) {
//			client := &http.Client {}
//            req, err := http.NewRequest("GET", "http://localhost:8080/downloadpart?&md5="+md5, nil)
//            if err != nil {
//				log.Println(err)
//				return
//			}
//            range_header := "bytes=" + strconv.FormatInt(min, 10) +"-" + strconv.FormatInt(max-1, 10)
//            req.Header.Add("Range", range_header)
//            resp,err := client.Do(req)
//            if err != nil {
//				log.Println(err)
//				return
//			}
//            defer resp.Body.Close()
//            tmp, err := os.Create("./clientdata/"+strconv.FormatInt(i, 10)+".txt")
//            if err != nil {
//				log.Println(err)
//				return
//			}
//            defer tmp.Close()
//            io.Copy(tmp, resp.Body)
//            //reader, _ := ioutil.ReadAll(resp.Body)
//            //ioutil.WriteFile("./clientdata/"+strconv.FormatInt(i, 10)+".txt",reader, 0x777) // Write to the file i as a byte array
//            wg.Done()
//		}(min, max, i)
//	}
//	wg.Wait()
//	//log.Println(body[1])
//	return
//}

//func DownloadAllFileHandler(w http.ResponseWriter, req *http.Request) {
//	md5 := KeyParser(req, "md5")
//	//log.Println(filepath)
//	//	filename := pathmap[md5]
//	filename := meta.FileMetas[md5].FileName
//	//filepath := meta.FileMetas[md5].Location
//	//filesize := meta.FileMetas[md5].FileSize
//	log.Println(filename)
//	src, err := os.Open(filename)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	FileStat, _ := src.Stat()
//	size := FileStat.Size()
//	limit := int64(10)
//	chunk := size / limit
//	remain := size % limit
//	done := make(chan []byte, limit)
//	defer src.Close()
//	for i := int64(0); i < limit; i++ {
//		min := chunk * i       // Min range
//		max := chunk * (i + 1) // Max range
//
//		if i == limit-1 {
//			max += remain // Add the remaining bytes in the last request
//		}
//		src, err := os.Open(filename)
//		if err != nil {
//			log.Println(err)
//		}
//
//		go func(min int64, max int64, i int64) {
//			client := &http.Client {}
//            req, _ := http.NewRequest("GET", "http://localhost/rand.txt", nil)
//            range_header := "bytes=" + strconv.Itoa(min) +"-" + strconv.Itoa(max-1) // Add the data for the Range header of the form "bytes=0-100"
//            req.Header.Add("Range", range_header)
//            resp,_ := client.Do(req)
//            defer resp.Body.Close()
//            reader, _ := ioutil.ReadAll(resp.Body)
//            body[i] = string(reader)
//            ioutil.WriteFile(strconv.Itoa(i), []byte(string(body[i])), 0x777) // Write to the file i as a byte array
//            wg.Done()
//		}(min, max, i)
//	}
//	w.Write(<-done)
//	return
//
//	//wg.Wait()
//}
