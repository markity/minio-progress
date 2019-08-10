# 为minio上传下载提供便捷的进度条接口

使用minio-progress快速为minio上传下载添加进度条

```
[===========         ] 856293376 / 1529094438 56%
```

## 从github安装

```
$ go get github.com/markity/minio-progress
```

## 快速开始

main.go

```
package main

import (
	"log"
	"os"

    // 导入进度条包
	progress "github.com/markity/minio-progress"

    // 导入minio SDK
    "github.com/minio/minio-go/v6"
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
	progressBar := progress.NewProgress(fileSize)

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