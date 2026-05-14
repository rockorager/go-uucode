module github.com/rockorager/go-uucode/bench/runewidth

go 1.21

require (
	github.com/mattn/go-runewidth v0.0.23
	github.com/rockorager/go-uucode v0.0.0
)

require github.com/clipperhouse/uax29/v2 v2.2.0 // indirect

replace github.com/rockorager/go-uucode => ../..
