package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-vgo/robotgo"
	"github.com/spf13/viper"
)

type config struct {
	interval  int
	actions   []string
	moveDis   int
	typeStr   string
	durations []string
}

var conf = config{}

func init() {
	viper.SetConfigName("conf") //把json文件换成yaml文件，只需要配置文件名 (不带后缀)即可
	viper.AddConfigPath(".")    //添加配置文件所在的路径
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file error: %s\n", err)
		os.Exit(1)
	}

	viper.WatchConfig() //监听配置变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置发生变更：", e.Name)
	})

	conf.interval = viper.GetInt("interval")
	conf.actions = strings.Split(viper.GetString("actions"), "|")
	conf.moveDis = viper.GetInt("move-distance")
	conf.typeStr = viper.GetString("typeStr")
	//conf.durations = strings.Split(viper.GetString("durations"), "|")

	fmt.Println(conf)
}

func main() {

	f, err := os.OpenFile("virus.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * time.Duration(conf.interval))
	go func() {
		for range ticker.C {
			for _, action := range conf.actions {
				if action == "move" {
					move()
				}
				if action == "typeStr" {
					typeStr()
				}
				if action == "click" {
					click()
				}
			}
		}
		ch <- 1
	}()
	<-ch
}

func move() {
	robotgo.MoveMouse(conf.moveDis, conf.moveDis)
}

func typeStr() {
	robotgo.TypeStr(conf.typeStr)
}

func click() {
	robotgo.MouseClick("left", true)
}
