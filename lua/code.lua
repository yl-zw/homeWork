local key =KEYS[1]
local Arg = ARGV[1]

local res=redis.call('get',key)
if  res then
    return res
else
    return false
end

function M(nu,nui)
    return nu*nui

end