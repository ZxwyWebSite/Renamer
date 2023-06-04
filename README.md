## ZxwyWebSite/Renamer
(当前处于测试阶段)，文件批量重命名工具
### 使用
+ 将程序添加到path
+ `ln -s /root/renamer-dev/renamer /usr/local/bin/renamer`
+ 在运行目录创建 `list.txt`，将要重命名成的名称写入
+ 每行一个文件名(最后一行不要回车)

示例 list.txt
```
第05话_王立骑士学院.mp4
第06话_试炼迷宫.mp4
第07话_天惠武姬的病.mp4
第08话_莉波护卫指令.mp4
第09话_新的魔印武具.mp4
第10话_天上人袭击计划.mp4
第11话_无双的见习骑士.mp4
```
+ 执行 `renamer -help` 查看帮助
```
Usage of renamer:
  -e string
        指定文件扩展名(不加点)
  -p string
        指定扫描目录(相对位置) (default "/root/renamer-dev")
```
### 其它
+ 测试阶段，没有错误处理，如有bug请上报