package servicectrl

import "github.com/danenmao/pterergate-dtf/dtf/dtfdef"

// Service Start Function
type ServiceStartFn func(cfg *dtfdef.ServiceConfig) error

// 各服务role的操作表
var gs_ServiceRoleAction = map[dtfdef.ServiceRole]ServiceStartFn{
	dtfdef.ServiceRole_Manager:   StartManager,
	dtfdef.ServiceRole_Generator: StartGenerator,
	dtfdef.ServiceRole_Scheduler: StartScheduler,
	dtfdef.ServiceRole_Executor:  StartExecutor,
	dtfdef.ServiceRole_Collector: StartCollector,
}
