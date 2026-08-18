# Bug Reproduction

## Bug
Offer acceptance fills a job only when accepted offers exceed headcount.

## Trigger
go test ./internal/service -run '^TestAcceptedOfferFillsSingleHeadcountJob$' -count=1

## Error
single-headcount job should be filled after accepted offer
