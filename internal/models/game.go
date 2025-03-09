package models

type Player struct {
	ID        string   `json:"id"`
	HoleCards []string `json:"hole_cards"`
	Chips     int      `json:"chips"`
	IsFolded  bool     `json:"is_folded"`
}

type GameState struct {
	TableID        string   `json:"table_id"`
	Players        []Player `json:"players"`
	Deck           []string `json:"deck"`
	CommunityCards []string `json:"community_cards"`
	Pot            int      `json:"pot"`
	CurrentTurn    string   `json:"current_turn"`
}
