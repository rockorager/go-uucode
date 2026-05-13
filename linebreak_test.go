package uucode

import "testing"

func TestLineIterator(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []LineSegment
	}{
		{
			name: "empty",
			s:    "",
		},
		{
			name: "plain word",
			s:    "hello",
			want: []LineSegment{{Start: 0, End: 5, Break: LineMustBreak}},
		},
		{
			name: "space break",
			s:    "hello world",
			want: []LineSegment{{Start: 0, End: 6, Break: LineCanBreak}, {Start: 6, End: 11, Break: LineMustBreak}},
		},
		{
			name: "newline",
			s:    "a\nb",
			want: []LineSegment{{Start: 0, End: 2, Break: LineMustBreak}, {Start: 2, End: 3, Break: LineMustBreak}},
		},
		{
			name: "crlf",
			s:    "a\r\nb",
			want: []LineSegment{{Start: 0, End: 3, Break: LineMustBreak}, {Start: 3, End: 4, Break: LineMustBreak}},
		},
		{
			name: "numeric punctuation",
			s:    "A.B",
			want: []LineSegment{{Start: 0, End: 3, Break: LineMustBreak}},
		},
		{
			name: "emoji zwj sequence",
			s:    "👩‍🚀x",
			want: []LineSegment{{Start: 0, End: len("👩‍🚀"), Break: LineCanBreak}, {Start: len("👩‍🚀"), End: len("👩‍🚀x"), Break: LineMustBreak}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			it := NewLineIterator(tc.s)
			var got []LineSegment
			for {
				seg, ok := it.Next()
				if !ok {
					break
				}
				got = append(got, seg)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("segment count: got %#v, want %#v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("segment %d: got %#v, want %#v", i, got[i], tc.want[i])
				}
			}
		})
	}
}
