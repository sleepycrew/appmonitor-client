package check

type CheckContainer interface {
	AddCheck(metadata Metadata, check Check)
	AddNestedCheck(parent *string, metadata Metadata, check Check)
}
