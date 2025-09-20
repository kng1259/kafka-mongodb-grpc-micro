home = ENV['HOME']
nodes = {
    "server-node" => { count: 1, cpus: 2, storage: "30GB", memory: 3072 },
    "agent-node" => { count: 2, cpus: 2, storage: "60GB", memory: 3072 }
}

node_names = nodes.map do |node_type, config|
  (1..config[:count]).map { |n| "#{node_type}-#{n}" }
end

Vagrant.configure("2") do |config|

    nodes.each do |node_prefix, specs|
        (1..specs[:count]).each do |i|
            node_name = "#{node_prefix}-#{i}"
            config.vm.define node_name do |vm|
                vm.vm.hostname = node_name
                vm.vm.box = "cloud-image/ubuntu-24.04"
                vm.vm.disk :disk, name: "boot", size: "30GB", primary: true
                vm.vm.box_check_update = true

                vm.vm.network :private_network, type: "dhcp",
                    :libvirt__network_name => "vagrant-routed",
                    :libvirt__forward_mode => "nat",
                    :libvirt__dhcp_start => "172.28.128.4",
                    :libvirt__dhcp_stop => "172.28.128.100"

                vm.vm.provider :libvirt do |domain|
                    domain.default_prefix = "kmgm-"
                    domain.cpus = specs[:cpus]
                    domain.memory = specs[:memory]
                    domain.storage :file, :size => specs[:storage], :type => 'qcow2', :disk_ext => 'ext4'
                end

                vm.vm.provision "shell", inline: <<-SHELL
                    sudo apt-get update
                    sudo apt-get upgrade
                    sudo apt-get install -y python3
                    # echo -e "n\np\n1\n\n\nw" | sudo fdisk /dev/vdb
                    # sudo mkfs.ext4 /dev/vdb1
                    # sudo mkdir -p /mnt/new_disk
                    # echo '/dev/vdb1 /mnt/new_disk ext4 defaults 0 0' | sudo tee -a /etc/fstab
                    # sudo mount -a
                SHELL

            end
        end
    end

    config.vm.provision "ansible" do |ansible|
        ansible.verbose = "v"
        ansible.config_file = "ansible.cfg"
        ansible.limit = "all"
        ansible.force_remote_user = false
        ansible.groups = {
            "poc4k_first_server_node" => [node_names[0][0]],
            "poc4k_agent_nodes" => node_names[1],
            "poc4k" => ["poc4k_first_server_node", "poc4k_agent_nodes"],
        }
        ansible.extra_vars = {
            ansible_user: 'vagrant',
        }
        ansible.playbook = "playbooks/rke2-local-setup.yaml"
    end

end
