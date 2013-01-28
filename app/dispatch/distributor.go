package dispatch

import (
    "goatd/app/model"
)


/*
 * Distributor
 */

type Distributor struct {
    *model.Store
}

func NewDistributor(store *model.Store) *Distributor {
    return &Distributor{store}
}