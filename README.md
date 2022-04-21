# file-stats

* collects the size of object at the given path

## outline

* gets files and directories at a given path (`-baseDirectory`)
* reports files and directories
* exposes results for prometheus to consume ( creates a file under `-targetFilePath`)

## Links

* https://stackoverflow.com/questions/14694088/is-it-safe-for-more-than-one-goroutine-to-print-to-stdout

## usage

* build binary from source code
  `go build -o bin/file-stats src/*.go`
* example call
  `/opt/file-stats -baseDirectory /opt/storage/ -outputFilePath /opt/prometheus_node_exporter_textfile/`
  * would evaluate `/opt/storage` and report the results to `/opt/prometheus_node_exporter_textfile/opt_storage.prom`
