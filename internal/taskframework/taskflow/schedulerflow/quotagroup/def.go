package quotagroup

const (

	// 资源组同步的间隔, 秒
	QuotaGroupSyncInterval uint32 = 120
)

// 资源组记录定义
type QuotaGroupRecord struct {
	ID          uint32
	Name        string
	Quota       float32
	Description string
	InsertTime  string
}
