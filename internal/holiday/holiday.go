package holiday

import (
	"time"
)

// Holiday 假期配置
type Holiday struct {
	Name      string // 假期名称
	StartDate string // 开始日期 (YYYY-MM-DD)
	EndDate   string // 结束日期 (YYYY-MM-DD)
}

// Festival 节假日配置
type Festival struct {
	Name      string // 节假日名称
	StartDate string // 开始日期 (YYYY-MM-DD)
	EndDate   string // 结束日期 (YYYY-MM-DD)
	Greeting  string // 节假日特色问候
}

// 中国2024-2026年国家法定假日
var Festivals = []Festival{
	// 2024年
	{
		Name:      "元旦",
		StartDate: "2024-01-01",
		EndDate:   "2024-01-01",
		Greeting:  "🎉 新年快乐！祝大家元旦假期开心！",
	},
	{
		Name:      "春节",
		StartDate: "2024-02-10",
		EndDate:   "2024-02-17",
		Greeting:  "🧧 春节快乐！祝大家新春快乐，恭喜发财！",
	},
	{
		Name:      "清明节",
		StartDate: "2024-04-04",
		EndDate:   "2024-04-06",
		Greeting:  "🌿 清明节安康！缅怀先人，珍惜当下。",
	},
	{
		Name:      "劳动节",
		StartDate: "2024-05-01",
		EndDate:   "2024-05-05",
		Greeting:  "💪 劳动节快乐！感谢所有劳动者的付出！",
	},
	{
		Name:      "端午节",
		StartDate: "2024-06-10",
		EndDate:   "2024-06-10",
		Greeting:  "🐉 端午节快乐！祝大家粽子香，生活甜！",
	},
	{
		Name:      "中秋节",
		StartDate: "2024-09-17",
		EndDate:   "2024-09-17",
		Greeting:  "🌕 中秋节快乐！祝大家月圆人圆事圆满！",
	},
	{
		Name:      "国庆节",
		StartDate: "2024-10-01",
		EndDate:   "2024-10-07",
		Greeting:  "🎊 国庆节快乐！祝伟大祖国繁荣昌盛！",
	},
	// 2025年
	{
		Name:      "元旦",
		StartDate: "2025-01-01",
		EndDate:   "2025-01-01",
		Greeting:  "🎉 新年快乐！祝大家元旦假期开心！",
	},
	{
		Name:      "春节",
		StartDate: "2025-01-29",
		EndDate:   "2025-02-06",
		Greeting:  "🧧 春节快乐！祝大家新春快乐，恭喜发财！",
	},
	{
		Name:      "清明节",
		StartDate: "2025-04-04",
		EndDate:   "2025-04-06",
		Greeting:  "🌿 清明节安康！缅怀先人，珍惜当下。",
	},
	{
		Name:      "劳动节",
		StartDate: "2025-05-01",
		EndDate:   "2025-05-05",
		Greeting:  "💪 劳动节快乐！感谢所有劳动者的付出！",
	},
	{
		Name:      "端午节",
		StartDate: "2025-06-02",
		EndDate:   "2025-06-02",
		Greeting:  "🐉 端午节快乐！祝大家粽子香，生活甜！",
	},
	{
		Name:      "中秋节",
		StartDate: "2025-09-07",
		EndDate:   "2025-09-07",
		Greeting:  "🌕 中秋节快乐！祝大家月圆人圆事圆满！",
	},
	{
		Name:      "国庆节",
		StartDate: "2025-10-01",
		EndDate:   "2025-10-07",
		Greeting:  "🎊 国庆节快乐！祝伟大祖国繁荣昌盛！",
	},
	// 2026年
	{
		Name:      "元旦",
		StartDate: "2026-01-01",
		EndDate:   "2026-01-01",
		Greeting:  "🎉 新年快乐！祝大家元旦假期开心！",
	},
	{
		Name:      "春节",
		StartDate: "2026-02-17",
		EndDate:   "2026-02-24",
		Greeting:  "🧧 春节快乐！祝大家新春快乐，恭喜发财！",
	},
	{
		Name:      "清明节",
		StartDate: "2026-04-04",
		EndDate:   "2026-04-06",
		Greeting:  "🌿 清明节安康！缅怀先人，珍惜当下。",
	},
	{
		Name:      "劳动节",
		StartDate: "2026-05-01",
		EndDate:   "2026-05-05",
		Greeting:  "💪 劳动节快乐！感谢所有劳动者的付出！",
	},
	{
		Name:      "端午节",
		StartDate: "2026-06-22",
		EndDate:   "2026-06-22",
		Greeting:  "🐉 端午节快乐！祝大家粽子香，生活甜！",
	},
	{
		Name:      "中秋节",
		StartDate: "2026-09-25",
		EndDate:   "2026-09-25",
		Greeting:  "🌕 中秋节快乐！祝大家月圆人圆事圆满！",
	},
	{
		Name:      "国庆节",
		StartDate: "2026-10-01",
		EndDate:   "2026-10-07",
		Greeting:  "🎊 国庆节快乐！祝伟大祖国繁荣昌盛！",
	},
}

// IsWorkday 检查是否为工作日（周一到周五）
func IsWorkday(t time.Time) bool {
	weekday := t.Weekday()
	// 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	return weekday >= 1 && weekday <= 5
}

// IsHoliday 检查是否在假期期间
func IsHoliday(t time.Time, holidays []Holiday) bool {
	dateStr := t.Format("2006-01-02")
	for _, h := range holidays {
		startDate, _ := time.Parse("2006-01-02", h.StartDate)
		endDate, _ := time.Parse("2006-01-02", h.EndDate)
		if dateStr >= h.StartDate && dateStr <= h.EndDate {
			return true
		}
		// 确保时间比较正确
		if t.After(startDate) && t.Before(endDate.AddDate(0, 0, 1)) {
			return true
		}
	}
	return false
}

// IsFestival 检查是否为节假日，返回节假日信息及特色问候
func IsFestival(t time.Time) (bool, *Festival) {
	dateStr := t.Format("2006-01-02")
	for i, festival := range Festivals {
		if dateStr >= festival.StartDate && dateStr <= festival.EndDate {
			return true, &Festivals[i]
		}
	}
	return false, nil
}

// ShouldSendReminder 检查是否应该发送提醒
// 返回: (是否发送, 是否为节假日, 节假日信息)
func ShouldSendReminder(holidays []Holiday) (bool, bool, *Festival) {
	now := time.Now()
	
	// 检查是否在假期期间
	if IsHoliday(now, holidays) {
		return false, false, nil
	}
	
	// 检查是否为节假日
	isFestival, festival := IsFestival(now)
	if isFestival {
		// 节假日: 不发送普通�报告，但发送节假日特色问候
		return false, true, festival
	}
	
	// 检查是否为工作日
	if IsWorkday(now) {
		return true, false, nil
	}
	
	// 周末不提醒
	return false, false, nil
}

// ShouldSendOffWorkReminder 检查是否应该发送下班提醒
// 返回: (是否发送, 是否为节假日, 节假日信息)
func ShouldSendOffWorkReminder(holidays []Holiday) (bool, bool, *Festival) {
	now := time.Now()
	
	// 检查是否在假期期间
	if IsHoliday(now, holidays) {
		return false, false, nil
	}
	
	// 检查是否为节假日（节假日不发送下班提醒）
	isFestival, _ := IsFestival(now)
	if isFestival {
		return false, false, nil
	}
	
	// 检查是否为工作日
	if IsWorkday(now) {
		return true, false, nil
	}
	
	// 周末不提醒
	return false, false, nil
}

// parseDate 解析日期字符串
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
