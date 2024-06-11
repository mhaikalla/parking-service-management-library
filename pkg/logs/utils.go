package logs

const (
	HandlerOps = "HANDLER OPS"
)

// WhenError when `err` not equal nil, log it using child.
func WhenError(name string, err error, logger ILog) {
	if err != nil && logger != nil {
		logger.Child(name).Error(err)
	}
}

// HandlerWhenError like `WhenError` but name set to `HandlerOps`.
func HandlerWhenError(err error, logger ILog) {
	WhenError(HandlerOps, err, logger)
}
