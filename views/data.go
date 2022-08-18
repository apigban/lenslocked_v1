package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// Catchall message
	AlertMsgGeneric = "Something went wont. Please try again. Contact us if the problem persists."
)

// Alert is used to rendned Bootstrap Alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data
// to come in
type Data struct {
	Alert *Alert
	Yield interface{}
}
