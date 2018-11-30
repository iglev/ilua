
-- G_MAIN_FILE main.lua
G_MAIN_FILE = G_SCRIPT_DIR .. "main.lua"

local extraPath = {
    package.path,
    G_SCRIPT_DIR .. "?.lua",
}
package.path = table.concat(extraPath, ";")

