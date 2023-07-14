package main

import "errors"

var (
	ErrorNoPathSelect                 = errors.New("no Path selected")
	ErrorEmptyPath                    = errors.New("no Key Entered")
	ErrorBadPathPrefix                = errors.New("path should start with '/'")
	ErrorPathHasNoData                = errors.New("path has No Data")
	ErrorPathAlreadyHasData           = errors.New("path already has data")
	ErrorSelectEtcdConnectionNotFound = errors.New("select Etcd Connection not found")
	ErrorDuplicateEtcdConnectionName  = errors.New("duplicate Etcd Connection Name")
	ErrorInputDataEmpty               = errors.New("input Data Can't Be Empty")
	ErrorConnNameAlreadyExist         = errors.New("conn name already exist")
	ErrorAppResetDone                 = errors.New("app reset done")
	ErrorDeleteSelectingConnection    = errors.New("can not delete selecting connection")
	ErrorConnectToEtcdFailed          = errors.New("connect to ETCD failed")
	ErrorConnectToEtcdFailedInfo      = errors.New("check ETCD Connection config or dial-timeout")
	ErrorConnectionNameEmpty          = errors.New("connection name is empty")
	ErrorInputNeedJsonFormat          = errors.New("input data need json format")
)
