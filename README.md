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

1. Redis

## Usage

1. Create the task plugin.

1. Register the task type.

1. Invoke the dtf services.

1. Create a task.

1. Wait for the task to complete.
