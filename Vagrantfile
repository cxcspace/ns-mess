# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"

  config.vm.synced_folder ".", "/go/src/github.com/rosenhouse/ns-mess"

  config.vm.provision "shell", inline: <<-SHELL
    set -e -x -u

    apt-get update -y || (sleep 40 && apt-get update -y)
    apt-get install -y git

    GO_VERSION=1.9.4

    wget -qO- https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz | tar -C /usr/local -xz

    echo "GOPATH=/go" >> /etc/environment
    echo "PATH=$PATH:/usr/local/go/bin:/go/bin" >> /etc/environment

    source /etc/environment
    export GOPATH
    go install github.com/rosenhouse/ns-mess
  SHELL
end
