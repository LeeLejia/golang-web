package model

import (
	"time"
	"fmt"
	"../pdb"
)
/*
DROP TABLE IF EXISTS t_vlog;
CREATE TABLE t_vlog (
"id" serial NOT NULL,
"tag" varchar(255) COLLATE "default",
"code" varchar(16) COLLATE "default",
"machine" varchar(128) COLLATE "default",
"content" varchar(255) COLLATE "default",
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
const (
	VERIFY_LOG_SUCCESS         = "验证成功"

	VERIFY_LOG_NOEXIST         = "邀请码不存在"
	VERIFY_LOG_INVALID         = "邀请码已经禁用"
	VERIFY_LOG_NOFIND_APP      = "找不到指定app"
	VERIFY_LOG_APP_NOVALID     = "该APP已失效"
	VERIFY_LOG_EMPTY           = "邀请码为空"
	VERIFY_LOG_MACHINE_NOVALID = "机器码不可用"
	VERIFY_LOG_TOOMUCH_MACHINE = "设备数超限制"
	VERIFY_LOG_BEFORE_TIME     = "尚未生效"
	VERIFY_LOG_AFTER_TIME      = "已经过期"

	VERIFY_LOG_PROTO_NOVALID = "协议码不可用"
)

type T_vlog struct {
	ID        int64     `json:"id"`
	Tag       string    `json:"tag"`
	Code      string    `json:"code"`
	Machine   string    `json:"machine"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func VLogTableName() string {
	return "t_vlog"
}

func (m *T_vlog) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(tag,code,machine,content,created_at) "+
		"VALUES($1,$2,$3,$4,$5)", VLogTableName()))
	if err != nil {
		return
	}
	m.CreatedAt = time.Now()
	_, err = stmt.Exec(m.Tag, m.Code, m.Machine, m.Content, m.CreatedAt)
	return
}

func FindVLogs(condition, limit, order string) ([]T_vlog, error) {
	result := []T_vlog{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,tag,code,machine,content,created_at FROM %s %s %s %s", VLogTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := T_vlog{}
		err = rows.Scan(&tmp.ID, &tmp.Tag, &tmp.Code, &tmp.Machine, &tmp.Content, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func CountVLogs(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", VLogTableName(), condition)).Scan(&count)
	return
}
