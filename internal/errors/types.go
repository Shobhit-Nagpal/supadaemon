package errors

import "fmt"

type FileError struct {
	baseError
}

func NewFileError(fileName, op string) FileError {
	errMsg := ""

	switch op {
	case WriteFile:
		errMsg = fmt.Sprintf("Failed to write file: %s", fileName)
	case ReadFile:
		errMsg = fmt.Sprintf("Failed to read file: %s", fileName)
	case DeleteFile:
		errMsg = fmt.Sprintf("Failed to delete file: %s", fileName)
	case FileNotFound:
		errMsg = fmt.Sprintf("File does not exist: %s", fileName)
	default:
		errMsg = fmt.Sprintf("Unexpected error for file: %s", fileName)
	}

	errMsg = buildErrorMsg(errMsg, FileErr)
	return FileError{
		baseError: newBaseError(errMsg),
	}
}

type ServiceError struct {
	baseError
}

func NewServiceError(service, op string) ServiceError {
	errMsg := ""
	switch op {
	case ReloadService:
		errMsg = fmt.Sprintf("Failed to reload service")
	case StopService:
		errMsg = fmt.Sprintf("Failed to stop service")
	case StartService:
		errMsg = fmt.Sprintf("Failed to start service")
	case EnableService:
		errMsg = fmt.Sprintf("Failed to enable service")
	case DisableService:
		errMsg = fmt.Sprintf("Failed to disable service")
	case RootAccess:
		errMsg = fmt.Sprintf("Require root access (run with sudo)")
	case GetServiceStatus:
		errMsg = fmt.Sprintf("Failed to get service status")
	default:
		errMsg = fmt.Sprintf("%s", service)
	}

	errMsg = buildErrorMsg(errMsg, ServiceErr)
	return ServiceError{
		baseError: newBaseError(errMsg),
	}
}

type SystemError struct {
	baseError
}

func NewSystemError(op string) ServiceError {
	errMsg := ""
	switch op {
	case UserLookup:
		errMsg = fmt.Sprintf("Failed to get current user")
	case ExecuteableLookup:
		errMsg = fmt.Sprintf("Failed to get executeable")
	case ReloadDaemon:
		errMsg = fmt.Sprintf("Failed to reload daemon")
	case DataSerialization:
		errMsg = fmt.Sprintf("Failed to serialize data")
	case DataDeserialization:
		errMsg = fmt.Sprintf("Failed to deserialize data")
	default:
		errMsg = fmt.Sprintf("Runtime error")
	}

	errMsg = buildErrorMsg(errMsg, SystemErr)
	return ServiceError{
		baseError: newBaseError(errMsg),
	}
}
