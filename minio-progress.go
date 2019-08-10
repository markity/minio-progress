package progress

import (
	"fmt"
	"io"

	"github.com/minio/minio-go"
)

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

func NewUploadProgress(total int64) *UploadProgress {
	return &UploadProgress{total: total, current: 0, percent: 0}
}

type UploadProgress struct {
	total   int64
	current int64
	percent int
}

func (progress *UploadProgress) Read(b []byte) (int, error) {
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

func CopyWithProgress(dst io.Writer, object *minio.Object) (int64, error) {
	var totalRead int64 = 0

	// 每次读32 * 1024个字节
	size := 32 * 1024
	data := make([]byte, size)

	objInfo, err := object.Stat()
	if err != nil {
		return totalRead, err
	}
	objSize := objInfo.Size

	for {
		// 读取
		nRead, errRead := object.Read(data)
		if errRead != nil && errRead != io.EOF {
			// 未知异常,读取失败
			return totalRead, errRead
		}
		data = data[:nRead]

		// 写入
		nWrite, errWrite := dst.Write(data)
		totalRead += int64(nWrite)
		if errWrite != nil {
			// 写入失败
			return totalRead, errWrite
		}

		object.Stat()
		// 打印进度条
		percent := int(float64(totalRead) * 100 / float64(objSize))
		draw(totalRead, objSize, percent)

		if errRead == io.EOF {
			fmt.Printf("\n")
			return totalRead, nil
		}
	}
}
