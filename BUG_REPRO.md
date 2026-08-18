# Bug Reproduction

## Bug
ListInterviews returns internal interview pointers.

## Trigger
go test ./internal/service -run '^TestInterviewListObjectsCannotMutateStore$' -count=1

## Error
list result mutation changed stored interview
