module("G_MAIN", package.seeall)

luafiles = {
    {
        dir = G_SCRIPT_DIR .. "hotfix/",
        files = {
            "svc_hotfix.lua",
        },
        decode = "none",
    },
    {
        dir = G_SCRIPT_DIR .. "mo1/",
        files = {
            "svc_mo1.lua",
        },
        decode = "none",
    },
    {
        dir = G_SCRIPT_DIR .. "mo2/",
        files = {
            "svc_mo2.lua",
        },
        decode = "confuse",
    }
}
