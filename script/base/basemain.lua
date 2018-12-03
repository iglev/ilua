module("G_BASE_MAIN", package.seeall)

luafiles = {
    {
        dir = G_BASE_DIR .. "gfl/",
        files = {
            "log.lua",
            "file.lua",
        },
        decode = "none",
    },
    {
        dir = G_BASE_DIR .. "lfg/",
        files = {
            "hotfix.lua",
            "call.lua",
        },
        decode = "none",
    },
}
