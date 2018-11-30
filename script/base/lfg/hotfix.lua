require("data.hotfix")

-- hotfix, lua for golang

function LFGHotFix()
    tryReloadAll("G_BASE_MAIN")
    tryReloadAll("G_MAIN")
end

function isNeedReload(fullpath, modTime)
    if not DATAHotfix.LastTimeCache[fullpath] then
        return false
    end
    if DATAHotfix.LastTimeCache[fullpath] ~= modTime then
        return true
    end
    return false
end

function tryReload(fullpath)
    modTime, err = FileModTime(fullpath)
    if err then
        LogError("FileModTime fail, err=%v", err)
        return
    end
    if isNeedReload(fullpath, modTime) then
        err = dofile(fullpath)
        if err then
            LogError("dofile fail file=%v err=%v", fullpath, err)
            return
        end
    end
    DATAHotfix.LastTimeCache[fullpath] = modTime
end

function tryReloadAll(modName)
    local m = _G[modName]
    if m == nil then
        return
    end
    if m.luafiles == nil or type(m.luafiles) ~= 'table' then
        LogError("mod=%v not found luafiles", modName)
        return
    end
    for _, one in ipairs(m.luafiles) do
        if type(one.dir) == 'string' and type(one.files) == 'table' then
            for _, file in ipairs(one.files) do
                tryReload(one.dir .. file)
            end
        end
    end
end


