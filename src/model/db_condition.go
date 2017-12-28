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
}

func (cond *DbCondition)And(compare string,key string) *DbCondition{



	// todo
	if cond.condStr==""{

	}
	 return cond
}



/**
从请求中获取condition表达式和参数值
t_key格式为类型首写和列名，如int类型id则为i_id,再如：s_name,b_valid
因为获取的是检索条件，所以一般只要int，bool和string就好了.
cdc为对比条件： > = <
 */
func GetCondition(r *http.Request,cdc []string, t_key []string)(cond string,args []interface{}){
	c:=1
	conds:= make([]string,0)
	args =make([]interface{},0)
	for i,k:=range t_key {
		if len(k)<=2 || k[1]!='_'{
			// todo 写到系统日志
			fmt.Println("是否错误调用了GetCondition？t_key格式为类型首写和列名，如int类型id则为i_id,再如：s_name,b_valid")
			continue
		}
		t:=k[2:]
		value :=r.PostFormValue(t)
		if value ==""{
			continue
		}
		switch k[0] {
		case 'b':
			if strings.ToLower(value)=="true"{
				args =append(args,true)
			}else{
				args =append(args,false)
			}
			conds=append(conds, fmt.Sprintf("%s %s $%d",t,cdc[i],c))
			c++
		case 'i':
			i,err:=strconv.Atoi(value)
			if err!=nil{
				// todo 写到系统日志
				fmt.Println(fmt.Sprintf("类型转化出错！key=%s,value=%s,err=%s",k,value,err.Error()))
				continue
			}
			args =append(args,i)
		default:
			args =append(args,value)
		}
	}
	cond=strings.Join(conds,"And")
	return
}