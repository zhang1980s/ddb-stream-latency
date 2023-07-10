# ddb-stream-latency
Test the latency of DDB stream
---

### 测试环境：
一台EC2，用于生成压力。尺寸取决于压测规模。
两个DDB table（timer/dashboard）。timer表用于模拟数据写入和stream订阅，dashboard表用户保存时间差信息
一个lambda函数，用于订阅ddb stream，计算时间差和写入dashboard表


### 成本
取决于压测负载规模及DDB预置读写能力。

### 测试步骤：
0. 部署测试环境
创建两个ddb table，一个叫timer，一个叫dashboard。使用provisioned capacity模式。WCU/RCU值可以使用默认，后续根据压力测试规模调整。partionkey指定为“timestamp”，类型字符串。

1. 编译stream subscriber lambda

进入ddbreader/lambda目录，执行make。在aws界面上创建一个名为reader的lambda函数，使用x86 go runtime，设置128M内存，订阅timer表的stream。代码使用上面编译产生的结果reader.zip。

2. 编译负载生成程序
进入/ddbwriter目录，执行go build

生成程序后，参考下面参数进行压测步骤

```
./ddbwriter -h
Usage of ./ddbwriter:
  -d int
        Duration in minutes (default 60)
  -r int
        Number of goroutines (default 1)
  -t int
        Write interval in seconds (default 1)
```

3. 生成压测负载

启动10个go rountine生成压力，每秒钟写一次timer表，程序执行市场720分钟（12小时）。
```
./ddbwriter -d 720 -r 10 -t 1
```

4. 获取结果
查询dashboard表，查看time_diff_ms值。

---
### TODO LIST
- [x] 通过CDK自动部署压测环境
- [x] 支持codepipeline提交更新
- [ ] Pipeline和APP的Stack嵌套
- [ ] 使用eks部署压测客户端
- [ ] 使用cdk8s管理压测任务
- [ ] 自动删除lambda cloudwatch log group日志
- [ ] 自动过期DDB数据
- [ ] 压测数据可视化
- [ ] 支持kinesis data stream延迟测试



