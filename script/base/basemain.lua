module("G_BASE_MAIN", package.seeall)

luafiles = {
    {
        dir = G_BASE_DIR .. "gfl/",
        files = {
            "log.lua",
        },
        decode = "none",
    },
    {
        dir = G_BASE_DIR .. "lfg/",
        files = {
        },
        decode = "none",
    },
}
