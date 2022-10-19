package common

import "time"

var (
	cst *time.Location
)

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"

// RFC3339ToCSTLayout 将time.RFC3339格式的时间字符串转成China标准时间字符串
func RFC3339ToCSTLayout(value string) (string, error) {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}
	return ts.In(time.Local).Format(CSTLayout), nil
}
