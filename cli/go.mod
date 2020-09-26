module alba

go 1.14

replace github.com/dadosjusbr/alba/git => /home/danielfireman/repos/dadosjusbr/alba/git

replace github.com/dadosjusbr/alba/storage => /home/danielfireman/repos/dadosjusbr/alba/storage

require (
	github.com/dadosjusbr/alba/git v0.0.0-20200925182652-7763ef421afb
	github.com/dadosjusbr/alba/storage v0.0.0-20200925182652-7763ef421afb
	github.com/matryer/is v1.4.0
	github.com/urfave/cli/v2 v2.2.0
)
