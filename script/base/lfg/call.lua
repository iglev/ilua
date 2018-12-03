
-- call, lua for golang

function LFGCall(functab, ...)
    local func = nil
    if functab.m and functab.f then
        func = _G[functab.m][functab.f]
        if not func then
            return {err="not found function"}
        end
    else
        func = _G[functab.f]
    end
    r = func(...)
    return {r=r}
end

