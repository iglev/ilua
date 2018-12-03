module("SVC_mod1", package.seeall)

--[[
LogInfo("str=%v num=%v", "str", 12345)
LogError("str=%v num=%v", "str", 12345)
local subtab = {1, 2, 3}
local tab = {1, 2, "string", subtab, 1002, 1.2, "abc"}
LogInfo("tab=%v", tab)
]]

--[[
for i = 0, 10000 do
    LogInfo("tab=%v", tab)
end
]]

--[[
local mp = {
    Num = 123,
    Str = "string",
}
LogInfo("----abcd--mp=%v", mp)
]]

p = Person()
p.Name = "Lily"
p.Age = 3
p:Print()

no = Node()
no.Num = 1111
p:PrintNode(no)
p:PrintNode(p:GenNode(1234))

mymod.func1()
LogInfo("mymod.Num=%v mymod.StringVal=%v", mymod.Num, mymod.StringVal)

function fun(a, b)
    return tostring(a) .. tostring(b)
end

function fib(n)
    if n == 0 then
        return 0
    elseif n == 1 then
        return 1
    end
    return fib(n-1)+fib(n-2)
end
