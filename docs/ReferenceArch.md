# Reference Astrolabe Server Architecture

# Overview
The reference server provides an implementation of the Astrolabe API specification.
It provides the REST API, S3 API and pluggable Protected Entity adapters.

# Reference Server
The reference server gets its configuration from a configuration directory specified on the
command line.  The configuration directory contains a number of JSON files named
*PE_type*.pe.json, for example, ivd.pe.json.

Each of the files contains configuration information for the given type.  Only PE types that
have config files will be exported by the server.  
# REST server implementation
The REST server HTTP interface is generated using the *go-swagger* tool from the *astrolabe_api.yaml* file.

