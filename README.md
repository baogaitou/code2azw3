# Code2Azw3
Convert golang project to a azw3 file(Or mobi file).
Note:
- Currently only support golang.
- Replace cover.png before generate azw3 file.(Size: 560px * 800px)
- Fonts SfMono、InputMono、FiraCode is recommend.

### Useage
Complier code2azw3:
```
go build .
```

Generate source html file:
```
./code2azw3 -dir ~/go/src/github.com/google/uuid -name uuid
```

Generate azw3:
```
sh build.sh -n uuid
```

Copy uuid/uuid.azw3 to kindle & happy reading ~

### Screenshots
![TOC](https://github.com/baogaitou/code2azw3/blob/master/screenshots/00.png)
![TOC](https://github.com/baogaitou/code2azw3/blob/master/screenshots/01.png)
![Code](https://github.com/baogaitou/code2azw3/blob/master/screenshots/02.png)
![Code](https://github.com/baogaitou/code2azw3/blob/master/screenshots/03.png)
![Code](https://github.com/baogaitou/code2azw3/blob/master/screenshots/04.png)
![Code](https://github.com/baogaitou/code2azw3/blob/master/screenshots/05.png)



# Code2Azw3 中文说明
为了方便在 Kindle 上阅读一些优秀的 golang 代码, 制作了这个小工具.
- 目前只适用于 golang
- templates 目录下的 cover.jpg 是默认封面,可自行替换.
- 推荐使用 SfMono、InputMono、FiraCode 等字体并横屏使用.

### 使用
编译
```
go build .
```

生成源文件
```
./code2azw3 -dir ~/go/src/github.com/google/uuid -name uuid
```

生成 azw3
```
sh build.sh -n uuid
```

至此,我们生成了 uuid.azw3 和 uuid 目录下的一堆源文件.将 uuid.azw3 拷贝至 kindle 的 documents 目录内即可.