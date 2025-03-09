package utils

// Đây là một ví dụ đơn giản, bạn có thể mở rộng logic để xếp hạng tay bài Poker
func EvaluateHand(holeCards, communityCards []string) string {
	// Placeholder: Xếp hạng tay bài (ví dụ: pair, straight, flush, v.v.)
	return "pair"
}

func DetermineWinner(hands []struct {
	PlayerID string
	Hand     string
}) string {
	// Placeholder: Xác định người thắng
	return hands[0].PlayerID
}
