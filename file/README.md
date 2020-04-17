# IsExist
判断文件或者目录是否存在

#ExpandTilde 
路径转义 ~/xx -> /home/xxx

#ExpandTildeIgnoreErr 
ExpandTilde如果错误，不转义


#HomeDir
当前用户家目录

#IsAbs
判断路径是否是绝对路径

#TempFile
生产临时文件 -> ioutil.TempFile

#CreateFile
创建文件并且写入内容 -> ioutil.WriteFile

#IsDir
判断路径是否是目录

#IsFile
判断路径是否是文件

#SearchFile 
从指定的目录下搜索文件

#DirNum
路径下的目录数量

#SelfPath
当前程序运行的路径

# NewAtomicFileWriter
原子的写入
先写入一个临时文件，close的时候会Rename到指定的文件