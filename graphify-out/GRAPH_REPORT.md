# Graph Report - .  (2026-04-13)

## Corpus Check
- Corpus is ~32,806 words - fits in a single context window. You may not need a graph.

## Summary
- 533 nodes · 763 edges · 68 communities detected
- Extraction: 99% EXTRACTED · 1% INFERRED · 0% AMBIGUOUS · INFERRED: 11 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Handler & Model Layer|Handler & Model Layer]]
- [[_COMMUNITY_Auth & Server Wiring|Auth & Server Wiring]]
- [[_COMMUNITY_Handler Tests|Handler Tests]]
- [[_COMMUNITY_SQLite Store|SQLite Store]]
- [[_COMMUNITY_Store Tests|Store Tests]]
- [[_COMMUNITY_API Client|API Client]]
- [[_COMMUNITY_Auth Tests|Auth Tests]]
- [[_COMMUNITY_Chronopost Tracker|Chronopost Tracker]]
- [[_COMMUNITY_Frontend App Graph|Frontend App Graph]]
- [[_COMMUNITY_Worker Tests|Worker Tests]]
- [[_COMMUNITY_Mondial Relay Tracker|Mondial Relay Tracker]]
- [[_COMMUNITY_Chronopost Tests|Chronopost Tests]]
- [[_COMMUNITY_Relais Colis Tracker|Relais Colis Tracker]]
- [[_COMMUNITY_GLS Tracker|GLS Tracker]]
- [[_COMMUNITY_Server Tests|Server Tests]]
- [[_COMMUNITY_VintedGo Tracker|VintedGo Tracker]]
- [[_COMMUNITY_La Poste Tracker|La Poste Tracker]]
- [[_COMMUNITY_La Poste Tests|La Poste Tests]]
- [[_COMMUNITY_VintedGo Tests|VintedGo Tests]]
- [[_COMMUNITY_Push Notifications|Push Notifications]]
- [[_COMMUNITY_Internationalization|Internationalization]]
- [[_COMMUNITY_DPD Tests|DPD Tests]]
- [[_COMMUNITY_GLS Tests|GLS Tests]]
- [[_COMMUNITY_Rate Limiter|Rate Limiter]]
- [[_COMMUNITY_Mondial Relay Tests|Mondial Relay Tests]]
- [[_COMMUNITY_Colis Prive Tests|Colis Prive Tests]]
- [[_COMMUNITY_Registry Tests|Registry Tests]]
- [[_COMMUNITY_DPD Tracker|DPD Tracker]]
- [[_COMMUNITY_Relais Colis Tests|Relais Colis Tests]]
- [[_COMMUNITY_Frontend Utilities|Frontend Utilities]]
- [[_COMMUNITY_Config Tests|Config Tests]]
- [[_COMMUNITY_Handler Helpers|Handler Helpers]]
- [[_COMMUNITY_Colis Prive Tracker|Colis Prive Tracker]]
- [[_COMMUNITY_Toast Notifications|Toast Notifications]]
- [[_COMMUNITY_Auth Handlers|Auth Handlers]]
- [[_COMMUNITY_App Entry Point|App Entry Point]]
- [[_COMMUNITY_Dashboard Page|Dashboard Page]]
- [[_COMMUNITY_Parcel Detail Page|Parcel Detail Page]]
- [[_COMMUNITY_Login Page|Login Page]]
- [[_COMMUNITY_Security Headers|Security Headers]]
- [[_COMMUNITY_SPA Handler|SPA Handler]]
- [[_COMMUNITY_Health Check|Health Check]]
- [[_COMMUNITY_Notification Handlers|Notification Handlers]]
- [[_COMMUNITY_Colissimo Tests|Colissimo Tests]]
- [[_COMMUNITY_Vite Config|Vite Config]]
- [[_COMMUNITY_Embed Frontend|Embed Frontend]]
- [[_COMMUNITY_Svelte Config|Svelte Config]]
- [[_COMMUNITY_Parcel Card|Parcel Card]]
- [[_COMMUNITY_Parcel Timeline|Parcel Timeline]]
- [[_COMMUNITY_Toast Container|Toast Container]]
- [[_COMMUNITY_Status Badge|Status Badge]]
- [[_COMMUNITY_Navbar|Navbar]]
- [[_COMMUNITY_Add Parcel Page|Add Parcel Page]]
- [[_COMMUNITY_Not Found Page|Not Found Page]]
- [[_COMMUNITY_Service Worker|Service Worker]]
- [[_COMMUNITY_Auth Middleware|Auth Middleware]]
- [[_COMMUNITY_Community 56|Community 56]]
- [[_COMMUNITY_Community 57|Community 57]]
- [[_COMMUNITY_Community 58|Community 58]]
- [[_COMMUNITY_Community 59|Community 59]]
- [[_COMMUNITY_Community 60|Community 60]]
- [[_COMMUNITY_Community 61|Community 61]]
- [[_COMMUNITY_Community 62|Community 62]]
- [[_COMMUNITY_Community 63|Community 63]]
- [[_COMMUNITY_Community 64|Community 64]]
- [[_COMMUNITY_Community 65|Community 65]]
- [[_COMMUNITY_Community 66|Community 66]]
- [[_COMMUNITY_Community 67|Community 67]]

## God Nodes (most connected - your core abstractions)
1. `newTestEnv()` - 34 edges
2. `newTestStore()` - 24 edges
3. `Handler` - 22 edges
4. `request()` - 18 edges
5. `withChiParam()` - 17 edges
6. `SQLiteStore` - 17 edges
7. `Tracker` - 14 edges
8. `NewRegistry()` - 14 edges
9. `ParcelStatus` - 13 edges
10. `CarrierCode` - 13 edges

## Surprising Connections (you probably didn't know these)
- `App Icon 192x192 - Package/Parcel Symbol` --conceptually_related_to--> `Parcel Tracker Application`  [INFERRED]
  web/public/icon-192.png → README.md
- `main()` --calls--> `Load()`  [EXTRACTED]
  cmd/parcel-tracker/main.go → internal/config/config.go
- `main()` --calls--> `NewSQLiteStore()`  [EXTRACTED]
  cmd/parcel-tracker/main.go → internal/store/sqlite.go
- `main()` --depends_on--> `Auth`  [EXTRACTED]
  cmd/parcel-tracker/main.go → internal/auth/middleware.go
- `main()` --depends_on--> `Handler`  [EXTRACTED]
  cmd/parcel-tracker/main.go → internal/handler/parcels.go

## Hyperedges (group relationships)
- **Supported Carrier Implementations** — readme_carrier_laposte, readme_carrier_colissimo, readme_carrier_chronopost, readme_carrier_mondial_relay, readme_carrier_gls, readme_carrier_dpd, readme_carrier_colis_prive, readme_carrier_relais_colis, readme_carrier_manual, claude_tracker_interface [EXTRACTED 1.00]
- **Layered Architecture Design Pattern** — claude_handler_layer, claude_store_layer, claude_tracker_layer, claude_auth_layer, claude_layered_architecture [EXTRACTED 1.00]
- **Bugs Discovered During Test Coverage Effort** — test_bug_unchecked_error_create_event, test_bug_race_condition_jwt_secret, test_bug_foreign_keys_not_enforced, test_bug_spa_cache_header [EXTRACTED 1.00]
- **Request Flow: Server -> Handler -> Store -> Model** — server_new, handler_handler, store_store, model_parcel, model_trackingevent [INFERRED 0.90]
- **Main Entrypoint Wires All Layers** — main_main, config_load, sqlite_newsqlitestore, auth_auth, handler_handler, tracker_newregistry, server_new, notifier_notifier [EXTRACTED 1.00]
- **All Layers Consume Model Types** — handler_handler, store_store, tracker_worker, notifier_notifier, model_parcel, model_trackingevent, model_carriercode, model_parcelstatus [EXTRACTED 1.00]

## Communities

### Community 0 - "Handler & Model Layer"
Cohesion: 0.04
Nodes (23): createEventRequest, createParcelRequest, updateParcelRequest, main(), CarrierCode, ParcelStatus, TrackingEvent, CarrierInfo (+15 more)

### Community 1 - "Auth & Server Wiring"
Cohesion: 0.13
Nodes (37): newTestEnv(), TestCheckAuthAuthenticated(), TestCheckAuthNotAuthenticated(), TestCreateEvent(), TestCreateEventDefaultStatus(), TestCreateEventMissingMessage(), TestCreateEventParcelNotFound(), TestCreateParcel() (+29 more)

### Community 2 - "Handler Tests"
Cohesion: 0.07
Nodes (15): Handler, HandlerNotifier, writeError(), writeJSON(), Parcel, Notifier, Payload, statusMessage() (+7 more)

### Community 3 - "SQLite Store"
Cohesion: 0.12
Nodes (21): ApiError, checkAuth(), createEvent(), createParcel(), deleteEvent(), deleteParcel(), getHealth(), getParcel() (+13 more)

### Community 4 - "Store Tests"
Cohesion: 0.13
Nodes (7): newID(), NewSQLiteStore(), scanParcel(), scanParcelFields(), scanParcelRow(), scanner, SQLiteStore

### Community 5 - "API Client"
Cohesion: 0.16
Nodes (24): newTestStore(), TestCreateAndGetParcel(), TestCreateAndListEvents(), TestCreateEventDefaultTimestamp(), TestCreateEventUpdatesParcelStatus(), TestCreateEventWithExplicitTimestamp(), TestDeleteEvent(), TestDeleteEventNotFound() (+16 more)

### Community 6 - "Auth Tests"
Cohesion: 0.17
Nodes (12): newTestAuth(), TestJWTSecretIsPersistent(), TestLoginNotConfigured(), TestLoginSuccess(), TestLoginWrongPassword(), TestMiddlewareAllowsValidToken(), TestMiddlewareRejectsInvalidToken(), TestMiddlewareRejectsNoToken() (+4 more)

### Community 7 - "Chronopost Tracker"
Cohesion: 0.17
Nodes (12): buildChronopostSOAPRequest(), buildLocation(), mapChronopostStatus(), parseChronopostDate(), parseChronopostResponse(), chronopostBody, chronopostEnvelope, chronopostEvent (+4 more)

### Community 8 - "Frontend App Graph"
Cohesion: 0.15
Nodes (15): Adding a New Carrier Guide, Handler Layer, No ORM - Raw SQL Rationale, Tracker Registry, SQLite WAL Mode, Store Interface, Store Layer, Tracker Interface (+7 more)

### Community 9 - "Worker Tests"
Cohesion: 0.2
Nodes (3): Auth, New(), New()

### Community 10 - "Mondial Relay Tracker"
Cohesion: 0.27
Nodes (8): newTestWorkerStore(), TestRefreshAllSkipsArchivedParcels(), TestRefreshAllSkipsManualCarrier(), TestRefreshAllSkipsTerminalStatuses(), TestRefreshParcelCreatesEvents(), TestRefreshParcelDeduplicatesEvents(), TestRefreshParcelUnknownCarrier(), TestWorkerRunStopsOnContextCancel()

### Community 11 - "Chronopost Tests"
Cohesion: 0.35
Nodes (8): collectTexts(), extractDateFromText(), hasClass(), mapMondialRelayStatus(), parseMondialRelayBlock(), parseMondialRelayDate(), parseMondialRelayHTML(), parseMondialRelayTimelineItem()

### Community 12 - "Relais Colis Tracker"
Cohesion: 0.2
Nodes (0): 

### Community 13 - "GLS Tracker"
Cohesion: 0.28
Nodes (8): mapRelaisColisStatus(), parseRelaisColisDate(), parseRelaisColisResponse(), relaisColisData, relaisColisEvent, relaisColisEventList, relaisColisResponse, relaisColisWrapper

### Community 14 - "Server Tests"
Cohesion: 0.31
Nodes (8): buildGLSLocation(), mapGLSStatus(), parseGLSDate(), parseGLSResponse(), glsAddress, glsEvent, glsResponse, glsShipment

### Community 15 - "VintedGo Tracker"
Cohesion: 0.39
Nodes (5): newTestServer(), TestHealthEndpoint(), TestPingEndpoint(), TestProtectedEndpointRequiresAuth(), TestSPAFallback()

### Community 16 - "La Poste Tracker"
Cohesion: 0.36
Nodes (7): vintedGoEvent, vintedGoLocation, vintedGoResponse, buildVintedGoLocation(), mapVintedGoStatus(), parseVintedGoDate(), parseVintedGoResponse()

### Community 17 - "La Poste Tests"
Cohesion: 0.32
Nodes (7): mapLaPosteStatus(), parseLaPosteDate(), parseLaPosteResponse(), laPosteContext, laPosteEvent, laPosteResponse, laPosteShipment

### Community 18 - "VintedGo Tests"
Cohesion: 0.25
Nodes (0): 

### Community 19 - "Push Notifications"
Cohesion: 0.25
Nodes (0): 

### Community 20 - "Internationalization"
Cohesion: 0.33
Nodes (2): getStatusLabel(), t()

### Community 21 - "DPD Tests"
Cohesion: 0.29
Nodes (0): 

### Community 22 - "GLS Tests"
Cohesion: 0.29
Nodes (0): 

### Community 23 - "Rate Limiter"
Cohesion: 0.29
Nodes (7): Embedded Frontend in Go Binary, Single Binary Deployment, Svelte 5 Runes - No Stores Rationale, Svelte SPA Frontend, Docker Deployment, Bug: SPA Cache Header Never Set, Gap: Frontend No Test Framework

### Community 24 - "Mondial Relay Tests"
Cohesion: 0.53
Nodes (5): Config, envDuration(), envInt(), envStr(), Load()

### Community 25 - "Colis Prive Tests"
Cohesion: 0.4
Nodes (2): bucket, RateLimiter

### Community 26 - "Registry Tests"
Cohesion: 0.33
Nodes (0): 

### Community 27 - "DPD Tracker"
Cohesion: 0.33
Nodes (0): 

### Community 28 - "Relais Colis Tests"
Cohesion: 0.33
Nodes (0): 

### Community 29 - "Frontend Utilities"
Cohesion: 0.6
Nodes (5): cleanDPDLabel(), mapDPDStatus(), parseDPDDate(), parseDPDHTML(), parseDPDRow()

### Community 30 - "Config Tests"
Cohesion: 0.33
Nodes (0): 

### Community 31 - "Handler Helpers"
Cohesion: 0.33
Nodes (6): Layered Architecture Design, App Icon 192x192 - Package/Parcel Symbol, App Icon 512x512 - Package/Parcel Symbol, Environment Variables Configuration, Parcel Tracker Application, Test Coverage Analysis

### Community 32 - "Colis Prive Tracker"
Cohesion: 0.4
Nodes (0): 

### Community 33 - "Toast Notifications"
Cohesion: 0.4
Nodes (0): 

### Community 34 - "Auth Handlers"
Cohesion: 0.7
Nodes (4): mapColisPriveStatus(), parseColisPriveDate(), parseColisPriveHTML(), parseColisPriveRow()

### Community 35 - "App Entry Point"
Cohesion: 0.5
Nodes (0): 

### Community 36 - "Dashboard Page"
Cohesion: 0.67
Nodes (2): loginRequest, setupRequest

### Community 37 - "Parcel Detail Page"
Cohesion: 1.0
Nodes (3): Auth Layer, JWT Authentication, Bug: Race Condition in jwtSecret

### Community 38 - "Login Page"
Cohesion: 1.0
Nodes (0): 

### Community 39 - "Security Headers"
Cohesion: 1.0
Nodes (0): 

### Community 40 - "SPA Handler"
Cohesion: 1.0
Nodes (0): 

### Community 41 - "Health Check"
Cohesion: 1.0
Nodes (0): 

### Community 42 - "Notification Handlers"
Cohesion: 1.0
Nodes (0): 

### Community 43 - "Colissimo Tests"
Cohesion: 1.0
Nodes (0): 

### Community 44 - "Vite Config"
Cohesion: 1.0
Nodes (0): 

### Community 45 - "Embed Frontend"
Cohesion: 1.0
Nodes (1): subscribeRequest

### Community 46 - "Svelte Config"
Cohesion: 1.0
Nodes (0): 

### Community 47 - "Parcel Card"
Cohesion: 1.0
Nodes (2): Colissimo Carrier, La Poste Carrier

### Community 48 - "Parcel Timeline"
Cohesion: 1.0
Nodes (0): 

### Community 49 - "Toast Container"
Cohesion: 1.0
Nodes (0): 

### Community 50 - "Status Badge"
Cohesion: 1.0
Nodes (0): 

### Community 51 - "Navbar"
Cohesion: 1.0
Nodes (0): 

### Community 52 - "Add Parcel Page"
Cohesion: 1.0
Nodes (0): 

### Community 53 - "Not Found Page"
Cohesion: 1.0
Nodes (0): 

### Community 54 - "Service Worker"
Cohesion: 1.0
Nodes (0): 

### Community 55 - "Auth Middleware"
Cohesion: 1.0
Nodes (0): 

### Community 56 - "Community 56"
Cohesion: 1.0
Nodes (0): 

### Community 57 - "Community 57"
Cohesion: 1.0
Nodes (0): 

### Community 58 - "Community 58"
Cohesion: 1.0
Nodes (0): 

### Community 59 - "Community 59"
Cohesion: 1.0
Nodes (0): 

### Community 60 - "Community 60"
Cohesion: 1.0
Nodes (1): Chronopost Carrier

### Community 61 - "Community 61"
Cohesion: 1.0
Nodes (1): Mondial Relay Carrier

### Community 62 - "Community 62"
Cohesion: 1.0
Nodes (1): GLS Carrier

### Community 63 - "Community 63"
Cohesion: 1.0
Nodes (1): DPD Carrier

### Community 64 - "Community 64"
Cohesion: 1.0
Nodes (1): Colis Prive Carrier

### Community 65 - "Community 65"
Cohesion: 1.0
Nodes (1): Relais Colis Carrier

### Community 66 - "Community 66"
Cohesion: 1.0
Nodes (1): Manual Carrier

### Community 67 - "Community 67"
Cohesion: 1.0
Nodes (1): CI/CD Pipeline

## Knowledge Gaps
- **58 isolated node(s):** `bucket`, `Config`, `updateParcelRequest`, `subscribeRequest`, `HandlerNotifier` (+53 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Login Page`** (2 nodes): `App.svelte`, `main.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Security Headers`** (2 nodes): `if()`, `Dashboard.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `SPA Handler`** (2 nodes): `if()`, `ParcelDetail.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Health Check`** (2 nodes): `if()`, `Login.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Notification Handlers`** (2 nodes): `security.go`, `SecurityHeaders()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Colissimo Tests`** (2 nodes): `spa.go`, `SPAHandler()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Vite Config`** (2 nodes): `HealthCheck()`, `health.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Embed Frontend`** (2 nodes): `subscribeRequest`, `notifications.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Svelte Config`** (2 nodes): `TestColissimoCodeAndName()`, `colissimo_test.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Parcel Card`** (2 nodes): `Colissimo Carrier`, `La Poste Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Parcel Timeline`** (1 nodes): `vite.config.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Toast Container`** (1 nodes): `embed.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Status Badge`** (1 nodes): `svelte.config.js`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Navbar`** (1 nodes): `ParcelCard.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Add Parcel Page`** (1 nodes): `ParcelTimeline.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Not Found Page`** (1 nodes): `ToastContainer.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Service Worker`** (1 nodes): `StatusBadge.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Auth Middleware`** (1 nodes): `Navbar.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 56`** (1 nodes): `AddParcel.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 57`** (1 nodes): `NotFound.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 58`** (1 nodes): `sw.js`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 59`** (1 nodes): `middleware.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 60`** (1 nodes): `Chronopost Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 61`** (1 nodes): `Mondial Relay Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 62`** (1 nodes): `GLS Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 63`** (1 nodes): `DPD Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 64`** (1 nodes): `Colis Prive Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 65`** (1 nodes): `Relais Colis Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 66`** (1 nodes): `Manual Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 67`** (1 nodes): `CI/CD Pipeline`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `main()` connect `Handler & Model Layer` to `Mondial Relay Tests`, `Worker Tests`, `Handler Tests`, `Store Tests`?**
  _High betweenness centrality (0.045) - this node is a cross-community bridge._
- **Why does `TrackingEvent` connect `Handler & Model Layer` to `Handler Tests`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `bucket`, `Config`, `updateParcelRequest` to the rest of the system?**
  _58 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Handler & Model Layer` be split into smaller, more focused modules?**
  _Cohesion score 0.04 - nodes in this community are weakly interconnected._
- **Should `Auth & Server Wiring` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._
- **Should `Handler Tests` be split into smaller, more focused modules?**
  _Cohesion score 0.07 - nodes in this community are weakly interconnected._
- **Should `SQLite Store` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._