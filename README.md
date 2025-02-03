## Prometheus Custom Metrics: Fork and Clone Syscall Details Exporter Using a SystemTap Script

To use this exporter, you need to have `go` and `systemtap` installed on your system. This exporter relies on **kernel debug symbols**. If they aren't available, it won't work. The kernel module created by `stap` uses kernel probes for the fork and clone system calls. Every time these events occur, it prints details to standard output, which is then directed into a pipe and consumed by this Go exporter.

1. Compile the `custom_metrics` kernel module. This will create a module named `custom_metrics.ko` in the current directory:

   `sudo stap -p4 -m custom_metrics custom_metrics.stp`

2. Run the exporter:
   
   `go run exporter.go`

3. Verify the metrics endpoint on port 9100:

   `curl http://localhost:9100/metrics`
   
* This was tested on Ubuntu 20.04.06 LTS

