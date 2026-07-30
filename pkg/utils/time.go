package utils

import "time"

// Now 返回当前时间的标准格式字符串 "2006-01-02 15:04:05"
func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// FormatTime 将 time.Time 格式化为标准字符串
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseTime 将标准字符串解析为 time.Time
func ParseTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// GetTodayStart 获取今天 00:00:00
func GetTodayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
