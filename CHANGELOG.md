# 2024/01/07 编写了关于 cmdb 相关的部分代码 

* api/cmdb.ts
* types/cmdb.ts

目前可以获取这些数据了，但是还不能在前端展示。

又编写了一个脚本，计划用脚本方便的生成 hosts 的 json 数据，然后 mock 时，读取 json 进行返回。

# 2024/01/14 cmdb 中添加了关于分页相关的代码

# 2024/01/28 后端代码想改为 gorm 未果

# 2024/02/03 调用真正的后端成功
* axios.get 的代码还得改改，否则类型报warn

# 2024/03/03 增加了 host 的 create 和 delele 功能

* update 的后端代码完成了部分
* create 的还需要完善，比如表单校验


# 2024/03/16 完成了 update 的功能

* 完成了 update 的功能
* 更新了 prettier 配置
