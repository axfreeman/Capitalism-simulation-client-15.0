# Frontend for the simulation of a capitalist economy  
* Interface with remote server via its API.  
* This version uses http://github.com/axfreeman/capitalist-simulation-api-15.0    
* There is more documentation there.
     
## Philosophy
This version focusses on the UI logic by reducing authentication to the bare minimum possible.  

This client uses api keys issued by the server to access user data.  

This version uses only basic go http authentication.  

It creates a local database of registered users at startup, and provides a facility to register a new user and synchronise that with the server.  

A users can log in with the credentials of any registered user. The client then retrieves the current state of the simulation from the api server.  
