package truenasapi

/*
NotFoundError is returned when the desired resource cannot be located
*/
type NotFoundError struct {
	Err error
}

func (e *NotFoundError) Error() string {
	return e.Err.Error()
}

/*
InternalError is returned when something on our end
or the api's is wrong.
*/
type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return e.Err.Error()
}

/*
 AlreadyExistsError is returned when an api
 resource has already be created, ie a volume with
 that same name already exists.
*/
type AlreadyExistsError struct {
	Err error
}

func (e *AlreadyExistsError) Error() string {
	return e.Err.Error()
}

/*
 ResourceBusyError is returned when something cannot be done to a resource
 ie. Delete a busy volume or a volume with children
*/
type ResourceBusyError struct {
	Err error
}

func (e *ResourceBusyError) Error() string {
	return e.Err.Error()
}
