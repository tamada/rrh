package decorator

type noDecorator struct {
}

func (nd *noDecorator) RepositoryID(repoID string) string {
	return repoID
}
func (nd *noDecorator) GroupName(name string) string {
	return name
}
func (nd *noDecorator) EnvironmentValue(value string) string {
	return value
}
func (nd *noDecorator) EnvironmentLabel(label string) string {
	return label
}
