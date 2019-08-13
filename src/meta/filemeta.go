package meta

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path/filepath"
)

type FileMeta struct {
	FileMD5   string
	FileName  string
	FileSize  string
	Directory string
	FullPath  string
}

var FileLife map[string]int
var db *sql.DB

func init() {
	FileLife = make(map[string]int)
	db, _ = sql.Open("mysql", "root:@/test?charset=utf8")
	db.SetMaxOpenConns(10)
}

func DBConn() *sql.DB {
	return db
}

func UploadMeta(filemeta FileMeta) bool {
	stmt, err := DBConn().Prepare("insert ignore into fileinfo (`file_md5`,`file_name`,`file_size`,`directory`,`full_path`) values (?,?,?,?,?)")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filemeta.FileMD5, filemeta.FileName, filemeta.FileSize, filemeta.Directory, filemeta.FullPath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with fileMD5:%s has been uploaded before", filemeta.FileMD5)
		}
		return true
	}
	return false
}

func GetFileMeta(fullpath string) (*FileMeta, error) {
	filename := filepath.Base(fullpath)
	directory := filepath.Dir(fullpath)
	log.Println(filename)
	log.Println(directory)
	stmt, err := DBConn().Prepare("select file_md5 ,directory,file_name,file_size, full_path from fileinfo where file_name=? and directory=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := FileMeta{}
	err = stmt.QueryRow(filename, directory).Scan(
		&tfile.FileMD5, &tfile.Directory, &tfile.FileName, &tfile.FileSize, &tfile.FullPath)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &tfile, nil
}

func DeleteFileMeta(fullpath string) (*FileMeta, error) {
	filename := filepath.Base(fullpath)
	directory := filepath.Dir(fullpath)
	log.Println(filename)
	log.Println(directory)
	stmt, err := DBConn().Prepare("delete from fileinfo where full_path=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tfile := FileMeta{}
	err = stmt.QueryRow(fullpath).Scan(
		&tfile.FileMD5, &tfile.Directory, &tfile.FileName, &tfile.FileSize, &tfile.FullPath)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return &tfile, nil
}

func GetFileMetaList() ([]FileMeta, error) {
	stmt, err := DBConn().Prepare("select file_md5 ,directory,file_name,file_size, full_path from fileinfo")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	var tfiles []FileMeta
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := FileMeta{}
		err = rows.Scan(&tfile.FileMD5, &tfile.Directory, &tfile.FileName, &tfile.FileSize, &tfile.FullPath)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	fmt.Println(len(tfiles))
	return tfiles, nil
}

func AddFileLife(FileMD5 string) bool {
	stmt, err := DBConn().Prepare("UPDATE filelife SET file_life = file_life + 1 WHERE file_md5 = ?;")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(FileMD5)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			stmt, err := DBConn().Prepare("insert ignore into filelife (`file_md5`, `file_life`) values (?,1);")
			if err != nil {
				fmt.Println("Failed to prepare statement, err:" + err.Error())
				return false
			}
			defer stmt.Close()

			_, err = stmt.Exec(FileMD5)
			if err != nil {
				fmt.Println(err.Error())
				return false
			}
		}
		return true
	}
	return false
}

func MinusFileLife(FileMD5 string) bool {
	stmt, err := DBConn().Prepare("delete from filelife WHERE file_md5 = ? and file_life = 1;")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(FileMD5)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			stmt, err = DBConn().Prepare("UPDATE filelife SET file_life = file_life - 1 WHERE file_md5 = ? and file_life > 1;")
			if err != nil {
				fmt.Println("Failed to prepare statement, err:" + err.Error())
				return false
			}
			defer stmt.Close()

			ret, err := stmt.Exec(FileMD5)
			if err != nil {
				fmt.Println(err.Error())
				return false
			}
			if rf, err := ret.RowsAffected(); nil == err {
				if rf <= 0 {
					fmt.Printf("File with fileMD5:%s has been removed", FileMD5)
				}
				return true
			}
		} else {
			os.Remove("./tmp/" + FileMD5)
		}
		return true
	}
	return false

}

//func UpdateFileDirectory(olddirectory string, newdirectory string) bool {
//	stmt, err := DBConn().Prepare(
//		"update fileinfo set`directory`=? where  `directory`=? limit 1")
//	if err != nil {
//		fmt.Println("预编译sql失败, err:" + err.Error())
//		return false
//	}
//	defer stmt.Close()
//
//	ret, err := stmt.Exec(newdirectory, olddirectory)
//	if err != nil {
//		fmt.Println(err.Error())
//		return false
//	}
//	if rf, err := ret.RowsAffected(); nil == err {
//		if rf <= 0 {
//			fmt.Printf("更新文件location失败, olddirectory:%s", olddirectory)
//		}
//		return true
//	}
//	return false
//}
