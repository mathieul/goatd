package event

const (
    KindNone Kind = iota
    OfferTask
    AcceptTask
    CompleteTask
    TeammateAvailable
    SkillCreated
    KindLast
)


/*
 * Basic types
 */
type Kind int


/*
 * initialization
 */

var allKinds []Kind

func init() {
    allKinds = make([]Kind, 0, KindLast - KindNone)
    for i := KindNone + 1; i < KindLast; i++ {
        allKinds = append(allKinds, i)
    }
}
