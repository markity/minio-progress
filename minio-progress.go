package progress

import "fmt"

func NewProgress(total int64) *Progress {
	return &Progress{total: total, current: 0, percent: 0}
}

type Progress struct {
	total   int64
	current int64
	percent int
}

// 拼接多字符字符串 mulitSign("=", 3) => "==="
func mulitSign(sign string, num int) string {
	res := ""
	for i := 0; i < num; i++ {
		res += sign
	}

	return res
}

// 打印进度条到标准输出 [==================  ] 94%
func draw(current int64, total int64, percent int) {
	num := percent / 5
	fmt.Printf("\r[%v%v] %v / %v %v%%", mulitSign("=", num), mulitSign(" ", 20-num), current, total, percent)
}

func (progress *Progress) Read(b []byte) (int, error) {
	n := int64(len(b))
	progress.current += n
	percent := int(float64(progress.current) * 100 / float64(progress.total))
	// 只有计算出百分比不同时,才输出并更新progress.percent
	if percent != progress.percent {
		progress.percent = percent
		draw(progress.current, progress.total, progress.percent)
		// 完成时,输出一个换行符
		if progress.current == progress.total {
			fmt.Printf("\n")
		}
	}
	return int(n), nil
}
