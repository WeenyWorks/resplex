module github.com/WeenyWorks/lib/resplexc

replace github.com/WeenyWorks/resplex/lib/regheader => ../../lib/regheader

go 1.16

require (
	github.com/WeenyWorks/resplex/lib/regheader v0.0.0-00010101000000-000000000000
	github.com/xtaci/kcp-go/v5 v5.6.1
	github.com/xtaci/smux v1.5.15
)
