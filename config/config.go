package config

import (
	"github.com/Unknwon/goconfig"
	"log"
	"path/filepath"
	"os"
	"fmt"
)

func GetValue(section,key string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	} else {
		file := filepath.Join(dir, "config.ini")
		cfg, err := goconfig.LoadConfigFile(file)
		if err != nil {
			log.Fatal(err)
		}
		value, err := cfg.GetValue(section, key)
		if err != nil {
			log.Fatalf("无法获取配置%s: %s", key, err)
		}
		return value
	}
	return ""
}
