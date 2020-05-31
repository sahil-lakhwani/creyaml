# creyaml
Generate example CR yaml from a CRD

## Install
 `make` creates a binary `bin/creyaml`

  Move the binary to your PATH
 
## Usage

`kubectl get crd example.domain.com -o yaml | creyaml`