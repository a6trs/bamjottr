#!/usr/local/bin/lua
require 'lfs'
inspect = require('tools/inspect')

dates = {}

function register_date(file)
    local attr = lfs.attributes(file)
    if attr then dates[file] = attr.modification
    else print('register_date: No such file: ' .. file) end
end

function register_date_dir(dir)
    for file in lfs.dir(dir) do
        if file ~= '.' and file ~= '..' then register_date(dir .. file) end
    end
end

register_date 'main.go'
register_date_dir 'soil/'
register_date_dir 'grass/'

-- Prevent crash on first run
local file = io.open('tools/record.lua', 'a')
file:close()
-- Load previous data
prev_dates = dofile('tools/record.lua')
-- prev_dates is nil on first run
if prev_dates == nil then prev_dates = {} end
needs_rebuild = false
for k, v in pairs(prev_dates) do if dates[k] ~= v then needs_rebuild = true break end end
if not needs_rebuild then
    for k, v in pairs(dates) do if prev_dates[k] ~= v then needs_rebuild = true break end end
end

if needs_rebuild then
    print('Needs rebuild. Working around...')
    local _, _, exitcode = os.execute('go build main.go')
    if exitcode ~= 0 then
        print('Errors... Exiting T^T')
        return exitcode
    end
    file = io.open('tools/record.lua', 'w')
    file:write('return ')
    file:write(inspect.inspect(dates))
    file:close()
end
os.execute('./main')
