
box.cfg { listen = 3301 }
local grpc = require('tests')
local grpc_fiber = require('fiber').create(function() grpc:start() end)

TestService = {}
TestService.MapIncorrectKeyType = function(context, req)
    return {
        ["map"] = {
            [42] = "value1",
            [146] = "value2",
            ["keyWithIncorrectType"] = "value3",
            [100500] = "value4",
        }
    }
end
TestService.MapIncorrectValueType = function(context, req)
    return {
        ["map"] = {
            ["key1"] = 1,
            ["key2"] = 2,
            ["keyWithIncorrectValueType"] = "IncorrectValueType",
            ["key4"] = 4,
            ["key5"] = 5
        }
    }
end
