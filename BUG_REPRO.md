# Bug Reproduction

## Bug
GetJob returns the store's internal job pointer.

## Trigger
go test ./internal/service -run '^TestJobReturnedObjectCannotMutateStore$' -count=1

## Error
external job pointer mutation changed store
