Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"

  config.vm.synced_folder ".", "/go/src/github.com/rosenhouse/nsmess"

  # config.vm.provider "virtualbox" do |vb|
  #   vb.memory = "1024"
  # end

  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install -y golang
    echo GOPATH="/go" >> /etc/environment
    echo PATH="$PATH:/go/bin" >> /etc/environment
  SHELL
end
