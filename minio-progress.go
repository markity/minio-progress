package progress

import (
	"fmt"
	"io"

	minio "github.com/minio/minio-go"
)

// 拼接多字符字符串 mulitSign("=", 3) => "==="
func mulitSign(sign string, num int) string {
	res := ""
	for i := 0; i < num; i++ {
		res += sign
	}

	return res
}

// 打印进度条到标准输出 [==================  ] 30 / 35 85%
func draw(current int64, total int64, percent int) {
	num := percent / 5
	fmt.Printf("\r[%v%v] %v / %v %v%%", mulitSign("=", num), mulitSign(" ", 20-num), current, total, percent)
}

// NewUploadProgress: Create a *UploadProgress object
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

// CopyWithProgress: Copy the stream and print the progress bar
func CopyWithProgress(dst io.Writer, object *minio.Object) (int64, error) {
	objInfo, err := object.Stat()
	if err != nil {
		return 0, err
	}
	objSize := objInfo.Size

	var written int64 = 0
	buf := make([]byte, 32*1024)
	for {
		// 读取
		nr, er := object.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			percent := int(float64(written) * 100 / float64(objSize))
			draw(written, objSize, percent)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		fmt.Printf("\n")
	}
	return written, err
}
