# Various notes to self

## Design
* hashed passwords should be recorded in the server: see `LoadRegisteredUsers`
* FetchRemoteTemplates creates a global template list. When and if users are allowed to create templates, this will have side effects
* // TODO handle programme errors more systematically


## For production
* Proper security on 'Super Secret' key for the cookie store
* (in the API) generate secure apikeys

## Bugs
* Does not resynchronise with the api server if a user login fails, which it should do.  

## Display
* More detail for individual classes and industries, showing the demand they generate, etc

## Efficiency
* objects should be indirected much more instead of passing copies

## Custom jsonMarshal
* https://deepsource.com/blog/custom-marshalers-in-go

## Worth trying
* https://codyhouse.co/templates
* https://codepen.io/Elodieb/pen/OJyLxXm

