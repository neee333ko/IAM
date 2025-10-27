# User接口
## 1. 创建用户
### 1.1 接口描述
创建用户。
### 1.2 请求方法
POST /v1/users
### 1.3 输入参数
**Body参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
metadata | 是 | ObjectMeta | 元数据
status | 否 | int | 用户启用状态
nickname | 是 | string | 用户别名
password | 是 | string | 密码
email | 是 | string | 邮箱地址
phone | 否 | string | 电话号码
### 1.4 输出参数
参数名称 | 类型 | 描述
------------ | ------------- | ------------
metadata | ObjectMeta | 元数据
status | int | 用户启用状态
nickname | string | 用户别名
password | string | 密码
email | string | 邮箱地址
phone | string | 电话号码
isAdmin | int | 用户身份
totalPolicy | int64 | 用户策略数
LoginedAt | time.Time | 上次登录时间
### 1.5 用例
## 2. 更新用户
### 2.1 接口描述
更新用户信息。
### 2.2 接口方法
PUT /v1/users/:name
### 2.3 输入参数
**Path参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
name | 是 | string | 用户名

**Body参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
metadata | 否 | ObjectMeta | 元数据
status | 否 | int | 用户启用状态
nickname | 否 | string | 用户别名
password | 否 | string | 密码
email | 否 | string | 邮箱地址
phone | 否 | string | 电话号码
### 2.4 输出参数
参数名称 | 类型 | 描述
------------ | ------------- | ------------
metadata | ObjectMeta | 元数据
status | int | 用户启用状态
nickname | string | 用户别名
password | string | 密码
email | string | 邮箱地址
phone | string | 电话号码
isAdmin | int | 用户身份
totalPolicy | int64 | 用户策略数
LoginedAt | time.Time | 上次登录时间
### 2.5 用例
## 3. 修改密码
### 3.1 接口说明
修改用户密码。
### 3.2 接口方法
PUT /v1/users/:name/change-password
### 3.3 输入参数
**Path参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
name | 是 | string | 用户名

**Body参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
oldPassword | 是 | string | 旧密码
newPassword | 是 | string | 新密码

### 3.4 输出参数
参数名称 | 类型 | 描述
------------ | ------------- | ------------
metadata | ObjectMeta | 元数据
status | int | 用户启用状态
nickname | string | 用户别名
password | string | 密码
email | string | 邮箱地址
phone | string | 电话号码
isAdmin | int | 用户身份
totalPolicy | int64 | 用户策略数
LoginedAt | time.Time | 上次登录时间
### 3.5 用例
## 4. 返回用户信息
### 4.1 接口说明
返回单位用户信息。
### 4.2 接口方法
GET /v1/users/:name
### 4.3 输入参数
**Path参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
name | 是 | string | 用户名
### 4.4 输出参数
参数名称 | 类型 | 描述
------------ | ------------- | ------------
metadata | ObjectMeta | 元数据
status | int | 用户启用状态
nickname | string | 用户别名
password | string | 密码
email | string | 邮箱地址
phone | string | 电话号码
isAdmin | int | 用户身份
totalPolicy | int64 | 用户策略数
LoginedAt | time.Time | 上次登录时间
### 4.5 用例
## 5. 返回用户列表
### 5.1 接口说明
返回用户列表。
### 5.2 接口方法
GET /v1/users
### 5.3 输入参数
**Query参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
limit | 否 | int | 返回上限数
offset | 否 | int | 偏移数
fieldSelector | 否 | string | 字段查询
### 5.4 输出参数 
参数名称 | 类型 | 描述
------------ | ------------- | ------------
metadata | ListMeta | 元数据
items | []*User | 用户列表
### 5.5 用例
## 6. 删除用户
### 6.1 接口说明
删除单个用户。
### 6.2 接口方法
DELETE /v1/users/:name
### 6.3 输入参数
**Path参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
name | 是 | string | 用户名
### 6.4 输出参数
Null
### 6.5 用例
## 7. 批量删除用户
### 7.1 接口说明
批量删除用户，仅对管理员开放。
### 7.2 接口方法
DELETE /v1/users
### 7.3 输入参数
**Body参数**
参数名称 | 必选 | 类型 | 描述
------------ | ------------- | ------------ | -------------
names | 是 | []string | 用户名数组
### 7.4 输出参数
Null
### 7.5 用例
