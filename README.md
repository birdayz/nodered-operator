# Node-RED Operator
This is currently work in progress, and yet in a brainstorming phase where i'm discovering how a cloud-native integration of Node-RED could look like.

## IDEAS
- CRD for Flows
- CRD for nodes, so nodes can be installed via CR! Makes automatic installation maybe easier? -> Could also be done via config, check if it's a good idea
- Auto restart node-red if necessary, i.e. when this nasty popup is shown
- high-level config for certain things (eg values),and snippets/mixins to add arbitrary config
- sidecars
- specific loggig config/support ?
- context store config
- credential store secret
