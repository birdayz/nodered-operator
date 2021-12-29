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
- secrets, user ref/ secrets via...? direct vault, msft vault, .. ?
- module/package CRD

## TODO
- Technical admin user for Operator to talk to node-red api


## Development

Generate stuff:
```
make
make generate
make manifests
```

Run locally (some things may not work, as the operator tries to communicate with Node-RED instances via in-cluster networking):
```
make run
```

Build+Push img:
```
make docker-build docker-push IMG=ghcr.io/birdayz/nodered-operator:latest
```

Deploy:
```
kubectl apply -k config/default
```
