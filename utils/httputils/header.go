package httputils

var headers = [...]string{
	HTTPHeaderContext,
	HTTPHeaderTrackingID,
	HTTPHeaderModule,
}

const (
	HTTPHeaderTrackingID = "Umni-Tracking-ID"
	HTTPHeaderContext    = "Umni-Context"
	HTTPHeaderModule     = "Umni-Module"
)
