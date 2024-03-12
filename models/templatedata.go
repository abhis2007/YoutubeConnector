package models

type TemplateData struct {
	StringMap  map[string]string
	ArrMap     map[string]int
	FloatMap   map[string]float32
	CustomData map[string]interface{}
	StringData string
	CSRFToken  string
}
