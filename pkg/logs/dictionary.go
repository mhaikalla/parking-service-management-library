package logs

import "os"

var (
	// LogLevel set log level eg: `trace`, `debug`, etc. Default: `error`.
	LogLevel = os.Getenv("LOG_LEVEL")

	// LogFormat set log format like `text` and `json`. Default: `text`,
	LogFormat = os.Getenv("LOG_FORMAT")

	// ServiceName this service name that run this log.
	ServiceName = os.Getenv("SERVICE_NAME")
)

const (
	// LogType entry log type.
	LogType = "log_type"

	// LogLocation entry log location.
	LogLocation = "log_location"

	// LogName entry log name.
	LogName = "log_name"

	// LogID entry for log ID.
	LogID = "log_id"

	// LogParent entry for log parent id.
	LogParent = "log_parent"

	// LogChilds entry for log childs.
	LogChilds = "log_childs"

	// LogCreatedAt time when this log created.
	LogCreatedAt = "log_created_at"

	// ClientOriginalPath path client request.
	ClientOriginalPath = "client_request_path"

	// ClientRequestID client request ID for this logger.
	ClientRequestID = "client_request_id"

	// ClientDeviceID client request
	ClientDeviceID = "client_device_id"

	// ClientSubsID client subs ID to identified user.
	ClientSubsID = "client_subs_id"

	// ClientSubsType client subs type.
	ClientSubsType = "client_subs_type"

	// ClientBearerType bearer type user have
	ClientBearerType = "client_bearer_type"

	// ClientBearerContent content of client bearer token.
	ClientBearerContent = "client_bearer_content"

	// ClientRequestIP client request IP.
	ClientRequestIP = "client_request_ip"

	// ClientUserAgent client user agent when requesting.
	ClientUserAgent = "client_user_agent"

	// SystemCommitSHA system source code commit SHA.
	SystemCommitSHA = "sys_commit_sha"

	// SystemBuildTime this system binary build.
	SystemBuildTime = "sys_build_time"

	// SystemBranch branch this system build.
	SystemBranch = "sys_branch"

	// SystemName the name of this service.
	SystemName = "sys_name"
)
