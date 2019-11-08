# nfs_status


```
// 使用golang编写Prometheus Exporter
1) 依赖
go get -v github.com/prometheus/client_golang/prometheus/promhttp
go get -v github.com/prometheus/client_golang/prometheus


2） 定义metrics

3） 实现Describe 和 Collect

4）编译 go build
```

```
 ~/Documents/go/src/nfs_status   master  ./nfs_status -h
Usage of ./nfs_status:
  -metric.namespace string
    	Prometheus metrics namespace, as the prefix of metrics name (default "nfs")
  -nfs.storage-path string
    	Path to nfs storage volume. (default "/data/images/lighting")
  -web.listen-port string
    	An port to listen on for web interface and telemetry. (default "9001")
  -web.telemetry-path string
    	A path under which to expose metrics. (default "/metrics")
    	

-nfs.storage-path 指定nfs本地挂载的路径。

该Exporter实现的功能和下面shell功能一样
#!/bin/bash
#
# author junun
#
#先检查是否mount
nfsiostat $1 | grep  'mounted' > /dev/null
if [ $? -eq 0 ];then
    # 检查改目录下的文件是否可写
    timeout 1  touch $1/check_nfs_status
    if [ $? -gt 0 ];then
        echo -n 11
    else
        echo -n 1
    fi
else
    echo -n 10
fi
```