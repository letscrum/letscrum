# 自动生成 releasenotes 

该项目能扫描输入文件夹中所有yaml文件，并根据指定格式进行集成输出releasenote文件

## yaml文件格式

该项目会自动识别yaml文件中的kind、area和notes，其中kind和area必须是值类型，nontes必须是数组类型，否则将会报错

kind和area的值必须是模版中给定的选项其中之一，其他的值将会报错

## 指令介绍

通过运行 main.go 的方式，使用``notesDir``指定 notes 存放的位置，使用``outDir``指定输出releasenote文件的位置

输出文件名为``notesDir``文件夹路径的最后一个文件名。例入``notesDir``为``./changelogs/0.1.1``，则输出的releasenote文件名为``0.1.1``

样例指令如下：

```bash
go run ./tools/releasenoter/main.go --notesDir ./changelogs/x.x.x --outPath ./changelogs/tool/release_note
```

### 参数

* (optional) `--notesDir`  yaml nots存放的文件夹。默认是``./changelogs/0.1.1``

* (optional) `--outPath` 生成的 releasenotes 的存放位置。默认为是`./changelogs/tool/release_note`
