package filter

import (
	"net/http"

	// "github.com/JessonChan/canlog"

	"github.com/JessonChan/cango"
)

type VisitFilter struct {
	cango.Filter
	// controller
	// *CanCtrl
}

var _ = cango.RegisterFilter(&VisitFilter{})

func (v *VisitFilter) PreHandle(req *http.Request) interface{} {
	// canlog.CanDebug(req.Method, req.URL.Path)
	return true
}
