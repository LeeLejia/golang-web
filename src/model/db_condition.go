package model

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
)

type DbCondition struct{
	condStr string
	condCount int
	order string
	limit_len int
	limit_pos int
	args []interface{}
}

/**
添加AND条件,compare: > < = >= <= != like
 */
func (cond *DbCondition)And(r *http.Request,compare string, t_key string) *DbCondition{
	return cond.andOr(r,compare,t_key,"AND")
}
/**
添加AND条件,compare: > < = >= <= != like
 */
func (cond *DbCondition)And2(compare string, key string,value interface{}) *DbCondition{
	if cond.condStr==""{
		// 初始化参数
		cond.args = make([]interface{},0)
		cond.condCount=1
	}
	cond.args =append(cond.args,value)
	cond.condStr=fmt.Sprintf("%s AND %s %s $%d",cond.condStr,key,compare,cond.condCount)
	cond.condCount++
	if len(cond.condStr)>4 && cond.condStr[0:4]==" AND"{
		cond.condStr=cond.condStr[5:]
	}
	return cond
}
/**
添加AND条件,compare: > < = >= <= != like
 */
func (cond *DbCondition)Or2(compare string, key string,value interface{}) *DbCondition{
	if cond.condStr==""{
		// 初始化参数
		cond.args = make([]interface{},0)
		cond.condCount=1
	}
	cond.args =append(cond.args,value)
	cond.condStr=fmt.Sprintf("%s AND %s %s $%d",cond.condStr,key,compare,cond.condCount)
	cond.condCount++
	if len(cond.condStr)>4 && cond.condStr[0:4]==" OR "{
		cond.condStr=cond.condStr[5:]
	}
	return cond
}
/**
添加OR条件
 */
func (cond *DbCondition)Or(r *http.Request,compare string, t_key string) *DbCondition{
	return cond.andOr(r,compare,t_key,"OR")
}
/**
设置LIMIT语句
 */
func (cond *DbCondition)Limit(r *http.Request,startKey string, lenKey string)*DbCondition{
	pos:=-1
	len :=-1
	t,err:=strconv.Atoi(r.PostForm.Get(startKey))
	if err==nil{
		pos=t
	}
	t,err=strconv.Atoi(r.PostForm.Get(lenKey))
	if err==nil{
		len =t
	}
	cond.limit_len = len
	cond.limit_pos = pos
	return cond
}
/**
设置order语句
 */
func (cond *DbCondition)Order(order string)*DbCondition {
	cond.order = order
	return cond
}

/**
获取WHERE表达式
 */
func (cond *DbCondition)GetWhere()string{
	if cond.limit_pos>0 && cond.limit_len>0{
		return fmt.Sprintf("WHERE %s %s limit $%d offset $%d",cond.condStr,cond.order,cond.condCount,cond.condCount+1)
	}else if cond.limit_pos>0{
		return fmt.Sprintf("WHERE %s %s offset $%d",cond.condStr,cond.order,cond.condCount+1)
	}else if cond.limit_len>0{
		return fmt.Sprintf("WHERE %s %s limit $%d",cond.condStr,cond.order,cond.condCount+1)
	}
	if cond.condStr ==""{
		return fmt.Sprintf("%s",cond.order)
	}
	return fmt.Sprintf("WHERE %s %s",cond.condStr,cond.order)
}
/**
获取参数
 */
func (cond *DbCondition)GetParams() []interface{}{
	rs:=cond.args
	if cond.limit_len>0{
		rs=append(rs,cond.limit_len)
	}
	if cond.limit_pos>0{
		rs=append(rs,cond.limit_pos)
	}
	return rs
}

func (cond *DbCondition)andOr(r *http.Request,compare string, t_key string,ao string) *DbCondition{
	if cond.condStr==""{
		// 初始化参数
		cond.args = make([]interface{},0)
		cond.condCount=1
	}
	if len(t_key)<=2 || t_key[1]!='_'{
		// todo 写到系统日志
		fmt.Println("是否错误调用了GetCondition？t_key格式为类型首写和列名，如int类型id则为i_id,再如：s_name,b_valid")
		return cond
	}
	t:= t_key[2:]
	value :=r.PostFormValue(t)
	if value ==""{
		return cond
	}
	switch t_key[0] {
	case 'b':
		if strings.ToLower(value)=="true"{
			cond.args =append(cond.args,true)
		}else{
			cond.args =append(cond.args,false)
		}
		cond.condStr=fmt.Sprintf("%s %s %s %s $%d",cond.condStr,ao,t,compare,cond.condCount)
		cond.condCount++
	case 'i':
		i,err:=strconv.Atoi(value)
		if err!=nil{
			// todo 写到系统日志
			fmt.Println(fmt.Sprintf("类型转化出错！key=%s,value=%s,err=%s", t_key,value,err.Error()))
			return cond
		}
		cond.args =append(cond.args,i)
		cond.condStr=fmt.Sprintf("%s %s %s %s $%d",cond.condStr,ao,t,compare,cond.condCount)
		cond.condCount++
	default:
		if strings.ToLower(compare)=="like"{
			cond.args =append(cond.args,"%"+value+"%")
			cond.condStr=fmt.Sprintf("%s %s %s like $%d",cond.condStr,ao,t,cond.condCount)
		}else{
			cond.args =append(cond.args,value)
			cond.condStr=fmt.Sprintf("%s %s %s %s $%d",cond.condStr,ao,t,compare,cond.condCount)
		}
		cond.condCount++
	}
	if len(cond.condStr)>4 && (cond.condStr[0:4]==" AND" || cond.condStr[0:4]==" OR "){
		cond.condStr=cond.condStr[5:]
	}
	return cond
}