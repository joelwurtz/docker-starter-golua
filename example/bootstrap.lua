-- This will be used to prefix all docker objects (network, images, containers)
project_name = 'app'
-- This is the root domain where the app will be available
-- The "frontend" container will receive all the traffic
root_domain = project_name .. '.test'
-- This contains extra domains where the app will be available
-- The "frontend" container will receive all the traffic
extra_domains = {}
-- This is the host directory containing your PHP application
project_directory = 'application'

-- Usually, you should not edit the file above this point
php_version = '7.4'
docker_compose_files = {
    'docker-compose.yml',
    'docker-compose.worker.yml',
}
services_to_build_first = {
    'php-base',
    'builder',
}

dinghy = False
power_shell = False
user_id = 1000
root_dir = '.'
start_workers = False

path = os.getenv("PATH")

function file_exists(name)
    local f= io.open(name,"r")

    if f~=nil then
        io.close(f)
        return true
    else
        return false
    end
end

require("tasks")

