package atd

import (
    "time"
)

type OverviewService struct {}

func (service OverviewService) Index(request []byte) (response []byte) {
    message := string(request)
    return "[OverviewService] " + message + "[" + time.Now().String() + "]"
}
