
box.cfg { listen = 3301 }
local grpc = require('tests')
local grpc_fiber = require('fiber').create(function() grpc:start() end)

TestService = {}
TestService.Procedure = function(context, req)
    return {message = req.message}
end
TestService.ExistsProcedure = function(context, req)
    return {message = req.message}
end
TestService.ErrorWhileCallProcedure = function(context, req)
    IHaveNoIdeaWhatIamDoing()
    return {message = req.message}
end
