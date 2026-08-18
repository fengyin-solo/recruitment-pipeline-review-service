# Bug Reproduction

## Bug
GetOffer returns the store's internal offer pointer.

## Trigger
go test ./internal/service -run '^TestOfferReturnedObjectCannotMutateStore$' -count=1

## Error
external offer pointer mutation changed store
