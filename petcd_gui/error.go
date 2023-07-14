package main

import "errors"

var (
	ErrorNoPathSelect               = errors.New("no Path selected")
	ErrorEmptyPath                  = errors.New("no Key Entered")
	ErrorBadPathPrefix              = errors.New("path should start with '/'")
	ErrorPathHasNoData              = errors.New("path has No Data")
	ErrorSelectEtcdEndPointNotFound = errors.New("select Etcd EndPoint not found")
	ErrorDuplicateEtcdEndPointName  = errors.New("duplicate Etcd EndPoint Name")
	ErrorInputDataEmpty             = errors.New("input Data Can't Be Empty")
	ErrorConnNameAlreadyExist       = errors.New("conn name already exist")
	ErrorAppResetDone               = errors.New("app reset done")
	ErrorDeleteSelectingEndPoint    = errors.New("can not delete selecting endpoint")
	ErrorConnectToEtcdFailed        = errors.New("connect to ETCD failed")
	ErrorConnectToEtcdFailedInfo    = errors.New("check ETCD endpoint config or dial-timeout")
	ErrorInvalidInput               = errors.New("invalid input")
)
