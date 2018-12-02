
-- log, golang for lua

function LogInfo(format, ...)
    mlog.Info("|lualog: " .. format, ...)
end

function LogError(format, ...)
    mlog.Error("|lualog: " .. format, ...)
end

