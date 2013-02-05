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

func (team *Team) Update(name string, value interface{}) bool {
    setAttributeValue(team, name, value)
    return team.store.Update(KindTeam, team.Uid(), name, value)
}

func (team *Team) Destroy() *Team {
    if destroyed := team.store.Destroy(KindTeam, team.Uid()); destroyed != nil {
        return destroyed.(*Team)
    }
    return nil
}

func (team Team) Reload() *Team {
    if found := team.store.Teams.Find(team.Uid()); found != nil {
        return found
    }
    return nil
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

func (proxy *TeamStoreProxy) Each(iterator func(interface{})) {
    proxy.store.Each(KindTeam, iterator)
}

func (proxy *TeamStoreProxy) Count() int {
    return proxy.store.Count(KindTeam)
}

func (proxy *TeamStoreProxy) DestroyAll() {
    proxy.store.DestroyAll(KindTeam)
}
