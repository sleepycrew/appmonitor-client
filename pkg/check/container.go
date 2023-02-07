package check

type CheckContainer interface {
	AddCheck(check Check)
	AddNestedCheck(parent *string, check Check)
}
