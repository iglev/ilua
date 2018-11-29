
-- exec dir
G_DIR = "./"

-- script root dir
G_SCRIPT_ROOT_DIR = "./test_script/"

local extraPath = {
    package.path,
    -- ...
}

package.path = table.concat(extraPath, ";")
print(package.path)
require(G_SCRIPT_ROOT_DIR .. "conf")

