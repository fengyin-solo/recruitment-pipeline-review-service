# Bug Reproduction

## Bug
GetResume returns the store's internal resume pointer.

## Trigger
go test ./internal/service -run '^TestResumeReturnedObjectCannotMutateStore$' -count=1

## Error
external resume pointer mutation changed store
