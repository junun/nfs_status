package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 指标结构体
type Metrics struct {
	metrics 	map[string]*prometheus.Desc
	mountpath 	[]string
	mutex   	sync.Mutex
}

/**
 * 函数：newGlobalMetric
 * 功能：创建指标描述符
 */
func newGlobalMetric (
	namespace string,
	metricName string,
	docString string,
	labels []string ) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}

/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics(namespace string, path string) *Metrics {
	volumes := strings.Split(path, ",")
	if len(volumes) < 1 {
		fmt.Println("No NFS storage mount path given.  Path: %v", path)
	}
	return &Metrics{
		mountpath: volumes,
		metrics: map[string]*prometheus.Desc{
			"status_metrics": newGlobalMetric(namespace, "status_metrics","The description of status_metrics", []string{"mountPath"}),
		},
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()  // 加锁
	defer c.mutex.Unlock()
	for _, v := range c.mountpath {
		mockGaugeMetricData := c.GenerateData(v)
		for host, currentValue := range mockGaugeMetricData {
			ch <- prometheus.MustNewConstMetric(c.metrics["status_metrics"], prometheus.GaugeValue, float64(currentValue), host)
		}
	}

}

/**
 * 函数：GenerateMockData
 * 功能：生成模拟数据
 */
func (c *Metrics) GenerateData(path string) (mockGaugeMetricData map[string]int) {
	mockGaugeMetricData = map[string]int{
		path: CheckNfsStatus(path),
	}
	return
}

func CheckNfsStatus(path string) int {
	pre_cmd := exec.Command("nfsiostat", path)

	stdout, e := pre_cmd.StdoutPipe()
	if e != nil {
		log.Fatal(e)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if e := pre_cmd.Start(); e != nil {
		log.Fatal(e)
	}
	// 读取输出结果
	opBytes, e := ioutil.ReadAll(stdout)
	if e != nil {
		log.Fatal(e)
	}

	match, _ := regexp.MatchString("mounted", string(opBytes))
	if match {
		checkFile :=  path + "/check_nfs_status"
		cmd := exec.Command("touch", checkFile)
		cmd.Start()

		timer := time.AfterFunc(time.Duration(1) * time.Second, func() {
			err := cmd.Process.Kill()
			if err != nil {
				panic(err) // panic as can't kill a process.
			}
		})
		err := cmd.Wait()
		timer.Stop()

		// read error from here, you will notice the kill from the
		if err !=nil {
			return 11
		}
	} else {
		return 10
	}

	return 1
}


