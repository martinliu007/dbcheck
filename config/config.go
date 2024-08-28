package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Config 结构体，表示配置文件中变量，与config.yaml文件里的配置对应
type Config struct {
	DBHost       string `yaml:"DBHost"` // 用来存放校验结果的数据库Host，其中`yaml:"DBHost"`，表示将yaml文件中的DBHost键对应的值复制给这个结构体中的DBHost
	DBPort       string    `yaml:"DBPort"`  // 用来存放校验结果的数据库端口
	DBUsername   string `yaml:"DBUsername"`  // 用来存放校验结果的数据库用户名
	DBPassword   string `yaml:"DBPassword"` // 用来存放校验结果的数据库密码
	DBName       string `yaml:"DBName"`  // 用来存放校验结果的数据库库名
	CheckDBList  []string `yaml:"CheckDBList"`  //需要校验的数据库实例，这里注意，是切片类型，也就是可以有多个MySQL实例
	CheckDBPort	 string `yaml:"CheckDBPort"`  // 需要校验的数据库端口
	CheckDBUsername	 string `yaml:"CheckDBUsername"`  // 需要校验的数据库用户名
	CheckDBPassword	 string `yaml:"CheckDBPassword"`  // 需要校验的数据库用户密码

}

// LoadConfig 初始化配置
func LoadConfig(configPath string) (*Config, error) {
	// 读取配置文件内容
	configFile, err := os.ReadFile(configPath)

	// 如果读取文件过程中发生错误，返回nil和错误信息
	if err != nil {
		return nil, err
	}
	// 创建一个空的Config结构体实例
	config := &Config{}

	// 将配置文件内容解析为Config结构体对象
	err = yaml.Unmarshal(configFile, config)

	// 如果解析配置文件过程中发生错误，返回nil和错误信息
	if err != nil {
		return nil, err
	}

	// 返回解析后的Config结构体对象和nil
	return config, nil
}
