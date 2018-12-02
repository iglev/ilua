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

function fun()
    return 10
end
