package taskplugin


//
// 任务插件注册结构
//
type TaskPluginRegister struct {
	TaskType    	uint32              // 任务类型
	NewPluginFn 	NewTaskPluginFn 	// 插件创建函数
	Description 	string				// 描述任务的类型及主要操作
}




