
-- exec dir
G_DIR = "./"

-- script root dir
G_SCRIPT_ROOT_DIR = "./script/"

-- script base dir
G_BASE_DIR = G_SCRIPT_ROOT_DIR .. "base/"

-- script basemain file
G_BASE_MAIN_FILE = G_BASE_DIR .. "basemain.lua"

local extraPath = {
    package.path,
    G_BASE_DIR .. "?.lua",
    G_BASE_DIR .. "gfl/?.lua",
    G_BASE_DIR .. "lfg/?.lua",
    -- ...
}

package.path = table.concat(extraPath, ";")
require(G_SCRIPT_ROOT_DIR .. "conf")
-- require("std")
