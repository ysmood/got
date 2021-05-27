package diff

// TokenHunk is a pair of Del
type TokenHunk struct {
	Del []*TokenLine
	Add []*TokenLine
}

func (h *TokenHunk) From() int {
	if len(h.Del) > 0 {
		return h.Del[0].N
	}
	return h.Add[0].N
}

func (h *TokenHunk) To() int {
	if len(h.Add) > 0 {
		return h.Add[len(h.Add)-1].N
	}
	return h.Del[len(h.Del)-1].N
}

// TokenLine of tokens
type TokenLine struct {
	N      int
	Type   Type
	Symbol *Token
	Tokens []*Token
}

// ParseTokenLines of tokens
func ParseTokenLines(ts []*Token) []*TokenLine {
	list := []*TokenLine{}
	var buf *TokenLine
	for i, t := range ts {
		switch t.Type {
		case LineNum:
			if buf != nil {
				list = append(list, buf)
			}
			buf = &TokenLine{N: i + 1}
		case SameSymbol, AddSymbol, DelSymbol:
			buf.Type = t.Type
			buf.Symbol = t
		default:
			buf.Contents = append(buf.Contents, t)
		}
	}
	list = append(list, buf)
	return list
}

// ParseTokenHunks of changes
func ParseTokenHunks(lines []*TokenLine) []*TokenHunk {
	list := []*TokenHunk{}
	var buf *TokenHunk
	for _, l := range lines {
		switch l.Type {
		case DelSymbol:
			if buf == nil {
				buf = &TokenHunk{Del: []*TokenLine{}, Add: []*TokenLine{}}
			}
			buf.Del = append(buf.Del, l)
		case AddSymbol:
			if buf == nil {
				buf = &TokenHunk{Del: []*TokenLine{}, Add: []*TokenLine{}}
			}
			buf.Add = append(buf.Add, l)
		case SameSymbol:
			if buf != nil {
				list = append(list, buf)
			}
			buf = nil
		}
	}
	if buf != nil {
		list = append(list, buf)
	}
	return list
}
