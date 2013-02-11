package atd

import (
    "time"
)

type TeamsService struct {}

func (service TeamsService) Index(request []byte) (response []byte) {
    message := string(request)
    return "[TeamsService] " + message + "[" + time.Now().String() + "]"
}
