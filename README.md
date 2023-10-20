# pterergate-dtf

Pterergate-dtf (Pterergate Distributed Task Framework) is a high-performance distributed task framework that supports parallelly scheduling thousands 
of running tasks deployed in a cluster consisting of tens thousands of nodes.

[![Go](https://github.com/danenmao/pterergate-dtf/actions/workflows/go.yml/badge.svg)](https://github.com/danenmao/pterergate-dtf/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/danenmao/pterergate-dtf?status.svg)](https://godoc.org/github.com/danenmao/pterergate-dtf)
[![Go Report Card](https://goreportcard.com/badge/github.com/danenmao/pterergate-dtf)](https://goreportcard.com/report/github.com/danenmao/pterergate-dtf)
![GitHub](https://img.shields.io/github/license/danenmao/pterergate-dtf)

![GitHub last commit (branch)](https://img.shields.io/github/last-commit/danenmao/pterergate-dtf/main)
![GitHub commit activity (branch)](https://img.shields.io/github/commit-activity/t/danenmao/pterergate-dtf)

## Install

```console
go get github.com/danenmao/pterergate-dtf
```

## Requirement

1. MySQL

    DTF uses a MySQL table `tbl_task` to store the information of created tasks. Users shoule provide a MySQL server, and create this table in a database.

2. Redis

   DTF uses Redis frequently to store kinds of intermediate data.Users should provide a Redis server.

## Usage

1. Implement ITaskGenerator, ITaskExecutor, ITaskSchedulerCallback and ITaskCollectorCallback. Users can perform their business logic in these interface.

    ```Go
    // implement taskmodel.ITaskGenerator
    type SampleGenerator struct{}
    
    // implement taskmodel.ITaskExecutor
    type SampleExecutor struct{}

    // implement taskmodel.ITaskSchedulerCallback
    type SampleSchedulerCallback struct{}

    // implement taskmodel.ITaskCollectorCallback
    type SampleCollectorCallback struct{}
    ```

2. Implement a task plugin.

    ``` Go
    // implement taskplugin.ITaskPlugin
    type SamplePlugin struct{
        TaskBody taskmodel.TaskBody
        TaskConf taskmodel.TaskConf
    }

    func (p * SamplePlugin) GetTaskConf(taskConf *taskmodel.TaskConf) error{
        *taskConf = p.TaskConf
        return nil
    }

    func  (p * SamplePlugin) GetTaskBody(taskBody *taskmodel.TaskBody) error{
        *taskBody = p.TaskBody
        return nil
    }
    
    var plugin = SamplePlugin{
        TaskBody: taskmodel.TaskBody{
            Generator: &SampleGenerator{},
            Executor: &SampleExecutor{},
            SchedulerCallback: &SampleSchedulerCallback{},
            CollectorCallback: &SampleCollectorCallback{},
        },
        TaskConf: taskmodel.TaskConf{
            IterationMode: taskmodel.IterationMode_No,
            TaskTypeTimeout: time.Hour,
        },
    }
    ```

3. Register the task type.

    ``` Go
    const SampleTaskType = 1
    register := taskplugin.TaskPluginRegistration{
        TaskType: SampleTaskType,
        Name: "SampleTaskType",
        Description: "a sample task type",
        PluginFactoryFn: func(p *ITaskPlugin) error{
            *p = &plugin
        }
    }

    err := dtf.RegisterTaskType(&register)
    ```

4. Invoke the dtf services.

    ``` Go
    err := dtf.StartService(
        dtfdef.ServiceRole_Manager, 
        dtf.WithMySQL(...),
        dtf.WithRedis(...),
    )
    ```

5. Create a task to perform some operation.

    ``` Go
    taskParam := taskmodel.TaskParam{
        ...
    }

    taskId, err := dtf.CreateTask(
        SampleTaskType,
        taskParam,
    )
    ```

6. Wait for the service to exit.

    ``` Go
    dtf.Join()
    ```
