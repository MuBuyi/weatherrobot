package main

import (
	"fmt"
	"strings"
	"time"
	"wechatrobot/internal/config"
	"wechatrobot/internal/holiday"
)

func main() {
	// 加载配置
	config.Load()

	// 获取当前时间
	now := time.Now()

	sep := strings.Repeat("=", 50)
	fmt.Println(sep)
	fmt.Println("工作日/假期判断测试")
	fmt.Println(sep)

	fmt.Printf("当前时间: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("星期: %v\n", now.Weekday())

	// 检查是否为工作日
	isWorkday := holiday.IsWorkday(now)
	fmt.Printf("是否为工作日: %v\n\n", isWorkday)

	// 检查是否为节假日
	isFestival, festival := holiday.IsFestival(now)
	fmt.Printf("是否为节假日: %v\n", isFestival)
	if isFestival && festival != nil {
		fmt.Printf("节假日名称: %s\n", festival.Name)
		fmt.Printf("节假日问候: %s\n\n", festival.Greeting)
	}

	// 检查是否在假期期间
	isHoliday := holiday.IsHoliday(now, config.Cfg.Holidays)
	fmt.Printf("是否在假期期间: %v\n", isHoliday)
	if isHoliday {
		for _, h := range config.Cfg.Holidays {
			fmt.Printf("假期名称: %s (%s 至 %s)\n", h.Name, h.StartDate, h.EndDate)
		}
		fmt.Println()
	}

	// 检查是否应该发送天气报告
	shouldSendReport, isFestival2, festival2 := holiday.ShouldSendReminder(config.Cfg.Holidays)
	fmt.Println(sep)
	fmt.Println("天气报告判断结果:")
	fmt.Printf("是否应该发送天气报告: %v\n", shouldSendReport)
	fmt.Printf("是否为节假日: %v\n", isFestival2)
	if festival2 != nil {
		fmt.Printf("节假日信息: %s - %s\n", festival2.Name, festival2.Greeting)
	}
	fmt.Println()

	// 检查是否应该发送下班提醒
	shouldSendOffWork, _, _ := holiday.ShouldSendOffWorkReminder(config.Cfg.Holidays)
	fmt.Println(sep)
	fmt.Println("下班提醒判断结果:")
	fmt.Printf("是否应该发送下班提醒: %v\n", shouldSendOffWork)
	fmt.Println()

	// 显示配置的假期
	fmt.Println(sep)
	fmt.Println("已配置的假期列表:")
	if len(config.Cfg.Holidays) == 0 {
		fmt.Println("无自定义假期配置")
	} else {
		for _, h := range config.Cfg.Holidays {
			fmt.Printf("- %s: %s 至 %s\n", h.Name, h.StartDate, h.EndDate)
		}
	}
	fmt.Println(sep)
}
