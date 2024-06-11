package contexts

import (
	"context"

	"github.com/mhaikalla/parking-service-management-library/pkg/logs"
)

var (
	requestContextBaggage = new(string)
)

// Baggage contains value related to user that will be injected to `context.Context`.
type Baggage struct {
	DeviceID    string
	SubsID      string
	Substype    string
	RequestID   string
	PathOrigin  string
	Logger      logs.ILog
	BearerToken string
}

// GetBaggage get baggage from `context.Context`,
// If `ILog` exist, create child of it.
func GetBaggage(ctx context.Context) Baggage {
	defaultBag := Baggage{Logger: logs.NewLogrus("CONTEXT_LOG")}
	if ctx == nil {
		return defaultBag
	}
	value := ctx.Value(requestContextBaggage)
	if baggage, ok := value.(Baggage); ok {
		if baggage.Logger == nil {
			baggage.Logger = logs.NewLogrus("CONTEXT_LOG")
			return baggage
		}
		baggage.Logger = baggage.Logger.Child("")
		return baggage
	}

	return defaultBag
}

// InjectBaggage inject baggage to `context.Context`.
func InjectBaggage(ctx context.Context, baggage Baggage) context.Context {
	if ctx == nil {
		return context.WithValue(context.TODO(), requestContextBaggage, baggage)
	}
	return context.WithValue(ctx, requestContextBaggage, baggage)
}

// CopyBaggage copy baggage from `src` to `dst`.
func CopyBaggage(src context.Context, dst context.Context) context.Context {
	bag := GetBaggage(src)
	return InjectBaggage(dst, bag)
}
