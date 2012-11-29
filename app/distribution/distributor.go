package distribution

import (
    "goatd/app/models"
)

/*
 * Distributor
 */
type Distributor struct {
    team *models.Team
}

func (distributor Distributor) Team() *models.Team {
    return distributor.team
}

func NewDistributor(team *models.Team) *Distributor {
    return &Distributor{team}
}
