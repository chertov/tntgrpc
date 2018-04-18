
box.cfg { listen = 3301 }
local grpc = require('tests')
local grpc_fiber = require('fiber').create(function() grpc:start() end)

TestService = {}
TestService.MessageNotLuaTable = function(context, req)
    return { rr = { outsider = 'This Is Not Lua Table' } }
end
TestService.NilReturn = function(context, req)
    return;
end
TestService.IncorrectFieldType = function(context, req)
    return { rr = { outsider = { message = 123 } } }
end
TestService.NonExistingField = function(context, req)
    return { rr = { outsider = { message = "hello, world", thisFieldIsNtExists = "hello, world" } } }
end
