package atd

import (
    // "fmt"
    // "time"
    "goatd/app/event"
    "goatd/app/model"
    // "goatd/app/dispatch"
)


/*
 * Global
 */
var instance *ATD

func GetInstance() *ATD {
    if instance == nil {
        instance = new(ATD)
        instance.init()
    }
    return instance
}

/*
 * ATD
 */
type ATD struct {
    busManager *event.BusManager
    store *model.Store
}

func (atd *ATD) init() {
    atd.busManager = event.NewBusManager()
    atd.busManager.Start()
    atd.store = model.NewStore(atd.busManager)
}

func (atd *ATD) Run() {
}

func (atd *ATD) BusManager() *event.BusManager {
    return atd.busManager
}

func (atd *ATD) Store() *model.Store {
    return atd.store
}