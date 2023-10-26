package quotagroup

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
	"github.com/danenmao/pterergate-dtf/internal/routine"
	"github.com/danenmao/pterergate-dtf/internal/taskframework/taskflow/schedulerflow/schedulingqueue"
)

// 资源组结构
type QuotaGroup struct {
	ID          uint32                           // ID
	Name        string                           // 资源组名
	Quota       float32                          // 资源组的调度资源配额
	Description string                           // 资源组的描述
	InsertTime  uint64                           // 资源组的插入时间
	QueueGroup  *schedulingqueue.SchedulingGroup // 资源组的调度队列组
}

// 资源组的fit值结构
type Quota struct {
	Name  string  // 资源组名
	Quota float32 // 资源组的配额值
}

// 资源组管理器
type QuotaGroupMgr struct {
	GroupMap      map[string]*QuotaGroup // 资源组表
	QuotaList     []Quota                // 所有资源组的quota值的列表
	MaxQuota      float32                // 资源组中最大的quota值
	MaxQuotaIndex int                    // 最大quota值元素的索引
	Mutex         sync.Mutex             // 访问锁
}

// 全局的资源组管理器对象
var gs_quotaGroupMgr = QuotaGroupMgr{
	GroupMap:  map[string]*QuotaGroup{},
	QuotaList: []Quota{},
}

// 获取模块的资源组管理器对象
func GetQuotaGroupMgr() *QuotaGroupMgr {
	return &gs_quotaGroupMgr
}

// 初始化资源组
func (rg *QuotaGroupMgr) Init() error {

	// 初始化管理器结构
	err := rg.initMgr()
	if err != nil {
		glog.Warning("failed to init mgr: ", err.Error())
		return err
	}

	// 启动同步例程
	go rg.syncRecordRoutine()

	// 创建调度监控例程
	go schedulingqueue.MonitorCurrentTaskRoutine()

	glog.Info("succeeded to init rg mgr")
	return nil
}

// 向资源组中添加任务
func (rg *QuotaGroupMgr) AddTask(
	groupName string,
	taskId taskmodel.TaskIdType,
	taskType uint32,
	priority uint32,
) error {

	rg.Mutex.Lock()
	defer rg.Mutex.Unlock()

	// 选择指定的资源组
	group, ok := rg.GroupMap[groupName]
	if !ok {
		glog.Warning("unknown quota group name: ", groupName)
		return errors.New("unknown quota group name")
	}

	// 向资源组中添加任务
	err := group.QueueGroup.AddTask(taskId, taskType, priority)
	if err != nil {
		glog.Warning("failed to add task to group: ", taskId, ", ", groupName)
		return err
	}

	glog.Info("succeeded to add task to group: ", taskId, ", ", groupName)
	return nil
}

// 选择要调度的资源组, 调度任务, 返回任务下被调度到的子任务列表
// 无法选出子任务时, retTaskId为0, subtasks返回的元素为空
func (rg *QuotaGroupMgr) Select(
	retTaskId *taskmodel.TaskIdType,
	subtasks *[]taskmodel.SubtaskBody,
) error {

	rg.Mutex.Lock()
	defer rg.Mutex.Unlock()

	// 选择资源组的索引
	i, err := rg.stochasticAccept()
	if err != nil {
		glog.Warning("failed to select a resource group: ", err)
		return err
	}

	// 根据索引取资源组
	group, ok := rg.GroupMap[rg.QuotaList[i].Name]
	if !ok {
		glog.Error("unknown quota group name: ", rg.QuotaList[i].Name)
		return err
	}

	// 从资源组的调度队列组中选择任务
	err = group.QueueGroup.Schedule(retTaskId, subtasks)
	if err != nil {
		glog.Warning("failed to schedule tasks from queue array: ", err.Error())
		return err
	}

	if *retTaskId != 0 {
		glog.Info("succeeded to schedule tasks from queue array: ", *retTaskId, len(*subtasks))
	}

	return nil
}

// 获取调度中的任务总数
func (rg *QuotaGroupMgr) GetTaskCount() (taskCount uint, err error) {

	rg.Mutex.Lock()
	defer rg.Mutex.Unlock()

	taskCount = 0
	count := uint(0)
	for _, group := range rg.GroupMap {
		count, err = group.QueueGroup.GetTaskCount()
		if err != nil {
			return
		}

		taskCount += count
	}

	return
}

// 初始化管理器结构
func (rg *QuotaGroupMgr) initMgr() error {

	// 读取资源组记录
	err := rg.syncRecord()
	if err != nil {
		glog.Warning("failed to sync rg record: ", err.Error())
		return err
	}

	glog.Info("succeeded to init rg mgr")
	return nil
}

// 基于随机接受（Stochastic Acceptance）的算法来实现资源组之间的配额调度
func (rg *QuotaGroupMgr) stochasticAccept() (int, error) {

	i := 0
	const MaxTryCount = 20
	n := len(rg.QuotaList)
	rand.Seed(time.Now().Unix())

	for {
		// 选择一个随机的元素
		idx := rand.Intn(n)

		// 按概率 Wi/Wmax 来接受选择
		if rand.Float32() <= rg.QuotaList[idx].Quota/rg.MaxQuota {
			return idx, nil
		}

		// 当超过最大次数后，未得到元素，选择当前随机的元素
		i++
		if i >= MaxTryCount {
			glog.Error("too many count to select an idx")
			return idx, nil
		}
	}
}

// 同步资源组的记录的例程
func (rg *QuotaGroupMgr) syncRecordRoutine() error {

	routine.ExecRoutineWithInterval(
		"syncRecordRoutine",
		func() {
			rg.syncRecord()
		},
		time.Duration(QuotaGroupSyncInterval)*time.Second,
	)

	return nil
}

// 同步资源组的记录
func (rg *QuotaGroupMgr) syncRecord() error {

	// 读取资源组记录
	records := []QuotaGroupRecord{}
	err := readQuotaGroupRecord(&records)
	if err != nil {
		glog.Warning("failed to read quota group records: ", err.Error())
		return err
	}

	rg.Mutex.Lock()
	defer rg.Mutex.Unlock()

	// 创建资源组结构
	for i, record := range records {

		// 创建或更新资源组结构
		group, err := rg.initOrUpdateGroup(&record)
		if err != nil {
			glog.Warning("failed to init quota group: ", record.Name, ",", err)
			continue
		}

		// 添加到fit数组尾部
		rg.QuotaList = append(rg.QuotaList, Quota{
			Name:  group.Name,
			Quota: group.Quota,
		})

		// 记录最大的fit值及索引
		if group.Quota > rg.MaxQuota {
			rg.MaxQuotaIndex = i
			rg.MaxQuota = group.Quota
		}
	}

	glog.Info("succeeded to sync rg record")
	return nil
}

// 从数据库中读取资源组的记录
func readQuotaGroupRecord(
	records *[]QuotaGroupRecord,
) error {

	predefinedRG := []QuotaGroupRecord{
		{
			ID: 2, Name: " 1", Quota: 0.6,
		},
		{
			ID: 3, Name: "2", Quota: 0.4,
		},
	}

	*records = append(*records, predefinedRG...)

	return nil
}

// 初始化或更新资源组记录
func (rg *QuotaGroupMgr) initOrUpdateGroup(
	record *QuotaGroupRecord,
) (*QuotaGroup, error) {

	groupName := record.Name
	groupQuota := record.Quota

	group, ok := rg.GroupMap[groupName]
	if ok {
		// 若资源组已存在, 仅更新配置
		group.Quota = groupQuota
		rg.GroupMap[groupName] = group
		return group, nil
	}

	// 若为新资源组，创建记录
	group = &QuotaGroup{
		Name:        groupName,
		ID:          record.ID,
		Quota:       groupQuota,
		Description: record.Description,
		InsertTime:  uint64(time.Now().Unix()),
		QueueGroup:  &schedulingqueue.SchedulingGroup{},
	}

	// 初始化调度队列组
	err := group.QueueGroup.Init(group.Name)
	if err != nil {
		glog.Warning("failed to init schedule queue array: ", record)
		return nil, err
	}

	// 记录到资源组map中
	rg.GroupMap[groupName] = group

	glog.Info("succeeded to init schedule queue array: ", record)
	return group, nil
}
