# Bug Reproduction

## Bug
ListJobs computes a negative slice start when page is zero.

## Trigger
go test ./internal/service -run '^TestListJobsZeroPageDoesNotPanic$' -count=1

## Error
list jobs panicked for page zero
