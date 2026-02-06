package app

import "fmt"

type GroupNotActiveError struct {
    GroupName string
}
	
func (e *GroupNotActiveError) Error() string {
    return fmt.Sprintf("dynamic group '%s' is not active", e.GroupName)
}