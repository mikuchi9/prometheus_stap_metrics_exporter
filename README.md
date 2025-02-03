## Prometheus Custom Metrics: Fork and Clone Syscall Details Exporter Using a SystemTap Script

To use this exporter, you need to have `go` and `systemtap` installed on your system. This exporter relies on **kernel debug symbols**. If they aren't available, it won't work. The kernel module created by `stap` uses kernel probes for the fork and clone system calls. Every time these events occur, it prints details to standard output, which is then directed into a pipe and consumed by this Go exporter.

1. Compile the `custom_metrics` kernel module. This will create a module named `custom_metrics.ko` in the current directory:

   `sudo stap -p4 -m custom_metrics custom_metrics.stp`

2. Run the exporter:
   
   `go run exporter.go`

3. Verify the metrics endpoint on port 9100:

   `curl http://localhost:9100/metrics`
   
* *This was tested on Ubuntu 20.04.06 LTS*\
  *kernel: 5.4.0-204-generic*\
  *kernel debug symbols: linux-image-5.4.0-204-generic-dbgsym*\
  *SystemTap version:*\
  `Systemtap translator/driver (version 4.2/0.176, Debian version 4.2-3ubuntu0.1 (focal))`
  `Copyright (C) 2005-2019 Red Hat, Inc. and others`
  `This is free software; see the source for copying conditions.`
  `tested kernel versions: 2.6.32 ... 5.4-rc6`
  `enabled features: AVAHI BPF LIBSQLITE3 NLS NSS`
