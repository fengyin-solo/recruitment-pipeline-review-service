# Bug Reproduction

## Bug
Candidate keyword matching no longer checks the email field.

## Trigger
go test ./internal/service -run '^TestCandidateKeywordMatchesEmailAddress$' -count=1

## Error
email keyword should match one candidate
