# creyaml
Generate example CR yaml from a CRD

## Install
 `make` creates a binary `bin/creyaml`

  Move the binary to your PATH
 
## Usage

### Using pipes

  `cat crd.yaml | creyaml`

### Using file

   `creyaml --file crd.yaml`
