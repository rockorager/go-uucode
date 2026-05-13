package uucode

// LineBreakKind describes the line break opportunity after a line segment.
type LineBreakKind uint8

// Line break opportunity values.
const (
	// LineDontBreak means a line must not be broken at this boundary.
	LineDontBreak LineBreakKind = iota
	// LineCanBreak means a line may be broken at this boundary.
	LineCanBreak
	// LineMustBreak means a line must be broken at this boundary.
	LineMustBreak
)

// String returns the line break opportunity name for kind.
func (kind LineBreakKind) String() string {
	switch kind {
	case LineCanBreak:
		return "can_break"
	case LineMustBreak:
		return "must_break"
	default:
		return "dont_break"
	}
}

// LineSegment identifies a non-breaking line segment by byte offsets into the
// original string. Break describes the line break opportunity after End.
type LineSegment struct {
	// Start is the byte offset of the first byte in the segment.
	Start int
	// End is the byte offset just after the segment.
	End int
	// Break is the line break opportunity after End.
	Break LineBreakKind
}

// LineIterator iterates over non-breaking line segments in a string according
// to the Unicode Line Breaking Algorithm (UAX #14).
type LineIterator struct {
	s        string
	segStart int
	next     int
	lastEnd  int
	done     bool
	state    lineState
}

// NewLineIterator returns a line break iterator for s.
func NewLineIterator(s string) LineIterator {
	if len(s) == 0 {
		return LineIterator{s: s, done: true}
	}
	item, ok := nextLineItem(s, 0)
	if !ok {
		return LineIterator{s: s, done: true}
	}
	var state lineState
	state.consume(item, true)
	return LineIterator{s: s, next: item.end, lastEnd: item.end, state: state}
}

// Next returns the next non-breaking line segment.
//
// The returned LineSegment contains byte offsets into the original string. ok
// is false after the iterator is exhausted. In accordance with UAX #14 LB3,
// the final segment is returned with Break set to LineMustBreak.
func (it *LineIterator) Next() (LineSegment, bool) {
	if it.done || !it.state.havePrev {
		return LineSegment{}, false
	}
	for {
		if it.next >= len(it.s) {
			it.done = true
			return LineSegment{Start: it.segStart, End: it.lastEnd, Break: LineMustBreak}, true
		}

		item, ok := nextLineItem(it.s, it.next)
		if !ok {
			it.done = true
			return LineSegment{Start: it.segStart, End: it.lastEnd, Break: LineMustBreak}, true
		}
		kind := it.state.boundaryBefore(item, it.s[item.end:])
		it.state.consume(item, kind == LineDontBreak)
		it.next = item.end
		it.lastEnd = item.end
		if kind != LineDontBreak {
			seg := LineSegment{Start: it.segStart, End: item.start, Break: kind}
			it.segStart = item.start
			return seg, true
		}
	}
}

// Peek returns the next line segment without advancing the iterator.
func (it LineIterator) Peek() (LineSegment, bool) {
	return (&it).Next()
}

type lineItem struct {
	start    int
	end      int
	lb       LineBreakClass
	gc       GeneralCategoryClass
	eaw      EastAsianWidthClass
	zwj      bool
	cm       bool
	dotted   bool
	extPicCn bool
}

type lineState struct {
	havePrev  bool
	prev      lineItem
	prev2     lineItem
	havePrev2 bool

	haveBeforeSpaces bool
	beforeSpaces     lineItem

	prevInitialQuote         bool
	beforeSpacesInitialQuote bool

	afterZWJ bool
	riRun    int
	numState uint8
}

const (
	lineNumNone uint8 = iota
	lineNumSeq
	lineNumClose
)

func nextLineItem(s string, i int) (lineItem, bool) {
	r, end, ok := nextRuneInString(s, i)
	if !ok {
		return lineItem{}, false
	}
	row := runtimeLookup(r)
	lb := row.lineBreak()
	gc := row.generalCategory()

	// LB1. Resolve classes whose line breaking behavior depends on external
	// criteria to the default classes specified by UAX #14.
	switch lb {
	case LineBreakAI, LineBreakSG, LineBreakXX:
		lb = LineBreakAL
	case LineBreakSA:
		if gc == GeneralCategoryMn || gc == GeneralCategoryMc {
			lb = LineBreakCM
		} else {
			lb = LineBreakAL
		}
	case LineBreakCJ:
		lb = LineBreakNS
	}

	return lineItem{
		start:    i,
		end:      end,
		lb:       lb,
		gc:       gc,
		eaw:      row.eastAsianWidth(),
		zwj:      lb == LineBreakZWJ,
		cm:       lb == LineBreakCM,
		dotted:   r == 0x25cc,
		extPicCn: row.isExtendedPictographic() && gc == GeneralCategoryCn,
	}, true
}

func (st *lineState) boundaryBefore(cur lineItem, rest string) LineBreakKind {
	if !st.havePrev {
		return LineDontBreak // LB2.
	}
	prev := st.prev

	// LB4 and LB5.
	if prev.lb == LineBreakBK {
		return LineMustBreak
	}
	if prev.lb == LineBreakCR {
		if cur.lb == LineBreakLF {
			return LineDontBreak
		}
		return LineMustBreak
	}
	if prev.lb == LineBreakLF || prev.lb == LineBreakNL {
		return LineMustBreak
	}

	// LB6.
	if isLineHardBreak(cur.lb) {
		return LineDontBreak
	}

	// LB7.
	if cur.lb == LineBreakSP || cur.lb == LineBreakZW {
		return LineDontBreak
	}

	// LB8.
	if prev.lb == LineBreakZW || (prev.lb == LineBreakSP && st.haveBeforeSpaces && st.beforeSpaces.lb == LineBreakZW) {
		return LineCanBreak
	}

	// LB8a.
	if st.afterZWJ {
		return LineDontBreak
	}

	// LB9. Combining marks and ZWJ that follow an eligible base are ignored in
	// later rules and do not introduce a break.
	if (cur.cm || cur.zwj) && canAttachLineMark(prev.lb) {
		return LineDontBreak
	}

	// LB10. Remaining CM/ZWJ behave as AL for this boundary and subsequent
	// boundaries. The consume step applies the same transformation.
	if cur.cm || cur.zwj {
		cur = cur.asAL()
	}

	// LB11.
	if cur.lb == LineBreakWJ || prev.lb == LineBreakWJ {
		return LineDontBreak
	}

	// LB12.
	if prev.lb == LineBreakGL {
		return LineDontBreak
	}

	// LB12a.
	if cur.lb == LineBreakGL && prev.lb != LineBreakSP && prev.lb != LineBreakBA && prev.lb != LineBreakHY && prev.lb != LineBreakHH {
		return LineDontBreak
	}

	// LB13.
	if cur.lb == LineBreakCL || cur.lb == LineBreakCP || cur.lb == LineBreakEX || cur.lb == LineBreakSY {
		return LineDontBreak
	}

	// LB14.
	if prev.lb == LineBreakOP || (prev.lb == LineBreakSP && st.haveBeforeSpaces && st.beforeSpaces.lb == LineBreakOP) {
		return LineDontBreak
	}

	// LB15a.
	if st.prevInitialQuote || (prev.lb == LineBreakSP && st.beforeSpacesInitialQuote) {
		return LineDontBreak
	}

	// LB15b.
	if isLineFinalQuote(cur) {
		next, ok := firstLineItem(rest)
		if !ok || isLineAfterFinalQuote(next.lb) {
			return LineDontBreak
		}
	}

	// LB15c.
	if prev.lb == LineBreakSP && cur.lb == LineBreakIS {
		next, ok := firstLineItem(rest)
		if ok && next.lb == LineBreakNU {
			return LineCanBreak
		}
	}

	// LB15d.
	if cur.lb == LineBreakIS {
		return LineDontBreak
	}

	// LB16.
	if cur.lb == LineBreakNS && (prev.lb == LineBreakCL || prev.lb == LineBreakCP || (prev.lb == LineBreakSP && st.haveBeforeSpaces && (st.beforeSpaces.lb == LineBreakCL || st.beforeSpaces.lb == LineBreakCP))) {
		return LineDontBreak
	}

	// LB17.
	if cur.lb == LineBreakB2 && (prev.lb == LineBreakB2 || (prev.lb == LineBreakSP && st.haveBeforeSpaces && st.beforeSpaces.lb == LineBreakB2)) {
		return LineDontBreak
	}

	// LB18.
	if prev.lb == LineBreakSP {
		return LineCanBreak
	}

	// LB19.
	if cur.lb == LineBreakQU && !isLineInitialQuote(cur) {
		return LineDontBreak
	}
	if prev.lb == LineBreakQU && !isLineFinalQuote(prev) {
		return LineDontBreak
	}

	// LB19a.
	if cur.lb == LineBreakQU && !isLineEastAsian(prev) {
		return LineDontBreak
	}
	if cur.lb == LineBreakQU {
		next, ok := firstLineItem(rest)
		if !ok || !isLineEastAsian(next) {
			return LineDontBreak
		}
	}
	if prev.lb == LineBreakQU && !isLineEastAsian(cur) {
		return LineDontBreak
	}
	if prev.lb == LineBreakQU && (!st.havePrev2 || !isLineEastAsian(st.prev2)) {
		return LineDontBreak
	}

	// LB20.
	if cur.lb == LineBreakCB || prev.lb == LineBreakCB {
		return LineCanBreak
	}

	// LB20a.
	if (cur.lb == LineBreakAL || cur.lb == LineBreakHL) && (prev.lb == LineBreakHY || prev.lb == LineBreakHH) && (!st.havePrev2 || isLineWordInitialHyphenBefore(st.prev2.lb)) {
		return LineDontBreak
	}

	// LB21.
	if cur.lb == LineBreakBA || cur.lb == LineBreakHH || cur.lb == LineBreakHY || cur.lb == LineBreakNS || prev.lb == LineBreakBB {
		return LineDontBreak
	}

	// LB21a.
	if cur.lb != LineBreakHL && (prev.lb == LineBreakHY || prev.lb == LineBreakHH) && st.havePrev2 && st.prev2.lb == LineBreakHL {
		return LineDontBreak
	}

	// LB21b.
	if prev.lb == LineBreakSY && cur.lb == LineBreakHL {
		return LineDontBreak
	}

	// LB22.
	if cur.lb == LineBreakIN {
		return LineDontBreak
	}

	// LB23.
	if (prev.lb == LineBreakAL || prev.lb == LineBreakHL) && cur.lb == LineBreakNU {
		return LineDontBreak
	}
	if prev.lb == LineBreakNU && (cur.lb == LineBreakAL || cur.lb == LineBreakHL) {
		return LineDontBreak
	}

	// LB23a.
	if prev.lb == LineBreakPR && (cur.lb == LineBreakID || cur.lb == LineBreakEB || cur.lb == LineBreakEM) {
		return LineDontBreak
	}
	if (prev.lb == LineBreakID || prev.lb == LineBreakEB || prev.lb == LineBreakEM) && cur.lb == LineBreakPO {
		return LineDontBreak
	}

	// LB24.
	if (prev.lb == LineBreakPR || prev.lb == LineBreakPO) && (cur.lb == LineBreakAL || cur.lb == LineBreakHL) {
		return LineDontBreak
	}
	if (prev.lb == LineBreakAL || prev.lb == LineBreakHL) && (cur.lb == LineBreakPR || cur.lb == LineBreakPO) {
		return LineDontBreak
	}

	// LB25.
	if st.matchesLineNumericNoBreakBefore(cur, rest) {
		return LineDontBreak
	}

	// LB26.
	if prev.lb == LineBreakJL && (cur.lb == LineBreakJL || cur.lb == LineBreakJV || cur.lb == LineBreakH2 || cur.lb == LineBreakH3) {
		return LineDontBreak
	}
	if (prev.lb == LineBreakJV || prev.lb == LineBreakH2) && (cur.lb == LineBreakJV || cur.lb == LineBreakJT) {
		return LineDontBreak
	}
	if (prev.lb == LineBreakJT || prev.lb == LineBreakH3) && cur.lb == LineBreakJT {
		return LineDontBreak
	}

	// LB27.
	if isLineHangulSyllable(prev.lb) && cur.lb == LineBreakPO {
		return LineDontBreak
	}
	if prev.lb == LineBreakPR && isLineHangulSyllable(cur.lb) {
		return LineDontBreak
	}

	// LB28.
	if (prev.lb == LineBreakAL || prev.lb == LineBreakHL) && (cur.lb == LineBreakAL || cur.lb == LineBreakHL) {
		return LineDontBreak
	}

	// LB28a.
	if prev.lb == LineBreakAP && isLineAksaraStart(cur) {
		return LineDontBreak
	}
	if isLineAksaraStart(prev) && (cur.lb == LineBreakVF || cur.lb == LineBreakVI) {
		return LineDontBreak
	}
	if prev.lb == LineBreakVI && st.havePrev2 && isLineAksaraStart(st.prev2) && isLineAksaraCore(cur) {
		return LineDontBreak
	}
	if isLineAksaraStart(prev) && isLineAksaraStart(cur) {
		next, ok := firstLineItem(rest)
		if ok && next.lb == LineBreakVF {
			return LineDontBreak
		}
	}

	// LB29.
	if prev.lb == LineBreakIS && (cur.lb == LineBreakAL || cur.lb == LineBreakHL) {
		return LineDontBreak
	}

	// LB30.
	if (prev.lb == LineBreakAL || prev.lb == LineBreakHL || prev.lb == LineBreakNU) && cur.lb == LineBreakOP && !isLineEastAsian(cur) {
		return LineDontBreak
	}
	if prev.lb == LineBreakCP && !isLineEastAsian(prev) && (cur.lb == LineBreakAL || cur.lb == LineBreakHL || cur.lb == LineBreakNU) {
		return LineDontBreak
	}

	// LB30a.
	if prev.lb == LineBreakRI && cur.lb == LineBreakRI && st.riRun%2 == 1 {
		return LineDontBreak
	}

	// LB30b.
	if (prev.lb == LineBreakEB || prev.extPicCn) && cur.lb == LineBreakEM {
		return LineDontBreak
	}

	// LB31.
	return LineCanBreak
}

func (st *lineState) consume(item lineItem, noBreak bool) {
	ignore := false
	if st.havePrev && (item.cm || item.zwj) && canAttachLineMark(st.prev.lb) {
		ignore = noBreak
	}
	if item.cm || item.zwj {
		if !ignore {
			item = item.asAL()
		}
	}

	st.afterZWJ = item.zwj
	if ignore {
		return
	}

	initialQuote := isLineInitialQuote(item) && st.initialQuoteContext()

	if st.havePrev {
		st.prev2 = st.prev
		st.havePrev2 = true
	}
	st.prev = item
	st.havePrev = true

	if item.lb == LineBreakSP {
		if st.prev2.lb != LineBreakSP {
			st.beforeSpaces = st.prev2
			st.haveBeforeSpaces = st.havePrev2
			st.beforeSpacesInitialQuote = st.prevInitialQuote
		}
	} else {
		st.haveBeforeSpaces = false
		st.beforeSpacesInitialQuote = false
	}

	st.prevInitialQuote = initialQuote
	st.updateLineNumericState(item.lb)
	if item.lb == LineBreakRI {
		st.riRun++
	} else {
		st.riRun = 0
	}
}

func (item lineItem) asAL() lineItem {
	item.lb = LineBreakAL
	item.gc = GeneralCategoryLu
	item.eaw = EastAsianWidthNa
	item.cm = false
	// Preserve item.zwj so LB8a still applies after a leading or otherwise
	// unattached ZWJ that LB10 treats as AL.
	item.extPicCn = false
	return item
}

func (st *lineState) initialQuoteContext() bool {
	if !st.havePrev {
		return true
	}
	switch st.prev.lb {
	case LineBreakBK, LineBreakCR, LineBreakLF, LineBreakNL, LineBreakOP, LineBreakQU, LineBreakGL, LineBreakSP, LineBreakZW:
		return true
	default:
		return false
	}
}

func (st *lineState) updateLineNumericState(lb LineBreakClass) {
	switch lb {
	case LineBreakNU:
		st.numState = lineNumSeq
	case LineBreakSY, LineBreakIS:
		if st.numState != lineNumSeq {
			st.numState = lineNumNone
		}
	case LineBreakCL, LineBreakCP:
		if st.numState == lineNumSeq {
			st.numState = lineNumClose
		} else {
			st.numState = lineNumNone
		}
	default:
		st.numState = lineNumNone
	}
}

func (st *lineState) matchesLineNumericNoBreakBefore(cur lineItem, rest string) bool {
	prev := st.prev
	if (cur.lb == LineBreakPO || cur.lb == LineBreakPR) && (st.numState == lineNumSeq || st.numState == lineNumClose) {
		return true
	}
	if (prev.lb == LineBreakPO || prev.lb == LineBreakPR) && cur.lb == LineBreakOP {
		return restStartsLineNumber(rest)
	}
	if (prev.lb == LineBreakPO || prev.lb == LineBreakPR) && cur.lb == LineBreakNU {
		return true
	}
	if prev.lb == LineBreakHY && cur.lb == LineBreakNU {
		return true
	}
	if prev.lb == LineBreakIS && cur.lb == LineBreakNU {
		return true
	}
	if st.numState == lineNumSeq && cur.lb == LineBreakNU {
		return true
	}
	return false
}

func restStartsLineNumber(rest string) bool {
	item, ok := firstLineItem(rest)
	if !ok {
		return false
	}
	if item.lb == LineBreakNU {
		return true
	}
	if item.lb == LineBreakIS {
		next, ok := firstLineItem(rest[item.end:])
		return ok && next.lb == LineBreakNU
	}
	return false
}

func firstLineItem(s string) (lineItem, bool) {
	if len(s) == 0 {
		return lineItem{}, false
	}
	return nextLineItem(s, 0)
}

func isLineHardBreak(lb LineBreakClass) bool {
	return lb == LineBreakBK || lb == LineBreakCR || lb == LineBreakLF || lb == LineBreakNL
}

func canAttachLineMark(lb LineBreakClass) bool {
	return lb != LineBreakBK && lb != LineBreakCR && lb != LineBreakLF && lb != LineBreakNL && lb != LineBreakSP && lb != LineBreakZW
}

func isLineInitialQuote(item lineItem) bool {
	return item.lb == LineBreakQU && item.gc == GeneralCategoryPi
}

func isLineFinalQuote(item lineItem) bool {
	return item.lb == LineBreakQU && item.gc == GeneralCategoryPf
}

func isLineAfterFinalQuote(lb LineBreakClass) bool {
	switch lb {
	case LineBreakSP, LineBreakGL, LineBreakWJ, LineBreakCL, LineBreakQU, LineBreakCP, LineBreakEX, LineBreakIS, LineBreakSY, LineBreakBK, LineBreakCR, LineBreakLF, LineBreakNL, LineBreakZW:
		return true
	default:
		return false
	}
}

func isLineEastAsian(item lineItem) bool {
	return item.eaw == EastAsianWidthF || item.eaw == EastAsianWidthW || item.eaw == EastAsianWidthH
}

func isLineWordInitialHyphenBefore(lb LineBreakClass) bool {
	switch lb {
	case LineBreakBK, LineBreakCR, LineBreakLF, LineBreakNL, LineBreakSP, LineBreakZW, LineBreakCB, LineBreakGL:
		return true
	default:
		return false
	}
}

func isLineHangulSyllable(lb LineBreakClass) bool {
	return lb == LineBreakJL || lb == LineBreakJV || lb == LineBreakJT || lb == LineBreakH2 || lb == LineBreakH3
}

func isLineAksaraStart(item lineItem) bool {
	return item.lb == LineBreakAK || item.lb == LineBreakAS || item.dotted
}

func isLineAksaraCore(item lineItem) bool {
	return item.lb == LineBreakAK || item.dotted
}
