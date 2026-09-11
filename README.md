# common-sdk

Go networking utilities for multicast communication and supporting infrastructure experiments.

## Package map

| Package | Source-level scope |
| --- | --- |
| `multicast` | Fragmented messages and sender/receiver helpers |
| `compression` | gzip and LZ4 helpers with round-trip tests |
| `encryption` | Encryption helper implementation; not a security certification |
| `retry`, `pool` | Retry and connection-pool building blocks |
| `metrics`, `health`, `logging`, `errors` | Instrumentation and diagnostics |

```sh
go build -mod=readonly ./...
go test -mod=readonly ./...
```

Dependencies are committed in `go.mod` and `go.sum`. Run multicast programs only on a network you control.
Tests do not establish throughput, compression ratios, resource use or production reliability; no numerical performance claim is made here.

The automated test scope is currently compression. Other packages need targeted tests before production adoption. This is an experimental library, not a versioned stable SDK.

[Project map](https://github.com/swlee3306#selected-work)
