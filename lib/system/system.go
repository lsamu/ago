package system

import (
    agoNet "github.com/lsamu/ago/lib/net"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
    "github.com/shirou/gopsutil/process"
    "time"
)

type (
    Host struct {
        Name       string `json:"name"`        //服务器名称
        OS         string `json:"os"`          //操作系统
        IP         string `json:"ip"`          //服务器IP
        FrameWork  string `json:"frame_work"`  //系统架构
        UpdateTime uint64 `json:"update_time"` //更新时间
        BootTime   uint64 `json:"boot_time"`   //重启时间
    }
    CPU struct {
        Nums      int32     `json:"nums"`       //核心数
        Cores     int       `json:"cores"`      //核心数
        UsageRate []float64 `json:"usage_rate"` //cpu使用率
    }
    Mem struct {
        Total       uint64  `json:"total"`
        Free        uint64  `json:"free"`
        UsedPercent float64 `json:"used_percent"`
    }
    Disk struct {
        Drive     string  `json:"drive"`      //盘符
        FS        string  `json:"fs"`         //文件系统
        Total     uint64  `json:"total"`      //总大小
        Usable    uint64  `json:"usable"`     //可用大小
        Usage     uint64  `json:"usage"`      //已用大小
        UsageRate float64 `json:"usage_rate"` //已用百分比
    }

    Net struct {
        Name      string `json:"name"`
        BytesSent uint64 `json:"bytes_sent"`
        BytesRecv uint64 `json:"bytes_recv"`
    }

    Process struct {
        Name       string  `json:"name"`
        UserName   string  `json:"user_name"`
        Pid        int32   `json:"pid"`
        CpuPercent float64 `json:"cpu_percent"`
        MemPercent int64   `json:"mem_percent"`
        CreateTime float32 `json:"create_time"`
    }

    Go struct {
        Lang      string //语言环境
        Path      string //安装路径
        StartTime string //启动时间
        Version   string //版本
        RunTime   string //运行时长
    }
)

func GetHostInfo() (info Host, err error) {
    //主机
    hostInfo, err := host.Info()
    if err != nil {
        return info, err
    }
    info.Name = hostInfo.Hostname
    info.OS = hostInfo.OS
    info.FrameWork = hostInfo.Platform
    info.IP, _ = agoNet.GetLocalIP()
    info.UpdateTime = hostInfo.Uptime
    info.BootTime = hostInfo.BootTime
    return
}

func GetCPUInfo() (info CPU, err error) {
    cpuInfo, err := cpu.Info()
    if err != nil {
        return info, err
    }
    info.Nums = cpuInfo[0].Cores
    info.UsageRate, _ = cpu.Percent(time.Second, false)
    info.Cores, _ = cpu.Counts(true)
    return
}

func GetMemInfo() (info Mem, err error) {
    memInfo, err := mem.VirtualMemory()
    if err != nil {
        return info, err
    }
    info.Total = memInfo.Total
    info.Free = memInfo.Free
    info.UsedPercent = memInfo.UsedPercent
    return
}

func GetNetInfo() (info []Net, err error) {
    netInfo, err := net.IOCounters(true)
    if err != nil {
        return info, err
    }
    for _, v := range netInfo {
        if v.Name != "enp3s0" {
            continue
        }
        info = append(info, Net{
            Name:      v.Name,
            BytesSent: v.BytesSent,
            BytesRecv: v.BytesRecv,
        })
    }
    return
}

func GetDiskInfo() (info []Disk, err error) {
    parts, err := disk.Partitions(true)
    if err != nil {
        return info, err
    }
    for _, part := range parts {
        if part.Fstype != "ext4" {
            continue
        }
        diskInfo, _ := disk.Usage(part.Mountpoint)
        info = append(info, Disk{
            Drive:     part.Device,
            FS:        part.Fstype,
            Total:     diskInfo.Total,
            Usable:    diskInfo.Free,
            Usage:     diskInfo.Used,
            UsageRate: diskInfo.UsedPercent,
        })
    }
    return
}

func GetProcessInfo(filterUserName string) (info []Process, err error) {
    infProcess, err := process.Processes()
    if err != nil {
        return info, err
    }
    for _, pro := range infProcess {
        isBackground, _ := pro.Background()
        if isBackground {
            continue
        }
        userName, _ := pro.Username()
        if filterUserName != "" && userName == filterUserName {
            continue
        }
        name, _ := pro.Name()
        pid := pro.Pid
        cpuPer, _ := pro.CPUPercent()
        createTime, _ := pro.CreateTime()
        memPer, _ := pro.MemoryPercent()
        info = append(info, Process{
            Name:       name,
            UserName:   userName,
            Pid:        pid,
            CpuPercent: cpuPer,
            MemPercent: createTime,
            CreateTime: memPer,
        })
    }
    return
}

func GetGoInfo() (info Go, err error) {

    return
}