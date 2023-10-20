package taskplugin

// 任务插件注册结构
type TaskPluginRegistration struct {
	TaskType        uint32              // 任务类型
	PluginFactoryFn TaskPluginFactoryFn // 插件创建函数
	Name            string              // 插件类型名称
	Description     string              // 描述任务的类型及主要操作
}
