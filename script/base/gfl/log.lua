
-- log, golang for lua

function LogInfo(format, ...)
    pargs = {...}
    table.insert(pargs, #pargs)
    MLog.Info("|lualog: " .. format, pargs)
end

function LogError(format, ...)
    pargs = {...}
    table.insert(pargs, #pargs)
    MLog.Error("|lualog: " .. format, pargs)
end

