# Chef attributes used to override the default settings on a node used to understand 
# - The current state of the node
# - What the state of the node was at the end of the rpevious Chef infra client run 
# - what the state should be at the end of a run
ubuntu_mirror = 'https://mirror.its.sfu.ca/mirror/ubuntu/'
ubuntu_release = 'bionic'
ubuntu_version = '18.04'
username = 'vagrant'
user_home = '/home/' + username
project_home = user_home + '/project/demos' # you may need to change the working directory to match your project


python3_packages = '/usr/local/lib/python3.6/dist-packages'
ruby_gems = '/var/lib/gems/2.5.0/gems/'


# Get Ubuntu sources set up and packages up to date.

template '/etc/apt/sources.list' do
  variables(
    :mirror => ubuntu_mirror,
    :release => ubuntu_release
  )
  notifies :run, 'execute[apt-get update]', :immediately
end
execute 'apt-get update' do
  action :nothing
end
execute 'apt-get upgrade' do
  command 'apt-get dist-upgrade -y'
  only_if 'apt list --upgradeable | grep -q upgradable'
end
directory '/opt'
directory '/opt/installers'


# Basic packages many of us probably want. Includes gcc C and C++ compilers.

package ['build-essential', 'cmake']

# Other core language tools you might want

package ['python3', 'python3-pip', 'python3-dev']  # Python
package 'golang-go'  # Go



# .NET Core

#remote_file '/opt/installers/packages-microsoft-prod.deb' do
#  source "https://packages.microsoft.com/config/ubuntu/#{ubuntu_version}/packages-microsoft-prod.deb"
#end
#execute 'dpkg -i /opt/installers/packages-microsoft-prod.deb' do
#  creates '/etc/apt/sources.list.d/microsoft-prod.list'
#  notifies :run, 'execute[apt-get update]', :immediately
#end
#package ['dotnet-sdk-3.1']


# NodeJS (more modern than Ubuntu nodejs package) and NPM

remote_file '/opt/installers/node-setup.sh' do
 source 'https://deb.nodesource.com/setup_14.x'
 mode '0755'
end
execute '/opt/installers/node-setup.sh' do
 creates '/etc/apt/sources.list.d/nodesource.list'
 notifies :run, 'execute[apt-get update]', :immediately
end
package ['nodejs']


# SWIG

#package 'swig'


# RabbitMQ-related things

package ['rabbitmq-server']

# Python pika library
execute 'pip3 install pika==1.1.0' do
 creates "#{python3_packages}/pika/__init__.py"
end
# Ruby bunny library
#execute 'gem install bunny -v 2.15.0' do
#  creates "#{ruby_gems}/bunny-2.15.0/Gemfile"
#end
# Go amqp library
# execute 'go get github.com/streadway/amqp github.com/google/uuid' do
#  cwd project_home 
#  user username
#  environment 'HOME' => user_home
#  creates user_home + '/go/src/github.com/streadway/amqp/README.md'
# end
# Java amqp library
#package 'librabbitmq-client-java'


# install flask 
execute 'pip3 install flask' do 
  command '/usr/bin/pip3 install flask'
end

# install rabit mq 
execute 'go get github.com/streadway/amqp' do 
  cwd '/home/vagrant/project/api/go/src/main'
  command '/usr/bin/go get github.com/streadway/amqp'
end

# install image/webp 
execute 'go get golang.org/x/image/webp' do 
  cwd '/home/vagrant/project/api/go/src/main'
  command '/usr/bin/go get golang.org/x/image/webp'
end

# install node modules 
execute 'npm install' do 
  cwd '/home/vagrant/project/client'
  command '/usr/bin/npm install'
end

# try running app with a bash resource
bash 'start react app' do 
  user 'root'
  cwd '/home/vagrant/project/client'
  code <<-EOH
    npm start &
  EOH
end
