# Bug Reproduction

## Bug
Candidate skills returned by CreateCandidate share the stored slice.

## Trigger
go test ./internal/service -run '^TestCandidateSkillsSnapshotDoesNotLeak$' -count=1

## Error
candidate skills leaked through returned object
