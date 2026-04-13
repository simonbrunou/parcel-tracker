---
type: "query"
date: "2026-04-13T18:28:51.113082+00:00"
question: "What connects Parcel TrackingEvent bucket to the rest of the system?"
contributor: "graphify"
source_nodes: ["Parcel", "TrackingEvent", "Store", "Handler", "Worker"]
---

# Q: What connects Parcel TrackingEvent bucket to the rest of the system?

## Answer

Parcel and TrackingEvent are central model types used by every layer: Handler creates/reads them, Store persists them, Worker refreshes them, Notifier alerts on status changes. AST missed these cross-package type dependencies because tree-sitter cannot resolve Go type references. Deep mode semantic extraction recovered 147 cross-package edges, reducing isolated nodes from 252 to 34. The bucket struct in ratelimit.go is correctly isolated — it is a private implementation detail.

## Source Nodes

- Parcel
- TrackingEvent
- Store
- Handler
- Worker