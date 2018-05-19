# Sherlock
**Transparent Proxy for Inspecting REST APIs**

## Architecture
### Transparent Proxy

*Purpose*
    
This entity is responsible for proxying every HTTP request and HTTP response pair
without the client or server knowing, and providing a subscription for entities
to listen to HTTP request-response pair payloads.

*Considerations*

* Initial implementation should utilize `github.com/elazarl/goproxy`, this allows 
this project to focus on other aspects of this projects while utilizing an existing reliable proxy solution.

### Accumulator

*Purpose*

This entity is responsible for taking request-response pairs from **Transparent Proxy**, 
filtering blacklisted or whitelisted domains, filtering JSON APIs, storing the 
raw HTTP payloads in `./data/` grouped by domain name, e.g. `./data/api.company.io`,
and making it easily consumable by **API Praser**.

*Considerations*

* 

### API Parser

*Purpose*

This entity is responsible for inspecting raw HTTP request response payloads
and generating a high level API spec.

*Considerations*

* 

