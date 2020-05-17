package models

var (
	allModels = []interface{}{
		&Task{},
	}
)

// Models returns the slice of all models instances
func Models() []interface{} {
	return allModels
}
