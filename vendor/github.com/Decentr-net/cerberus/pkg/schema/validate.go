package schema

import (
	"fmt"

	valid "github.com/asaskevich/govalidator"
)

// Validate ...
type Validate interface {
	Validate() bool
}

// Validate ...
func (o *PDVObjectV1) Validate() bool {
	if !valid.IsURL(fmt.Sprintf("%s/%s", o.Host, o.Path)) {
		return false
	}
	return len(o.Data) > 0
}

// Validate ...
func (d *PDVDataCookieV1) Validate() bool {
	if d.Name == "" || d.Value == "" {
		return false
	}

	return true
}
