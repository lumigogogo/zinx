package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lumigogogo/zinx/zlog"
)

// Config ..
type Config struct {
	// TCPServer ziface.IServer
	Host    string
	TCPPort int
	Name    string

	MaxPacketSize  uint32
	MaxConnNum     uint32
	WorkPoolSize   int
	MaxWorkTaskNum int
	MaxMsgChanNum  uint32

	ConfFilePath string

	LogDir        string
	LogFile       string
	LogDebugClose bool
}

// PathExists ..
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//Reload ..
func (c *Config) Reload() {

	if confFileExists, _ := PathExists(c.ConfFilePath); confFileExists != true {
		fmt.Println("Config File ", c.ConfFilePath, " is not exist!!")
		return
	}

	data, err := ioutil.ReadFile(c.ConfFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}

	if c.LogFile != "" {
		fmt.Println("hahahaha")
		zlog.SetLogFile(c.LogDir, c.LogFile)
	}
	if c.LogDebugClose == true {
		fmt.Println("hehehehhehe")
		zlog.CloseDebug()
	}
}

// GlobalConf ..
var GlobalConf *Config

func init() {
	GlobalConf = &Config{
		// default params
		Host:           "0.0.0.0",
		TCPPort:        9527,
		Name:           "zinx-server",
		MaxPacketSize:  512,
		MaxConnNum:     24000,
		WorkPoolSize:   10,
		MaxWorkTaskNum: 5,
		MaxMsgChanNum:  20,
		ConfFilePath:   "./conf/zinx.json",
	}

	GlobalConf.Reload()
}
