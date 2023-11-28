# Tokamak

A simple tool to restrict the system users' resource such as CPU or memory by CGroup. It can be used in some situations 
such as HPC cluster common login node, multi-users server.


## Features

- Support CGroup v1
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
  disable-oom-killer: false
```


## Feedback

If you have any feedback, please reach out to us at sonmihpc@gmail.com.

