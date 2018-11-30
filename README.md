# ilua
base on github.com/yuin/gopher-lua. 
1) load lua files
2) hotfix
3) go call lua func, lua call go func

# script 目录
base:           基础libs( lua for go, go for lua)
pro:            工程模块( main.lua)
args.lua:       入口参数(第一个加载lua文件)
conf.lua:       配置相关
load_pre.lua:	  package.path相关的全局变量设置
load_after.lua: 工程main.lua入口设置


