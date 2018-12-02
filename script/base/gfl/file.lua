
-- file, golang for lua

function FileModTime(path)
    mt, err = mfile.ModTime(path)
    if err ~= nil then
        LogError("ModTime err=%v", err)
        return 0, err
    end
    return mt, nil
end

