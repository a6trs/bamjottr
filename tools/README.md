Automated Tools
===============

These little tools may be helpful when writing code. You can choose any one(s) of them to use, but remember to read the 'requirements' section of a tool if you want to use it.

Currently, there is only one tool: `serve.lua`.

`serve.lua`
-----------
Implements auto-rebuilding for this project. Can also be use in other golang projects.  
With this script you'll be able to run the previously-built program directly if no changes are made to the code instead of rebuilding the whole project every time. Simply replace all `go run` calls with this script.

**Requirements** 
* Install [Lua 5.3](http://www.lua.org/versions.html#5.3).
> Remember to run `make install` when building Lua from source.

* Install [LuaFileSystem](http://keplerproject.github.io/luafilesystem).

**Usage**

*NIX
```
$ tools/serve.lua
```

Windows (not tested)
```
> lua tools/serve.lua
```
