# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
    config.vm.box = "ubuntu-trusty"
    config.vm.box_url = "https://cloud-images.ubuntu.com/vagrant/trusty/current/trusty-server-cloudimg-amd64-vagrant-disk1.box"

    # database machines
    config.vm.define "database-0" do |db0|
        db0.vm.network :private_network, ip: "192.168.30.30"
        db0.vm.provider :virtualbox do |v|
            v.name = "database-0"
            v.customize ["modifyvm", :id, "--memory", 2048]
        end
    end

    # api/webserver machines
    config.vm.define "web-0" do |web0|
        web0.vm.network :private_network, ip: "192.168.20.20"
        web0.vm.provider :virtualbox do |v|
            v.name = "web-0"
            v.customize ["modifyvm", :id, "--memory", 2048]
        end
    end

    # ansible config
    config.vm.provision :ansible do |ansible|
        # list of machines to be considered API machines
        api = ["web-0"]

        # list of machines in the database cluster
        databases = ["database-0"]

        # copy one and concat to get "all" server list
        # used in "common" role to increase parallelism
        all = api.dup
        all.concat(databases)

        ansible.groups = {
            "db_servers" => databases,
            "api_servers" => api,
            "all_servers" => all
        }
        ansible.playbook = "ansible/false.yml"

        ## note: I do not use vagrant to provision
        ##       It currently does not allow a way
        ##       to bring up machines THEN run ansible
        ##       However, below is what you would use

        # ansible.limit = 'all'
        # ansible.host_key_checking = false
        # ansible.playbook = "ansible/all.yml"
    end
end
