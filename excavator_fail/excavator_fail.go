package fail

fail

/*
This is a non-compiling file that has been added to explicitly ensure that CI fails.
It also contains the command that caused the failure and its output.
Remove this file if debugging locally.

go mod operation failed. This may mean that there are legitimate dependency issues with the "go.mod" definition in the repository and the updates performed by the gomod check. This branch can be cloned locally to debug the issue.

Command that caused error:
./godelw exec -- go get -d github.com/palantir/conjure-go-runtime/v2@upgrade github.com/palantir/godel/v2@upgrade github.com/palantir/pkg/bytesbuffers@upgrade github.com/palantir/pkg/cobracli@upgrade github.com/palantir/pkg/matcher@upgrade github.com/palantir/pkg/metrics@upgrade github.com/palantir/pkg/refreshable@upgrade github.com/palantir/pkg/retry@upgrade github.com/palantir/pkg/safejson@upgrade github.com/palantir/pkg/specdir@upgrade github.com/palantir/pkg/tlsconfig@upgrade github.com/palantir/pkg/uuid@upgrade github.com/palantir/witchcraft-go-error@upgrade github.com/palantir/witchcraft-go-logging@upgrade github.com/palantir/witchcraft-go-params@upgrade github.com/palantir/witchcraft-go-tracing@upgrade github.com/spf13/cobra@upgrade github.com/spf13/pflag@upgrade github.com/stretchr/testify@upgrade golang.org/x/net@upgrade golang.org/x/text@upgrade golang.org/x/tools@upgrade gopkg.in/yaml.v2@upgrade gopkg.in/yaml.v3@upgrade

Output:
go: -d flag is deprecated. -d=true is a no-op
go: github.com/palantir/conjure-go-runtime/v2@upgrade: module github.com/palantir/conjure-go-runtime/v2: Get "https://proxy.golang.org/github.com/palantir/conjure-go-runtime/v2/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/godel/v2@upgrade: module github.com/palantir/godel/v2: Get "https://proxy.golang.org/github.com/palantir/godel/v2/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/bytesbuffers@upgrade: module github.com/palantir/pkg/bytesbuffers: Get "https://proxy.golang.org/github.com/palantir/pkg/bytesbuffers/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/cobracli@upgrade: module github.com/palantir/pkg/cobracli: Get "https://proxy.golang.org/github.com/palantir/pkg/cobracli/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/matcher@upgrade: module github.com/palantir/pkg/matcher: Get "https://proxy.golang.org/github.com/palantir/pkg/matcher/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/metrics@upgrade: module github.com/palantir/pkg/metrics: Get "https://proxy.golang.org/github.com/palantir/pkg/metrics/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/refreshable@upgrade: module github.com/palantir/pkg/refreshable: Get "https://proxy.golang.org/github.com/palantir/pkg/refreshable/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/retry@upgrade: module github.com/palantir/pkg/retry: Get "https://proxy.golang.org/github.com/palantir/pkg/retry/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/safejson@upgrade: module github.com/palantir/pkg/safejson: Get "https://proxy.golang.org/github.com/palantir/pkg/safejson/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/specdir@upgrade: module github.com/palantir/pkg/specdir: Get "https://proxy.golang.org/github.com/palantir/pkg/specdir/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/tlsconfig@upgrade: module github.com/palantir/pkg/tlsconfig: Get "https://proxy.golang.org/github.com/palantir/pkg/tlsconfig/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/pkg/uuid@upgrade: module github.com/palantir/pkg/uuid: Get "https://proxy.golang.org/github.com/palantir/pkg/uuid/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/witchcraft-go-error@upgrade: module github.com/palantir/witchcraft-go-error: Get "https://proxy.golang.org/github.com/palantir/witchcraft-go-error/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/witchcraft-go-logging@upgrade: module github.com/palantir/witchcraft-go-logging: Get "https://proxy.golang.org/github.com/palantir/witchcraft-go-logging/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/witchcraft-go-params@upgrade: module github.com/palantir/witchcraft-go-params: Get "https://proxy.golang.org/github.com/palantir/witchcraft-go-params/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/palantir/witchcraft-go-tracing@upgrade: module github.com/palantir/witchcraft-go-tracing: Get "https://proxy.golang.org/github.com/palantir/witchcraft-go-tracing/@v/list": dial tcp: lookup proxy.golang.org on 10.2.218.250:53: server misbehaving
go: github.com/spf13/cobra@upgrade: module github.com/spf13/cobra: Get "https://proxy.golang.org/github.com/spf13/cobra/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
go: github.com/spf13/pflag@upgrade: module github.com/spf13/pflag: Get "https://proxy.golang.org/github.com/spf13/pflag/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
go: github.com/stretchr/testify@upgrade: module github.com/stretchr/testify: Get "https://proxy.golang.org/github.com/stretchr/testify/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
go: golang.org/x/net@upgrade: module golang.org/x/net: Get "https://proxy.golang.org/golang.org/x/net/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
go: golang.org/x/text@upgrade: module golang.org/x/text: Get "https://proxy.golang.org/golang.org/x/text/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
go: golang.org/x/tools@upgrade: module golang.org/x/tools: Get "https://proxy.golang.org/golang.org/x/tools/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
go: gopkg.in/yaml.v2@upgrade: module gopkg.in/yaml.v2: Get "https://proxy.golang.org/gopkg.in/yaml.v2/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
go: gopkg.in/yaml.v3@upgrade: module gopkg.in/yaml.v3: Get "https://proxy.golang.org/gopkg.in/yaml.v3/@v/list": dial tcp [2607:f8b0:4006:80c::2011]:443: connect: network is unreachable
Error: exit status 1

*/
