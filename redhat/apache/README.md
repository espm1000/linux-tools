# Apache2 Fun

To use a prebuilt Apache package with Ubuntu, use the `/local` folder
For manually building fun-ness, use the `/manual_build` folder

# Local

To build:
* cd to `scripts/`
* Run `./start_server.sh --build`

To run:  
* cd to `scripts/`
* Run `./start_server.sh --run`

## Local SSL Support

Use the `create_ca.sh` to create a Certificate Authority with OpenSSL
Use `create"${CERTS_FOLDER}_csr.sh --csr` to create a CSR against the CA
Use `create_csr.sh --sign` to sign the CSR against the CA
Note: Need to use the CA name of `local` for now...

# Manual Build

You're on your own.
