package model

import (
    "goatd/app/event"
    "goatd/app/sm"
)

/*
 * Team
 */

type Team struct {
    *event.Identity
    busManager *event.BusManager
    store *Store
    AttrName string
}

func (team Team) Name() string {
    return team.AttrName
}

func (team Team) Status(_ ...sm.Status) sm.Status {
    return StatusNone
}

func (team *Team) Copy() Model {
    return &Team{team.Identity, team.busManager, team.store, team.AttrName}
}

func (team *Team) SetupComs(busManager *event.BusManager, store *Store) {
    team.busManager = busManager
    team.store = store
}

func NewTeam(attributes A) (team *Team) {
    team = newModel(&Team{}, &attributes).(*Team)
    team.Identity = event.NewIdentity("Team")
    return team
}


/*
 * TeamStoreProxy
 */

type TeamStoreProxy struct {
    store *Store
}

func toTeamSlice(source []Model) []*Team {
    teams := make([]*Team, 0, len(source))
    for _, team := range source {
        teams = append(teams, team.(*Team))
    }
    return teams
}

func (proxy *TeamStoreProxy) Create(attributes A) *Team {
    return proxy.store.Create(KindTeam, attributes).(*Team)
}

func (proxy *TeamStoreProxy) Find(uid string) *Team {
    if value := proxy.store.Find(KindTeam, uid); value != nil { return value.(*Team) }
    return nil
}

func (proxy *TeamStoreProxy) FindAll(uids []string) []*Team {
    values := proxy.store.FindAll(KindTeam, uids)
    return toTeamSlice(values)
}

func (proxy *TeamStoreProxy) Select(tester func(interface{}) bool) []*Team {
    values := proxy.store.Select(KindTeam, tester)
    return toTeamSlice(values)
}
