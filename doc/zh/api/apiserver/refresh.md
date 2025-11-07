# 刷新token
## 1. 刷新token
### 1.1 接口描述
刷新token。
### 1.2 接口方法
POST /refresh
### 1.3 输入参数
**Body参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
refresh_token | 是 | string | 刷新token
### 1.4 输出参数
参数名称 | 类型 | 描述
------------ | ------------- | ------------
access_token | string | 通行token
token_type | string | token类型
refresh_token | string | 刷新token
expires_at | int64 | 到期日期
created_at | int64 | 创建日期
### 1.5 用例 