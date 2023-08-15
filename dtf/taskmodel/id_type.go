package taskmodel

import (
	"strconv"
)

//
// 定义任务ID的类型
//
type TaskIdType uint64

// Marshal
func (taskId TaskIdType) MarshalBinary() (data []byte, err error) {
	return []byte(strconv.FormatUint(uint64(taskId), 10)), nil
}

// Unmarshal
func (taskId *TaskIdType) UnmarshalBinary(data []byte) error {
	ret, err := strconv.ParseUint(string(data), 10, 64)
	*taskId = TaskIdType(ret)
	return err
}


//
// 定义子任务ID的类型
//
type SubtaskIdType uint64

// Marshal
func (id SubtaskIdType) MarshalBinary() (data []byte, err error) {
	return []byte(strconv.FormatUint(uint64(id), 10)), nil
}

// Unmarshal
func (id *SubtaskIdType) UnmarshalBinary(data []byte) error {
	ret, err := strconv.ParseUint(string(data), 10, 64)
	*id = SubtaskIdType(ret)
	return err
}
