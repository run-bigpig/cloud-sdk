package utils

import "time"

func DateTimeToTimeStamp(datatime string) int64 {
	t, err := time.Parse("2006-01-02 15:04:05", datatime)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// FormatTimeWithTimezone 时间戳格式化为时间
func FormatTimeWithTimezone(timestamp int64, timezone string) string {
	location, _ := time.LoadLocation(timezone)
	return time.Unix(timestamp, 0).In(location).Format("2006-01-02 15:04:05")
}

// DateTimeToTimaStampWithTimezone 时间为时间戳
func DateTimeToTimeStampWithTimezone(dateTime string, timezone string) int64 {
	location, _ := time.LoadLocation(timezone)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", dateTime, location)
	return t.Unix()
}

// CalcTimeStampsWithInterval 计算两个时间戳间指定粒度的所有的时间戳
func CalcTimeStampsWithInterval(startTime int64, endTime int64, interval int64) []int64 {
	result := make([]int64, 0)
	for i := startTime; i <= endTime; i += interval {
		result = append(result, i)
	}
	return result
}
