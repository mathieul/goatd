package atd

type OverviewService struct {}

/*
 * Overview.List
 */

type OverviewIndexReq struct {}

type OverviewRow struct {
    Name string `json:"name"`
    Count int   `json:"count"`
}

type OverviewIndexRes struct {
    Rows []OverviewRow `json:"rows"`
}

func (service OverviewService) List(req OverviewIndexReq) OverviewIndexRes {
    res := new(OverviewIndexRes)
    res.Rows = []OverviewRow{
        {"Team",     2},
        {"Teammate", 5},
        {"Queue",    3},
        {"Skill",   10},
        {"Task",    42},
    }

    return *res
}
