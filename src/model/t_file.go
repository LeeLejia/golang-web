package model

import (
	"fmt"
	"../pdb"
	"time"
)

type T_File struct {
	ID          int64       `json:"id"`
	FileKey     string      `json:"file_key"`
	FileName	string		`json:"file_name"`
	Owner       int         `json:"owner"`
	Path        string      `json:"path"`
	CreatedTime time.Time   `json:"created_time"`
}
func FileTableName() string {
	return "t_file"
}

func (f *T_File) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(file_key,file_name,owner,path,created_time) "+
			  "VALUES($1,$2,$3,$4,$5)", FileTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	f.CreatedTime = time.Now()
	_, err = stmt.Exec(f.FileKey,f.FileName,f.Owner,f.Path,f.CreatedTime)
	return
}

func FindFiles(condition, limit, order string) (result []T_File,err error) {
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,file_key,file_name,owner,path,created_time FROM %s %s %s %s", FileTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := T_File{}
		err = rows.Scan(&tmp.ID,&tmp.FileKey,&tmp.FileName,&tmp.Owner,&tmp.Path,&tmp.CreatedTime)
		if err==nil {
			result = append(result, tmp)
		}
	}
	return result, err
}

func CountFiles(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", FileTableName(), condition)).Scan(&count)
	return
}
