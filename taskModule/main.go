package main

import (
	"fmt"
	"github.com/alex023/clock"
	"time"
	"sync"
	"regexp"
	"net/http"
	"io/ioutil"
	"github.com/frank2019/mahonia"
	fm "github.com/cjwddz/fast-model"
	"strings"
	"github.com/fatih/set"
	"./model"
	"./conf"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strconv"
	"github.com/go-redis/redis"
)
var (
	// 上一次采集的任务
	tasksSet *set.SetNonTS
	taskDetail = regexp.MustCompile(`thread-([\d]{7,10}-[\d]{0,3}-[\d]{0,3}).html.+?target=_blank>(.+?)</a></td><td><font color=#FF6600><b>￥([0-9.]{1,10})</b></font></td><td>(.+?)</td>.+?([0-9]{4}-[0-9]{1,2}-[0-9]{1,2} [0-9]{1,2}:[0-9]{1,2}).+?home.php\?mod=space&uid=([\d]+)`)
	url = "https://bbs.125.la/plugin.php?id=e3600%3Atask&mod=show"
	wg sync.WaitGroup
	m fm.DbModel
	finish  = true
)
func connect_redis() *redis.Client  {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "admin123", // no password set
		DB:       0,  // use default DB
	})
	return client
}
func main(){
	// 初始化redis
	client:= connect_redis()
	defer client.Close()
	cache := connect_redis()
	cache.Set("name", "shanhuhai", 0)

	defer notifySimple("taskCatcher","<h2>服务停止！</h2>")
	conf.Init("./task.toml")
	notifySimple("taskCatcher","<h2>服务启动！</h2>")
	fm.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName, conf.App.DBDriver)
	m,_=model.GetTaskModel()
	tasksSet = set.NewNonTS()
	timerTask()
	job, _ := clock.NewClock().AddJobRepeat(time.Second*conf.App.Duration,0, timerTask)
	defer job.Cancel()
	wg.Add(1)
	wg.Wait()
}
var enc = mahonia.NewDecoder("gbk")
func timerTask(){
	if !finish{
		return
	}
	finish = false
	rep,err:=http.Get(url)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	defer rep.Body.Close()
	ds,err:=ioutil.ReadAll(rep.Body)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	allstring:= enc.ConvertString(strings.Replace(strings.Replace(string(ds),`"`,``,-1),"\\","",-1))
	rs:=taskDetail.FindAllStringSubmatch(allstring,-1)
	lenght :=len(rs)
	tasks:=make([]model.Task, lenght)
	newTs:=set.NewNonTS()
	for i:=0;i< lenght;i++{
		tasks[i].Id = rs[i][1]
		tasks[i].Uid = rs[i][6]
		tasks[i].Content = rs[i][2]
		tasks[i].Money = rs[i][3]
		tasks[i].Label = rs[i][4]
		tasks[i].Time = rs[i][5]
		newTs.Add(tasks[i])
	}
	dfs:=set.Difference(newTs,tasksSet)
	fmt.Println(fmt.Sprintf("different count:%d",dfs.Size()))
	tasksSet = newTs
	// 处理新加入的tasks
	notifyTasks:=make([]model.Task,0)
	dfs.Each(func(i interface{}) bool {
		tmp:=i.(model.Task)
		t:=model.Task{
			Id:tmp.Id,
			Uid:tmp.Uid,
			Content:tmp.Content,
			Money:tmp.Money,
			Label:tmp.Label+"源码",
			Time:tmp.Time,
		}
		err:=m.Insert(t)
		if err!=nil{
			fmt.Println(err.Error())
			return true
		}
		if taskFilter(&t){
			notifyTasks=append(notifyTasks,t)
		}else{
			fmt.Println("排除："+t.Label)
		}
		return true
	})
	if len(notifyTasks)>0{
		go notify(notifyTasks)
	}
	finish = true
}
func taskFilter(task *model.Task)bool{
	// 按金额过滤
	if money,err:=strconv.Atoi(task.Money);err==nil && conf.App.MinMoney>0 && conf.App.MoreThanMoney>0{
		// 不足最低金额
		if money<conf.App.MinMoney{
			task.Label=task.Label+",低报酬"
			return false
		}
		// 回报丰厚
		if money>=conf.App.MoreThanMoney{
			task.Label=task.Label+",报酬丰富"
			return true
		}
	}
	// 按关键字过滤
	keywords:=strings.Split(conf.App.KeyWords,",")
	tc:=strings.ToLower(task.Content)
	for _,w:=range keywords{
		if strings.Contains(tc,w){
			task.Label=fmt.Sprintf("%s,%s相关",task.Label,w)
			return true
		}
	}
	task.Label=task.Label+",非感兴趣项目"
	return false
}
func notify(tasks []model.Task){
	if tasks==nil || len(tasks)==0{
		return
	}
	fmt.Println("准备发送邮件通知管理员..")
	html:=""
	for _,t:=range tasks{
		html=fmt.Sprintf(`%s<h2>%s</h1><a href="https://bbs.125.la/home.php?mod=space&uid=%s">用户主页</a><br><a href="https://bbs.125.la/thread-%s.html">内容：%s</a><p>价格：%s</p><p>发布时间：%s</p><hr>`,html,t.Label,t.Uid,t.Id,t.Content,t.Money,t.Time)
	}
	e := email.NewEmail()
	e.From = conf.App.EmailFrom
	e.To = strings.Split(conf.App.EmailTolist,",")
	e.Subject = fmt.Sprintf("taskCatcher[%s]",tasks[0].Label)
	e.HTML = []byte(html)
	err:=e.Send(conf.App.EmailAddress, smtp.PlainAuth("", conf.App.EmailAccount, conf.App.EmailPassword, conf.App.EmailHost))
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(fmt.Sprintf("成功发送通知,管理员将接受到%d条新任务提醒！",len(tasks)))
}

func notifySimple(title,content string){
	e := email.NewEmail()
	e.From = conf.App.EmailFrom
	e.To = strings.Split(conf.App.EmailTolist,",")
	e.Subject = title
	e.HTML = []byte(content)
	err:=e.Send(conf.App.EmailAddress, smtp.PlainAuth("", conf.App.EmailAccount, conf.App.EmailPassword, conf.App.EmailHost))
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
}