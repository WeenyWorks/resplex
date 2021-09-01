module github.com/WeenyWorks/resplex

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/WeenyWorks/resplex/lib/regheader v0.0.0-00010101000000-000000000000
	github.com/WeenyWorks/resplex/lib/visheader v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/xtaci/kcp-go/v5 v5.6.1
	github.com/xtaci/smux v1.5.15
	golang.org/x/mod v0.5.0 // indirect
	golang.org/x/sys v0.0.0-20210831042530-f4d43177bf5e // indirect
	golang.org/x/tools v0.1.6-0.20210802203754-9b21a8868e16 // indirect
	golang.org/x/tools/gopls v0.7.1 // indirect
	honnef.co/go/tools v0.2.1 // indirect
	mvdan.cc/xurls/v2 v2.3.0 // indirect
)

replace github.com/WeenyWorks/resplex/lib/regheader => ./lib/regheader

replace github.com/WeenyWorks/resplex/lib/visheader => ./lib/visheader
