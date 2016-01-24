package controllers

import (
	"github.com/astaxie/beego"
	"os"
	"io/ioutil"
	"errors"
	"strings"
	"os/exec"
	"encoding/json"
	"net/http"
	"fmt"
	"bufio"
	"time"
)
type MainController struct {
	beego.Controller
}

type task struct {
	Urls []string `json:"urls"`
}

type Callback struct {
	Url string `json:"url"`
	Email []string `json:"email"`
	AcptNotice bool `json:"acptNotice"`
}
type Task struct {
	Urls []string `json:"urls"`
	Dirs []string `json:"dirs"`
	Callback *Callback `json:"callback"`
}
type Query struct{
	Username string `json:"username"`
	Password string `json:"password"`
	Task *Task `json:"task"`
}

func (this *MainController) Get() {
	this.TplNames = "index.html"
}

func (this  *MainController) Refresh(){
	url := strings.TrimSpace(this.GetString("url"))
    if strings.TrimSpace(url) == ""{
        this.ServeErrJson("URL can not be empty")
        return
    }
	err := reVarnish(url)
	if err != nil {
        WriteLog("error",fmt.Sprintf("Refresh varnish error:%s",err.Error()))
		this.ServeErrJson("refresh varnish error:" + err.Error())
		return
	}
	err = reFile(url)
	if err != nil {
        WriteLog("error",fmt.Sprintf("Refresh file error:$s",err.Error()))
		this.ServeErrJson("refresh file error:" + err.Error())
		return
	}
    WriteLog("info",fmt.Sprintf("refreshed successfully -> %s",url))
	this.ServeOKJson()
}

func reFile(url string) error{
	rawQuery := new(Query)
	rawQuery.Username = "people-93"
	rawQuery.Password = "Cjw%7A9j$m"
	rawQuery.Task = &Task{Urls:[]string{url}}
	b, err := json.Marshal(rawQuery)
	if err != nil {
		return err
	}
	client := &http.Client{}
	request, err := http.NewRequest("POST","https://r.chinacache.com/content/refresh",strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type","application/json")
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func reVarnish(url string) error{
	if !IsExist("cache_server_list.log"){
		return errors.New("Config file is not existent")
	}
	b, err := ioutil.ReadFile("cache_server_list.log")
	if err != nil {
		return err
	}
	servers := strings.Split(string(b),"\n")
	for _, server := range servers{
		s := strings.TrimSpace(server)
		cmd := exec.Command("/bin/sh","-c","curl -X PURGE " + url + " -x " + s + ":80")
		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MainController) ServeErrJson(msg string) {
	this.Data["json"] = map[string]interface{}{
		"msg": msg,
	}
	this.ServeJson()
}

func (this *MainController) ServeOKJson() {
	this.Data["json"] = map[string]interface{}{
		"msg": "",
	}
	this.ServeJson()
}

func IsExist(cfg string) bool{
	_,err := os.Open(cfg)
	return err == nil || os.IsExist(err)
}

func WriteLog(tpe,log string){
    var fi *os.File
    var err error
	var logMsg string
	fi, err = os.OpenFile("refresh.log",os.O_WRONLY|os.O_CREATE|os.O_APPEND,0777)
	if err != nil {
		fmt.Println(err)
	}
	defer fi.Close()
	out := bufio.NewWriter(fi)
	switch strings.ToLower(tpe) {
	case "error":
		logMsg = fmt.Sprintf("[ERROR] %s: %s\r\n",time.Now().Format("2006-01-02 15:04:05"),log)
	case "info":
		logMsg = fmt.Sprintf("[INFO] %s: %s\r\n",time.Now().Format("2006-01-02 15:04:05"),log)
	}
	fi.WriteString(logMsg)
	out.Flush()
    switch strings.ToLower(tpe){
    case "info":
        CaptureLog(logMsg)
    case "error":
        CaptureError(logMsg)
    }
}
