package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const statusUrl = "http://192.168.9.8/include/auth_action.php?k="
const loginUrl = "http://192.168.9.8/include/auth_action.php"
var userName = "username"
var passWord = "password"
var pidFile = "connect.pid"
var settingFile = "settings.ini"
var checkInterval time.Duration = 600

func main() {
	fmt.Println("Usage:\n    command {setting_file_path} {pid_file_path}")
	// 第一个参数是设置文件的路径
	if len(os.Args) > 1 {
		settingFile = os.Args[1]
	}
	// 第二个参数是 pid 文件的路径
	if len(os.Args) > 2 {
		pidFile = os.Args[2]
	}
	// 第三个是日志文件的路径
	// TODO
	res := loadConfig(settingFile)
	if !res {
		return
	}
	if !writePid(pidFile) {
		log.Println("Failed to write pid file: " + pidFile)
	}
	for {
		isLogin := checkStatus()
		if !isLogin {
			loginStatus := login()
			if !loginStatus {
				log.Println("Login failed, sleep for 1 minute")
				time.Sleep(time.Minute * 1)
			} else {
				log.Printf("Login successfully, sleep for %d seconds\n", checkInterval)
				time.Sleep(time.Second * checkInterval)
			}
			continue
		}
		log.Printf("Network has been online, sleep for %d seconds\n", checkInterval)
		time.Sleep(time.Second * checkInterval)
	}
}

func loadConfig(filename string) bool {
	config, err := ini.Load(filename)
	if err != nil {
		log.Println("Load settings file error: " + err.Error())
		return false
	}
	userName = config.Section("user").Key("username").String()
	passWord = config.Section("user").Key("password").String()
	checkTime, err := config.Section("user").Key("check_interval").Int()
	if err != nil {
		log.Println("Error occurred while parsing check interval")
	} else {
		checkInterval = time.Duration(checkTime)
	}
	log.Printf("Load settings from %s successfully.\n", filename)
	return true
}

func checkStatus() bool {
	keys := rand.Intn(10000)
	urls := statusUrl + strconv.Itoa(keys)
	payload := url.Values{
		"action": {"get_online_info"},
		"keys": {strconv.Itoa(keys)},
	}
	res, err := http.PostForm(urls, payload)
	if err != nil {
		log.Println("Check connection status failed")
		return false
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	text := string(body)
	if text == "not_online" {
		return false
	} else {
		data := strings.Split(text, ",")
		dataBytes, _ := strconv.Atoi(data[0])
		linkTime, _ := strconv.Atoi(data[1])
		log.Println("Account Status: Data Usage: " + strconv.Itoa(dataBytes >> 20) + " MB, Connection Time: " +
			(time.Duration(linkTime) * time.Second).String() + ", IP Address: " + data[5])
		return true
	}
}

func login() bool {
	payload := url.Values{
		"action": {"login"},
		"username": {userName},
		"password": {passWord},
		"ac_id": {"1"},
		"user_mac": {""},
		"user_ip": {""},
		"nas_ip": {""},
		"save_me": {"0"},
		"domain": {"@uestc"},
		"ajax": {"1"},
	}
	res, err := http.PostForm(loginUrl, payload)
	if err != nil {
		log.Println("Login connect failed.")
		return false
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer res.Body.Close()
	text := string(body)
	if strings.Contains(text, "login_ok") {
		log.Println("Login successfully.")
		return true
	}
	log.Println("Login Failed.")
	return false
}

func writePid(filename string) bool {
	pid := strconv.Itoa(os.Getpid())
	file, err := os.Create(filename)
	if err != nil {
		return false
	}
	_, _ = file.WriteString(pid)
	defer file.Close()
	log.Printf("Write to pid file(%s) successfully, pid: %s", filename, pid)
	return true
}
