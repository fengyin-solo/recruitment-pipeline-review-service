# Bug Reproduction

## Bug
Rejected interview score update mutates the stored interview before validation returns.

## Trigger
go test ./internal/service -run '^TestFailedInterviewScoreUpdateKeepsStoredInterview$' -count=1

## Error
failed update changed stored score
