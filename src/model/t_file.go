package model

import (
	"fmt"
	"../pdb"
	"time"
)

/*
DROP TABLE IF EXISTS t_file;
CREATE TABLE t_file (
"id" serial NOT NULL,
"file_key" varchar(16) UNIQUE COLLATE "default",
"file_type" varchar(16) COLLATE "default",
"file_name" varchar(128) COLLATE "default",
"owner" int4,
"path" varchar(255) COLLATE "default",
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */

type T_File struct {
	ID        int64     `json:"id"`
	FileKey   string    `json:"file_key"`
	FileType  string    `json:"file_type"`
	FileName  string    `json:"file_name"`
	Owner     int       `json:"owner"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

const(
	FILE_TYPE_PIC = "picture"
	FILE_TYPE_FILE = "file"
)
func FileTableName() string {
	return "t_file"
}

func (f *T_File) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(file_key,file_name,file_type,owner,path,created_time) "+
			  "VALUES($1,$2,$3,$4,$5,$6)", FileTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	f.CreatedAt = time.Now()
	_, err = stmt.Exec(f.FileKey,f.FileName,f.FileType,f.Owner,f.Path,f.CreatedAt)
	return
}

func FindFiles(condition, limit, order string) (result []T_File,err error) {
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,file_key,file_name,file_type,owner,path,created_time FROM %s %s %s %s", FileTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := T_File{}
		err = rows.Scan(&tmp.ID,&tmp.FileKey,&tmp.FileName,&tmp.FileType,&tmp.Owner,&tmp.Path,&tmp.CreatedAt)
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
