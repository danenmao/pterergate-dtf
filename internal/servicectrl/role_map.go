package servicectrl

import "pterergate-dtf/dtf/dtfdef"

type ServiceStartFn func(cfg *dtfdef.ServiceConfig) error

// 各服务role的操作表
var gs_ServiceRoleAction = map[dtfdef.ServiceRole]ServiceStartFn{
	dtfdef.ServiceRole_Manager:   StartManager,
	dtfdef.ServiceRole_Generator: nil,
	dtfdef.ServiceRole_Scheduler: nil,
	dtfdef.ServiceRole_Executor:  nil,
	dtfdef.ServiceRole_Iterator:  nil,
}
