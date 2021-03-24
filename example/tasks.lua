task = require("task")
compose = require("compose")

task.create(
    "build",
    function()

    end,
    function()
        args = {
            "--build-arg",
            "PROJECT_NAME=" .. project_name,
            "--build-arg",
            "USER_ID=" .. user_id,
            "--build-arg",
            "PHP_VERSION=" .. php_version,
        }

        for _, service in ipairs(services_to_build_first) do
            serviceArgs = args
            table.insert(serviceArgs, service)

            compose.build(unpack(serviceArgs))
        end

        compose.build(unpack(args))
    end
)