compose = require("compose")

-- @TODO
create_task(
    "redis.start",
    function()
        print("define redis start command")
    end,
    function()
        print("redis start command call")
    end
)

compose.build()
compose.up("-d")
compose.exec("redis", "/bin/bash")