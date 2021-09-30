package so
/*
1. 编写go模块，并导出需要给lua使用的函数：
//export add
func add(operand1 int, operand2 int) int {
    return operand1 + operand2
}

2. 将go模块编译为静态库：
go build -buildmode=c-shared -o example.so example.go

3. 编写lua文件，加载自己的.so文件：
local lua2go = require('lua2go')
local example = lua2go.Load('./example.so')

4. 在lua文件与头文件模块中注册导出的函数：
lua2go.Externs[[
extern GoInt add(GoInt p0, GoInt p1);
]]

5. 在lua文件中调用导出的函数并将结果转化为lua格式的数据：
local goAddResult = example.add(1, 1)
local addResult = lua2go.ToLua(goAddResult)
print('1 + 1 = ' .. addResult)
*/