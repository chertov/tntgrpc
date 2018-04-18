
local fiber = require('fiber')
box.cfg {
    listen = 3301
}

ExampleService = {}
ExampleService.SayHello = function(context, req)
    io.write('Hello, ' .. req.name .. '!\n')
    io.flush()
    return { message = 'Hello, ' .. req.name .. '!' }
end

local grpc = require('mytntgrpclib')
local grpc_fiber = fiber.create(function() grpc:start("0.0.0.0:50051") end)
