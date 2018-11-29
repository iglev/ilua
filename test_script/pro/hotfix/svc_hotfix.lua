require("hotfix.data.DATA_hotfix")

function HotFix()
    -- curr = os.time()
    -- if curr >= (DATA_hotfix.lastUpdateTime+10) then
        for _, v in ipairs(G_MAIN.luafiles) do
            if type(v) == 'table' and type(v.files) == 'table' then
                for _, filename in ipairs(v.files) do
                    fullname = v.dir .. filename
                    dofile(fullname)
                end
            end
        end
    -- end
end

