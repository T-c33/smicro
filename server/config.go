package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/T-c33/smicro/util"
	"gopkg.in/yaml.v2"
)

var (
	smicroConf = &SmicroConf{
		Port: 8080,
		Prometheus: PrometheusConf{
			SwitchOn: true,
			Port:     8081,
		},
		ServiceName: "smicro_server",
		Regiser: RegisterConf{
			SwitchOn:     false,
			RegisterPath: "/smicro/service/",
			Timeout:      time.Second,
			HeartBeat:    10,
			RegisterName: "etcd",
			RegisterAddr: "127.0.0.1:2379",
		},
		Log: LogConf{
			Level:      "debug",
			Dir:        "./logs/",
			ChanSize:   10000,
			ConsoleLog: true,
		},
		Limit: LimitConf{
			SwitchOn: true,
			QPSLimit: 50000,
		},
	}
)

type SmicroConf struct {
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	ServiceName string         `yaml:"service_name"`
	Regiser     RegisterConf   `yaml:"register"`
	Log         LogConf        `yaml:"log"`
	Limit       LimitConf      `yaml:"limit"`
	Trace       TraceConf      `yaml:"trace"`

	//内部的配置项
	ConfigDir  string `yaml:"-"`
	RootDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

type TraceConf struct {
	SwitchOn   bool    `yaml:"switch_on"`
	ReportAddr string  `yaml:"report_addr"`
	SampleType string  `yaml:"sample_type"`
	SampleRate float64 `yaml:"sample_rate"`
}

type LimitConf struct {
	QPSLimit int  `yaml:"qps"`
	SwitchOn bool `yaml:"switch_on"`
}

type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

type RegisterConf struct {
	SwitchOn     bool          `yaml:"switch_on"`
	RegisterPath string        `yaml:"register_path"`
	Timeout      time.Duration `yaml:"timeout"`
	HeartBeat    int64         `yaml:"heart_beat"`
	RegisterName string        `yaml:"register_name"`
	RegisterAddr string        `yaml:"register_addr"`
}

type LogConf struct {
	Level      string `yaml:"level"`
	Dir        string `yaml:"path"`
	ChanSize   int    `yaml:"chan_size"`
	ConsoleLog bool   `yaml:"console_log"`
}

func initDir(serviceName string) (err error) {

	exeFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}

	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid exe path:%v", exeFilePath)
		return
	}
	smicroConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIdx]), "..")
	smicroConf.ConfigDir = path.Join(smicroConf.RootDir, "./conf/", util.GetEnv())
	smicroConf.ConfigFile = path.Join(smicroConf.ConfigDir, fmt.Sprintf("%s.yaml", serviceName))
	return
}

func InitConfig(serviceName string) (err error) {

	err = initDir(serviceName)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(smicroConf.ConfigFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &smicroConf)
	if err != nil {
		return
	}

	fmt.Printf("init smicro conf succ, conf:%#v\n", smicroConf)
	return
}

func GetConfigDir() string {
	return smicroConf.ConfigDir
}

func GetRootDir() string {
	return smicroConf.RootDir
}

func GetServerPort() int {
	return smicroConf.Port
}

func GetConf() *SmicroConf {
	return smicroConf
}
