# IAM
### IAM 是一个基于 Go 语言开发的身份识别与访问管理系统，用于对资源访问进行授权。
#### 1.IAM-APISERVER 为控制流服务,包含对于User,Secret和Policy资源的CURD.
#### 2.IAM-AUTHZSERVER 为数据流服务,使用ladon对授权申请进行验证.
#### 3.IAM-PUMP 抓取redis中的授权记录存取到MongoDB中.
#### 详情API参考doc