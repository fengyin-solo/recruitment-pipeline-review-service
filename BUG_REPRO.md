# Bug Reproduction

## Bug
Rejected salary update mutates the stored job before validation returns.

## Trigger
go test ./internal/service -run '^TestFailedJobSalaryUpdateKeepsStoredJob$' -count=1

## Error
failed update changed stored salary range
