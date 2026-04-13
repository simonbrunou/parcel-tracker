---
type: "query"
date: "2026-04-13T18:28:47.025152+00:00"
question: "Why does Registry bridge Handler & Model Layer, Auth & Server Wiring, and Chronopost Tracker?"
contributor: "graphify"
source_nodes: ["Registry", "Handler", "TrackingEvent", "CarrierCode", "Worker"]
---

# Q: Why does Registry bridge Handler & Model Layer, Auth & Server Wiring, and Chronopost Tracker?

## Answer

Registry (betweenness centrality 0.061) is the fan-out point where HTTP requests meet carrier implementations. Handler.RefreshParcel() calls Registry.Get(carrierCode) which routes to one of 11 tracker implementations. It bridges C0 (model types + carrier trackers), C1 (handler/auth/server wiring), and C7 (Chronopost, structurally unique due to SOAP/XML). Without Registry, Handler would need direct knowledge of every carrier, breaking the plugin architecture.

## Source Nodes

- Registry
- Handler
- TrackingEvent
- CarrierCode
- Worker