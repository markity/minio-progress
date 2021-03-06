# 为minio上传下载提供便捷的进度条接口

使用minio-progress快速为minio上传和下载添加进度条

```
[===========         ] 856293376 / 1529094438 56%
```

## 从github安装

```
$ go get github.com/markity/minio-progress
```

## 快速开始

### 上传使用进度条

```
// 创建上传进度条对象
progressBar := progress.NewUploadProgress(fileSize)

// 然后将进度条对象包含在minio.PutObjectOptions中的Progress配置即可
n, err = client.PutObject(bucketName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progressBar})
if err != nil {
    fmt.Printf("上传文件失败:%v\n", err)
}
```

#### 完整案例

main.go

```
package main

import (
    "log"
    "os"

    // 导入进度条包
    progress "github.com/markity/minio-progress"

    // 导入minio SDK
    "github.com/minio/minio-go"
)

// minio客户端基本信息
var endpoint = "127.0.0.1:9000"
var accessKeyID = "Your-AccessKeyID"
var secretAccessKey = "Your-SecretAccessKey"
var secure = false

// 数据桶信息
var bucketName = "music"
var location = "us-east-1"

// 上传文件信息
var objectName = "big.rar"
var filePath = "./big.rar"

func main() {
    client, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
    if err != nil {
        log.Fatalf("创建客户端失败:%v\n", err)
    }

    // 检查数据桶是否存在
    exists, err := client.BucketExists(bucketName)
    if err != nil {
        log.Fatalf("查询数据桶错误:%v\n", err)
    }
    // 不存在则创建数据桶
    if !exists {
        err := client.MakeBucket(bucketName, location)
        if err != nil {
            log.Fatalf("创建数据桶错误:%v\n", err)
        }
    }

    // 打开文件
    file, err := os.OpenFile(filePath, os.O_RDONLY, 0444)
    if err != nil {
        log.Fatalf("打开文件失败:%v\n", err)
    }

    // 获取文件大小
    fileInfo, err := file.Stat()
    if err != nil {
        log.Fatalf("获取文件信息失败:%v\n", err)
    }
    fileSize := fileInfo.Size()

    // 创建进度条对象, 需要在参数中输入文件的大小
    progressBar := progress.NewUploadProgress(fileSize)

    _, err = client.PutObject(bucketName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream", Progress: progressBar})
    if err != nil {
        log.Fatalf("上传失败:%v\n", err)
    }
    log.Printf("上传成功!\n")
}
```

```
$ ls
big.rar  main.go
$ go run ./main.go
[====================] 1529094438 / 1529094438 100%
2019/08/10 16:00:26 上传成功!
```

### 下载使用进度条

```
// minio-go没有为下载接口提供查询进度功能, 所以只能在本地拷贝流时获取下载进度
// 你可以使用progress.CopyWithProgress代替io.Copy来将对象下载到本地文件
// 说明: progress.CopyWithProgress只是简单封装了io.Copy, 使之能打印下载进度

// 首先获取对象
object, err := client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
if err != nil {
    log.Fatalf("获取文件失败:%v\n", err)
}

// 然后使用io.CopyWithProgress拷贝到文件即可
n, err := progress.CopyWithProgress(file, object)
if err != nil {
    log.Fatalf("下载文件失败:%v\n", err)
}
```

#### 完整案例

main.go

```
package main

import (
	"log"
	"os"

	progress "github.com/markity/minio-progress"
	"github.com/minio/minio-go"
)

// minio客户端基本信息
var endpoint = "127.0.0.1:9000"
var accessKeyID = "Your-AccessKeyID"
var secretAccessKey = "Your-SecretAccessKey"
var secure = false

// 数据桶信息
var bucketName = "music"
var location = "us-east-1"

// 下载文件信息
var objectName = "big.rar"
var filePath = "./big.rar"

func main() {
    client, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
    if err != nil {
        log.Fatalf("创建客户端失败:%v\n", err)
    }

    object, err := client.GetObject(bucketName, objectName, minio.GetObjectOptions{})
    if err != nil {
        log.Fatalf("获取文件失败:%v\n", err)
    }

    file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatalf("打开文件失败:%v\n", err)
    }
    defer file.Close()

    if _, err := progress.CopyWithProgress(file, object); err != nil {
        log.Fatalf("下载文件失败:%v\n", err)
    }

    log.Printf("下载文件成功!\n")
}
```

```
$ ls
main.go
$ go run ./main.go
[====================] 1529094438 / 1529094438 100%
2019/08/10 17:43:00 下载文件成功
```