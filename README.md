# Tokamak

A straightforward tool that utilizes CGroup to regulate system users' resource usage, including CPU and memory. This 
tool finds application in scenarios like shared HPC cluster login nodes and multi-user servers.


## Features

- Support CGroup v1 and v2, auto detect
- Limit the system users' CPU usage 
- Limit the system users' memory usage


## Installation

You can download the rpm file and directly install it.

```bash
  rpm -ivh tokamakd-1.0.0-1.el9.x86_64.rpm
```

Or manually compile from the source file.

```bash
  git clone https://github.com/sonmihpc/tokamak.git
  cd tokamak
  make build
  make install
```

## Usage/Configure

How to enable CGroup v1 in Rocky Linux 9?
```bash
  grubby --update-kernel=ALL --args="systemd.unified_cgroup_hierarchy=0 systemd.legacy_systemd_cgroup_controller"
  systemctl reboot
```

How to enable CGroup V2 in Rocky Linux 9?
```bash
  grubby --update-kernel=ALL --args="systemd.unified_cgroup_hierarchy=1"
  systemctl reboot
```

After install the tokamakd, you can start the service.
```bash
  systemctl start tokamakd
  systemctl enable tokamakd
```

You also can adjust the setup by edit the /etc/tokamak/config.yaml

```
cgroup:
  check-interval-ms: 1000 # 1000ms
  user-cpu-percent: 30    # 30% every user can consume the largest CPU percent of all logic cores.
  user-mem-percent: 30    # 30% every user can consume the largest memory percent of all memory.
  disable-oom-killer: true
  exclude-uids:           # which user exclude from restriction
    - 1240
    - 1250
```


## Feedback

If you have any feedback, please reach out to us at sonmihpc@gmail.com.

