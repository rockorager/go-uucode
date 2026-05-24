module github.com/rockorager/go-uucode/bench/xunicode

go 1.25.5

require (
	charm.land/xunicode v0.0.0-20260310120343-e34c6fc66641
	github.com/rockorager/go-uucode v0.0.0
)

require golang.org/x/text v0.34.0 // indirect

replace github.com/rockorager/go-uucode => ../..
