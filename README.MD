# Golang 调用Github API
- 创建token.json 文件，里面添加github的token，内容如下：
```json
{
  "token": "xxxxxxxxxxxx"
}
```

 

```golang
// 获取用户基本信息
go run main.go user rockyzsu // rocky 为你打算获取的用户

// 获取所有粉丝，把粉丝按照关注数排序，默认取前10个
go run main.go fans rockyzsu

// 获取所有仓库数据
go run main.go repo rockyzsu
```
