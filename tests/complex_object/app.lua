
box.cfg { listen = 3301 }
local grpc = require('tests')
local grpc_fiber = require('fiber').create(function() grpc:start() end)


local json=require('json')

local reply = {
    simple1 = "very simple1",
    simple2 = 123,
    beOrNotToBe = true,
    nextStation = "SRILANKA",
    mustBeLocal = { ["local"] = "ok" },
    mustBeOuter = {
        outer = true,
        mapInMap = {
            outsider1 = { ["outer"] = true, ["name"] = "outer 1" },
            outsider2 = { ["outer"] = true, ["name"] = "outer 2" },
            outsider3 = { ["outer"] = true, ["name"] = "outer 3" },
        }
    },
    OneTwoThree = {1, 2, 3},
    OneTwoThreeStrs = {"1", "2", "3"},
    -- nextStations = {"USA", "INDONESIA", 2, 1}, -- It's Ok! string or int
    outsiders = {
        { ["outer"] = true, ["name"] = "outer 1" },
        { ["outer"] = true, ["name"] = "outer 2" },
        { ["outer"] = true, ["name"] = "outer 3" },
    },

    friends = {
        [42]     = { ["local"] = "friend 1" },
        [146]    = { ["local"] = "friend 2" },
        [100500] = { ["local"] = "friend 3" },
    },
    enemies = {
        enemy1 = { ["outer"] = true, ["name"] = "enemy 1" },
        enemy2 = { ["outer"] = true, ["name"] = "enemy 2" },
        enemy3 = { ["outer"] = true, ["name"] = "enemy 3" },
    },
}


TestService = {}
TestService.ItsOk1 = function(context, req)

    io.write("req " .. dump(req) .. "\n")
    local res = reply
    reply.enemies["enemy100500"] = req.outsiders[3].mapInMap["outsider2"]
    return reply
end
TestService.ItsOk2 = function(context, req)
    return {
        simple1 = "very simple1",
        simple2 = 123,
        beOrNotToBe = true,
        nextStation = "SRILANKA",
        mustBeLocal = { ["local"] = "ok" },
        mustBeOuter = {
            outer = true,
            mapInMap = {
                outsider1 = { ["outer"] = true, ["name"] = "outer 1" },
                outsider2 = { ["outer"] = true, ["name"] = "outer 2" },
                outsider3 = { ["outer"] = true, ["name"] = "outer 3" },
            }
        },
        OneTwoThree = {1, 2, 3},
        OneTwoThreeStrs = {"1", "2", "3"},
        -- nextStations = {"USA", "INDONESIA", 2, 1}, -- It's Ok! string or int
        outsiders = {
            { ["outer"] = true, ["name"] = "outer 1" },
            { ["outer"] = true, ["name"] = "outer 2" },
            { ["outer"] = true, ["name"] = "outer 3" },
        },

        friends = {
            [42]     = { ["local"] = "friend 1" },
            [146]    = { ["local"] = "friend 2" },
            [100500] = { ["local"] = "friend 3" },
        },
        enemies = {
            enemy1 = { ["outer"] = true, ["name"] = "enemy 1" },
            enemy2 = { ["outer"] = true, ["name"] = "enemy 2" },
            enemy3 = { ["outer"] = true, ["name"] = "enemy 3" },
        },
    }
end

function dump(o, indx)
    local pad = true
    if indx == nil then indx = 1 end
    if type(o) == 'table' then
        local count = 0
        for _ in pairs(o) do count = count + 1 end
        if count == 0 then return '{}' end
        local s = '{'
        if pad then s = s .. '\n' end
        for k,v in pairs(o) do
            for i=1,indx do s = s .. '    ' end
            if type(k) ~= 'number' then k = '"'..k..'"' end
            s = s .. '['..k..'] = ' .. dump(v, indx+1) .. ','
            if pad then s = s .. '\n' end
        end
        for i=1,indx-1 do s = s .. '    ' end
        return s .. '}'
    else
        if type(o) == 'string' then
            return '"' .. tostring(o) .. '"'
        else
            return tostring(o)
        end
    end
end
