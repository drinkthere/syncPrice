package main

import (
	"fmt"
	"os"
	"runtime"
	"syncPrice/config"
	"syncPrice/context"
	"syncPrice/utils"
	"syncPrice/utils/logger"
	"syncPrice/watchdog"
	"time"
)

var globalConfig config.Config
var globalContext context.GlobalContext

func ExitProcess() {
	// 停止挂单
	logger.Info("[Exit] stop and exit.")
	os.Exit(1)
}

func startRiskControl() {
	watchdog.StartRiskControl(&globalConfig, &globalContext)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}

	// 监听退出消息，并调用ExitProcess进行处理
	utils.RegisterExitSignal(ExitProcess)

	// 加载配置文件
	globalConfig = *config.LoadConfig(os.Args[1])

	// 设置使用CPU的核数
	cpuNum := 1
	if globalConfig.CPUNum != nil {
		cpuNum = *globalConfig.CPUNum
	}
	runtime.GOMAXPROCS(cpuNum)

	// 设置日志级别, 并初始化日志
	logger.InitLogger(globalConfig.LogPath, globalConfig.LogLevel)

	// 初始化context
	globalContext.Init(&globalConfig)

	// 开始监听ticker消息
	startTickerMessage()

	// 等等ws数据
	time.Sleep(10 * time.Second)

	// 开启风险控制
	startRiskControl()

	// 开启取消出错的订单上报
	startSendMedianPrice()

	// 阻塞主进程
	for {
		time.Sleep(24 * time.Hour)
	}
}
