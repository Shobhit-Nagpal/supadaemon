package errors

import "fmt"

type FileError struct {
	baseError
}

func NewFileError(fileName, op string) FileError {
	errMsg := ""

	switch op {
	case WriteFile:
		errMsg = fmt.Sprintf("Failed to write file: %s\n", fileName)
	case ReadFile:
		errMsg = fmt.Sprintf("Failed to read file: %s\n", fileName)
	case DeleteFile:
		errMsg = fmt.Sprintf("Failed to delete file: %s\n", fileName)
	case FileNotFound:
		errMsg = fmt.Sprintf("File does not exist: %s\n", fileName)
	default:
		errMsg = fmt.Sprintf("Unexpected error for file: %s\n", fileName)
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
		errMsg = fmt.Sprintf("Failed to reload service\n")
	case StopService:
		errMsg = fmt.Sprintf("Failed to stop service\n")
	case StartService:
		errMsg = fmt.Sprintf("Failed to start service\n")
	case EnableService:
		errMsg = fmt.Sprintf("Failed to enable service\n")
	case DisableService:
		errMsg = fmt.Sprintf("Failed to disable service\n")
	case RootAccess:
		errMsg = fmt.Sprintf("Require root access (run with sudo)\n")
	default:
		errMsg = fmt.Sprintf("%s\n", service)
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
		errMsg = fmt.Sprintf("Failed to get current user\n")
	case ExecuteableLookup:
		errMsg = fmt.Sprintf("Failed to get executeable\n")
	case ReloadDaemon:
		errMsg = fmt.Sprintf("Failed to reload daemon\n")
	case DataSerialization:
		errMsg = fmt.Sprintf("Failed to serialize data\n")
	case DataDeserialization:
		errMsg = fmt.Sprintf("Failed to deserialize data\n")
	default:
		errMsg = fmt.Sprintf("Runtime error\n")
	}

	errMsg = buildErrorMsg(errMsg, SystemErr)
	return ServiceError{
		baseError: newBaseError(errMsg),
	}
}
