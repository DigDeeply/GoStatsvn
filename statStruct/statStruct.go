package statStruct

type AuthorStat struct {
	AppendLines int
	RemoveLines int
}

type AuthorTimeStat map[string]AuthorStat	//map[time]AuthorStat
type AuthorTimeStats map[string][]AuthorTimeStat //map[user]AuthorTimeStat
