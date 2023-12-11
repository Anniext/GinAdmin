package core

import (
	"GinAdmin/core/internal"
	"GinAdmin/global"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 优先级: 命令行 > 环境变量 > 默认值

func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		// 绑定-c的值到config上
		flag.StringVar(&config, "c", "", "choose config file")
		flag.Parse()
		if config == "" {
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" {
				switch gin.Mode() {
				case gin.DebugMode:
					config = internal.ConfigDebugFile
					fmt.Printf("您正在使用gin模式的%s环境, config路径为%s\n", gin.EnvGinMode, internal.ConfigDefaultFile)
				case gin.ReleaseMode:
					config = internal.ConfigReleaseFile
					fmt.Printf("您正在使用gin模式的%s环境, config路径为%s\n", gin.EnvGinMode, internal.ConfigReleaseFile)
				case gin.TestMode:
					config = internal.ConfigTestFile
					fmt.Printf("您正在使用gin模式的%s环境, config路径为%s\n", gin.EnvGinMode, internal.ConfigTestFile)
				}
			} else {
				config = configEnv
				fmt.Printf("您正在使用gin模式的%s环境, config路径为%s\n", gin.EnvGinMode, config)
			}
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else {
		config = path[0]
		fmt.Printf("您正在使用gin模式的%s环境, config路径为%s\n", gin.EnvGinMode, config)
	}
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed: ", e.Name)
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			panic(err)
		}
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		panic(err)
	}
  global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return v
}
