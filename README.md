# go-distributed-lock

A thread-safe distributed lock implementation with monotonic fencing tokens and lease expirations in Go.

## Features
- **Fencing Tokens**: Monotonically increasing 64-bit integer counters preventing split-brain writes during GC pauses.
- **Context-Aware Acquisition**: Bounded polling with `context.Context` timeout cancellation.
- **Lease Expiration**: Automatic lock expiration preventing deadlocks on node crash.
