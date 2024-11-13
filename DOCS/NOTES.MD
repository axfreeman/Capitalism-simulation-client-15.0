# Various notes to self

## Design
* hashed passwords should be recorded in the server: see `LoadRegisteredUsers`
* FetchRemoteTemplates creates a global template list. When users are allowed to create templates, this will have side effects
* Handle programme errors more systematically
* Proper security on 'Super Secret' key for the cookie store
* (in the API) generate secure apikeys

## Custom jsonMarshal
* https://deepsource.com/blog/custom-marshalers-in-go

## Worth trying as regards accordions
* https://codyhouse.co/templates
* https://codepen.io/Elodieb/pen/OJyLxXm

## Display Improvements
* Individual tables (eg Industries, stocks ) should switch between size, value, and price
* More detail for individual classes and industries, showing the demand they generate, etc

## Simulation as such
* Invest: don't transfer money to the owners if it is needed to maintain current output level

## BUGS
* Does not resynchronise with the api server if a user login fails, which it should do.  
* check that objects are indirected everywhere instead of passing copies
* Class stocks: doesn't display
* Missing column in Class value display

## Cycle analysis
* https://github.com/KyleBanks/depth

## Serving static files
* https://www.codershood.info/2020/02/16/serving-static-files-in-golang-using-gorilla-mux/
* NOTE: cannot mix static file server managed by http.Handle and by Router (set to mux.NewRouter())
