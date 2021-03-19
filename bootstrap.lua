compose = require("compose")
task = require("task")

task.create(
    "redis.start",
    function()
        task.set_short_description("start redis service")
    end,
    function()
        compose.up("redis", "-d")
    end
)
task.create(
    "redis.stop",
    function()
        task.set_short_description("stop redis service")
    end,
    function()
        compose.down("redis")
    end
)