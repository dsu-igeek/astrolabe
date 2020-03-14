# S3 Repository

# Repository format

The S3 repository stores Protected Entities as objects with three prefixes.

## Protected Entity Info objects
<*prefix*>/\<*pe type*>/peinfo/\<*pe id*>

The PE Info object is formatted as a JSON object
## Metadata objects
<*prefix*>/\<*pe type*>/md/\<*pe id*>.md
## Data objects
\<*prefix*>/\<*pe type*>/md/\<*pe id*>.\<*starting offset*>.data
