# Graph Report - .  (2026-04-13)

## Corpus Check
- 68 files · ~33,422 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 533 nodes · 679 edges · 79 communities detected
- Extraction: 98% EXTRACTED · 2% INFERRED · 0% AMBIGUOUS · INFERRED: 11 edges (avg confidence: 0.82)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]
- [[_COMMUNITY_Community 12|Community 12]]
- [[_COMMUNITY_Community 13|Community 13]]
- [[_COMMUNITY_Community 14|Community 14]]
- [[_COMMUNITY_Community 15|Community 15]]
- [[_COMMUNITY_Community 16|Community 16]]
- [[_COMMUNITY_Community 17|Community 17]]
- [[_COMMUNITY_Community 18|Community 18]]
- [[_COMMUNITY_Community 19|Community 19]]
- [[_COMMUNITY_Community 20|Community 20]]
- [[_COMMUNITY_Community 21|Community 21]]
- [[_COMMUNITY_Community 22|Community 22]]
- [[_COMMUNITY_Community 23|Community 23]]
- [[_COMMUNITY_Community 24|Community 24]]
- [[_COMMUNITY_Community 25|Community 25]]
- [[_COMMUNITY_Community 26|Community 26]]
- [[_COMMUNITY_Community 27|Community 27]]
- [[_COMMUNITY_Community 28|Community 28]]
- [[_COMMUNITY_Community 29|Community 29]]
- [[_COMMUNITY_Community 30|Community 30]]
- [[_COMMUNITY_Community 31|Community 31]]
- [[_COMMUNITY_Community 32|Community 32]]
- [[_COMMUNITY_Community 33|Community 33]]
- [[_COMMUNITY_Community 34|Community 34]]
- [[_COMMUNITY_Community 35|Community 35]]
- [[_COMMUNITY_Community 36|Community 36]]
- [[_COMMUNITY_Community 37|Community 37]]
- [[_COMMUNITY_Community 38|Community 38]]
- [[_COMMUNITY_Community 39|Community 39]]
- [[_COMMUNITY_Community 40|Community 40]]
- [[_COMMUNITY_Community 41|Community 41]]
- [[_COMMUNITY_Community 42|Community 42]]
- [[_COMMUNITY_Community 43|Community 43]]
- [[_COMMUNITY_Community 44|Community 44]]
- [[_COMMUNITY_Community 45|Community 45]]
- [[_COMMUNITY_Community 46|Community 46]]
- [[_COMMUNITY_Community 47|Community 47]]
- [[_COMMUNITY_Community 48|Community 48]]
- [[_COMMUNITY_Community 49|Community 49]]
- [[_COMMUNITY_Community 50|Community 50]]
- [[_COMMUNITY_Community 51|Community 51]]
- [[_COMMUNITY_Community 52|Community 52]]
- [[_COMMUNITY_Community 53|Community 53]]
- [[_COMMUNITY_Community 54|Community 54]]
- [[_COMMUNITY_Community 55|Community 55]]
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
- [[_COMMUNITY_Community 68|Community 68]]
- [[_COMMUNITY_Community 69|Community 69]]
- [[_COMMUNITY_Community 70|Community 70]]
- [[_COMMUNITY_Community 71|Community 71]]
- [[_COMMUNITY_Community 72|Community 72]]
- [[_COMMUNITY_Community 73|Community 73]]
- [[_COMMUNITY_Community 74|Community 74]]
- [[_COMMUNITY_Community 75|Community 75]]
- [[_COMMUNITY_Community 76|Community 76]]
- [[_COMMUNITY_Community 77|Community 77]]
- [[_COMMUNITY_Community 78|Community 78]]

## God Nodes (most connected - your core abstractions)
1. `newTestEnv()` - 34 edges
2. `newTestStore()` - 24 edges
3. `request()` - 18 edges
4. `Handler` - 17 edges
5. `withChiParam()` - 17 edges
6. `SQLiteStore` - 17 edges
7. `newTestAuth()` - 12 edges
8. `newTestWorkerStore()` - 8 edges
9. `Auth` - 7 edges
10. `parseMondialRelayHTML()` - 6 edges

## Surprising Connections (you probably didn't know these)
- `Parcel Tracker Application` --conceptually_related_to--> `App Icon 192x192 - Package/Parcel Symbol`  [INFERRED]
  README.md → web/public/icon-192.png
- `Parcel Tracker Application` --references--> `Test Coverage Analysis`  [INFERRED]
  README.md → TEST_COVERAGE_ANALYSIS.md
- `Docker Deployment` --conceptually_related_to--> `Single Binary Deployment`  [INFERRED]
  README.md → CLAUDE.md
- `REST API Routes` --references--> `Handler Layer`  [INFERRED]
  README.md → CLAUDE.md
- `Environment Variables Configuration` --conceptually_related_to--> `Layered Architecture Design`  [INFERRED]
  README.md → CLAUDE.md

## Hyperedges (group relationships)
- **Supported Carrier Implementations** — readme_carrier_laposte, readme_carrier_colissimo, readme_carrier_chronopost, readme_carrier_mondial_relay, readme_carrier_gls, readme_carrier_dpd, readme_carrier_colis_prive, readme_carrier_relais_colis, readme_carrier_manual, claude_tracker_interface [EXTRACTED 1.00]
- **Layered Architecture Design Pattern** — claude_handler_layer, claude_store_layer, claude_tracker_layer, claude_auth_layer, claude_layered_architecture [EXTRACTED 1.00]
- **Bugs Discovered During Test Coverage Effort** — test_bug_unchecked_error_create_event, test_bug_race_condition_jwt_secret, test_bug_foreign_keys_not_enforced, test_bug_spa_cache_header [EXTRACTED 1.00]
- **Request Flow: Server -> Handler -> Store -> Model** — server_new, handler_handler, store_store, model_parcel, model_trackingevent [INFERRED 0.90]
- **Main Entrypoint Wires All Layers** — main_main, config_load, sqlite_newsqlitestore, auth_auth, handler_handler, tracker_newregistry, server_new, notifier_notifier [EXTRACTED 1.00]
- **All Layers Consume Model Types** — handler_handler, store_store, tracker_worker, notifier_notifier, model_parcel, model_trackingevent, model_carriercode, model_parcelstatus [EXTRACTED 1.00]

## Communities

### Community 0 - "Community 0"
Cohesion: 0.13
Nodes (37): newTestEnv(), TestCheckAuthAuthenticated(), TestCheckAuthNotAuthenticated(), TestCreateEvent(), TestCreateEventDefaultStatus(), TestCreateEventMissingMessage(), TestCreateEventParcelNotFound(), TestCreateParcel() (+29 more)

### Community 1 - "Community 1"
Cohesion: 0.13
Nodes (7): newID(), NewSQLiteStore(), scanParcel(), scanParcelFields(), scanParcelRow(), scanner, SQLiteStore

### Community 2 - "Community 2"
Cohesion: 0.16
Nodes (24): newTestStore(), TestCreateAndGetParcel(), TestCreateAndListEvents(), TestCreateEventDefaultTimestamp(), TestCreateEventUpdatesParcelStatus(), TestCreateEventWithExplicitTimestamp(), TestDeleteEvent(), TestDeleteEventNotFound() (+16 more)

### Community 3 - "Community 3"
Cohesion: 0.12
Nodes (4): Handler, HandlerNotifier, writeError(), writeJSON()

### Community 4 - "Community 4"
Cohesion: 0.18
Nodes (19): ApiError, checkAuth(), createEvent(), createParcel(), deleteEvent(), deleteParcel(), getHealth(), getParcel() (+11 more)

### Community 5 - "Community 5"
Cohesion: 0.17
Nodes (12): newTestAuth(), TestJWTSecretIsPersistent(), TestLoginNotConfigured(), TestLoginSuccess(), TestLoginWrongPassword(), TestMiddlewareAllowsValidToken(), TestMiddlewareRejectsInvalidToken(), TestMiddlewareRejectsNoToken() (+4 more)

### Community 6 - "Community 6"
Cohesion: 0.15
Nodes (13): buildChronopostSOAPRequest(), buildLocation(), mapChronopostStatus(), parseChronopostDate(), parseChronopostResponse(), chronopostBody, chronopostEnvelope, chronopostEvent (+5 more)

### Community 7 - "Community 7"
Cohesion: 0.22
Nodes (9): collectTexts(), extractDateFromText(), hasClass(), mapMondialRelayStatus(), parseMondialRelayBlock(), parseMondialRelayDate(), parseMondialRelayHTML(), parseMondialRelayTimelineItem() (+1 more)

### Community 8 - "Community 8"
Cohesion: 0.15
Nodes (15): Adding a New Carrier Guide, Handler Layer, No ORM - Raw SQL Rationale, Tracker Registry, SQLite WAL Mode, Store Interface, Store Layer, Tracker Interface (+7 more)

### Community 9 - "Community 9"
Cohesion: 0.16
Nodes (7): CarrierCode, Parcel, ParcelStatus, TrackingEvent, Notifier, Payload, statusMessage()

### Community 10 - "Community 10"
Cohesion: 0.19
Nodes (9): mapRelaisColisStatus(), parseRelaisColisDate(), parseRelaisColisResponse(), relaisColisData, relaisColisEvent, relaisColisEventList, relaisColisResponse, RelaisColisTracker (+1 more)

### Community 11 - "Community 11"
Cohesion: 0.2
Nodes (9): buildGLSLocation(), mapGLSStatus(), parseGLSDate(), parseGLSResponse(), glsAddress, glsEvent, glsResponse, glsShipment (+1 more)

### Community 12 - "Community 12"
Cohesion: 0.22
Nodes (8): vintedGoEvent, vintedGoLocation, vintedGoResponse, VintedGoTracker, buildVintedGoLocation(), mapVintedGoStatus(), parseVintedGoDate(), parseVintedGoResponse()

### Community 13 - "Community 13"
Cohesion: 0.21
Nodes (8): mapLaPosteStatus(), parseLaPosteDate(), parseLaPosteResponse(), laPosteContext, laPosteEvent, laPosteResponse, laPosteShipment, LaPosteTracker

### Community 14 - "Community 14"
Cohesion: 0.24
Nodes (2): Auth, New()

### Community 15 - "Community 15"
Cohesion: 0.27
Nodes (8): newTestWorkerStore(), TestRefreshAllSkipsArchivedParcels(), TestRefreshAllSkipsManualCarrier(), TestRefreshAllSkipsTerminalStatuses(), TestRefreshParcelCreatesEvents(), TestRefreshParcelDeduplicatesEvents(), TestRefreshParcelUnknownCarrier(), TestWorkerRunStopsOnContextCancel()

### Community 16 - "Community 16"
Cohesion: 0.29
Nodes (6): cleanDPDLabel(), mapDPDStatus(), parseDPDDate(), parseDPDHTML(), parseDPDRow(), DPDTracker

### Community 17 - "Community 17"
Cohesion: 0.31
Nodes (5): mapColisPriveStatus(), parseColisPriveDate(), parseColisPriveHTML(), parseColisPriveRow(), ColisPriveTracker

### Community 18 - "Community 18"
Cohesion: 0.2
Nodes (0): 

### Community 19 - "Community 19"
Cohesion: 0.25
Nodes (5): CarrierInfo, NewRegistry(), Registry, Tracker, TrackResult

### Community 20 - "Community 20"
Cohesion: 0.39
Nodes (5): newTestServer(), TestHealthEndpoint(), TestPingEndpoint(), TestProtectedEndpointRequiresAuth(), TestSPAFallback()

### Community 21 - "Community 21"
Cohesion: 0.25
Nodes (0): 

### Community 22 - "Community 22"
Cohesion: 0.25
Nodes (0): 

### Community 23 - "Community 23"
Cohesion: 0.33
Nodes (2): togglePush(), urlBase64ToUint8Array()

### Community 24 - "Community 24"
Cohesion: 0.33
Nodes (2): getStatusLabel(), t()

### Community 25 - "Community 25"
Cohesion: 0.29
Nodes (0): 

### Community 26 - "Community 26"
Cohesion: 0.43
Nodes (3): Worker, WorkerNotifier, EventKey()

### Community 27 - "Community 27"
Cohesion: 0.29
Nodes (0): 

### Community 28 - "Community 28"
Cohesion: 0.29
Nodes (7): Embedded Frontend in Go Binary, Single Binary Deployment, Svelte 5 Runes - No Stores Rationale, Svelte SPA Frontend, Docker Deployment, Bug: SPA Cache Header Never Set, Gap: Frontend No Test Framework

### Community 29 - "Community 29"
Cohesion: 0.4
Nodes (2): bucket, RateLimiter

### Community 30 - "Community 30"
Cohesion: 0.53
Nodes (5): Config, envDuration(), envInt(), envStr(), Load()

### Community 31 - "Community 31"
Cohesion: 0.33
Nodes (0): 

### Community 32 - "Community 32"
Cohesion: 0.33
Nodes (0): 

### Community 33 - "Community 33"
Cohesion: 0.33
Nodes (0): 

### Community 34 - "Community 34"
Cohesion: 0.33
Nodes (0): 

### Community 35 - "Community 35"
Cohesion: 0.33
Nodes (6): Layered Architecture Design, App Icon 192x192 - Package/Parcel Symbol, App Icon 512x512 - Package/Parcel Symbol, Environment Variables Configuration, Parcel Tracker Application, Test Coverage Analysis

### Community 36 - "Community 36"
Cohesion: 0.4
Nodes (0): 

### Community 37 - "Community 37"
Cohesion: 0.4
Nodes (0): 

### Community 38 - "Community 38"
Cohesion: 0.4
Nodes (1): MockTracker

### Community 39 - "Community 39"
Cohesion: 0.4
Nodes (1): ManualTracker

### Community 40 - "Community 40"
Cohesion: 0.4
Nodes (1): ColissimoTracker

### Community 41 - "Community 41"
Cohesion: 0.4
Nodes (4): PaginatedParcels, ParcelFilter, PushSubscription, Store

### Community 42 - "Community 42"
Cohesion: 0.5
Nodes (0): 

### Community 43 - "Community 43"
Cohesion: 0.67
Nodes (2): createParcelRequest, updateParcelRequest

### Community 44 - "Community 44"
Cohesion: 0.67
Nodes (2): loginRequest, setupRequest

### Community 45 - "Community 45"
Cohesion: 1.0
Nodes (3): Auth Layer, JWT Authentication, Bug: Race Condition in jwtSecret

### Community 46 - "Community 46"
Cohesion: 1.0
Nodes (0): 

### Community 47 - "Community 47"
Cohesion: 1.0
Nodes (0): 

### Community 48 - "Community 48"
Cohesion: 1.0
Nodes (0): 

### Community 49 - "Community 49"
Cohesion: 1.0
Nodes (0): 

### Community 50 - "Community 50"
Cohesion: 1.0
Nodes (0): 

### Community 51 - "Community 51"
Cohesion: 1.0
Nodes (0): 

### Community 52 - "Community 52"
Cohesion: 1.0
Nodes (0): 

### Community 53 - "Community 53"
Cohesion: 1.0
Nodes (0): 

### Community 54 - "Community 54"
Cohesion: 1.0
Nodes (0): 

### Community 55 - "Community 55"
Cohesion: 1.0
Nodes (1): subscribeRequest

### Community 56 - "Community 56"
Cohesion: 1.0
Nodes (1): createEventRequest

### Community 57 - "Community 57"
Cohesion: 1.0
Nodes (0): 

### Community 58 - "Community 58"
Cohesion: 1.0
Nodes (2): Colissimo Carrier, La Poste Carrier

### Community 59 - "Community 59"
Cohesion: 1.0
Nodes (0): 

### Community 60 - "Community 60"
Cohesion: 1.0
Nodes (0): 

### Community 61 - "Community 61"
Cohesion: 1.0
Nodes (0): 

### Community 62 - "Community 62"
Cohesion: 1.0
Nodes (0): 

### Community 63 - "Community 63"
Cohesion: 1.0
Nodes (0): 

### Community 64 - "Community 64"
Cohesion: 1.0
Nodes (0): 

### Community 65 - "Community 65"
Cohesion: 1.0
Nodes (0): 

### Community 66 - "Community 66"
Cohesion: 1.0
Nodes (0): 

### Community 67 - "Community 67"
Cohesion: 1.0
Nodes (0): 

### Community 68 - "Community 68"
Cohesion: 1.0
Nodes (0): 

### Community 69 - "Community 69"
Cohesion: 1.0
Nodes (0): 

### Community 70 - "Community 70"
Cohesion: 1.0
Nodes (0): 

### Community 71 - "Community 71"
Cohesion: 1.0
Nodes (1): Chronopost Carrier

### Community 72 - "Community 72"
Cohesion: 1.0
Nodes (1): Mondial Relay Carrier

### Community 73 - "Community 73"
Cohesion: 1.0
Nodes (1): GLS Carrier

### Community 74 - "Community 74"
Cohesion: 1.0
Nodes (1): DPD Carrier

### Community 75 - "Community 75"
Cohesion: 1.0
Nodes (1): Colis Prive Carrier

### Community 76 - "Community 76"
Cohesion: 1.0
Nodes (1): Relais Colis Carrier

### Community 77 - "Community 77"
Cohesion: 1.0
Nodes (1): Manual Carrier

### Community 78 - "Community 78"
Cohesion: 1.0
Nodes (1): CI/CD Pipeline

## Knowledge Gaps
- **67 isolated node(s):** `TrackingEvent`, `bucket`, `Config`, `createParcelRequest`, `updateParcelRequest` (+62 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Community 46`** (2 nodes): `App.svelte`, `main.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 47`** (2 nodes): `if()`, `Dashboard.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 48`** (2 nodes): `if()`, `ParcelDetail.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 49`** (2 nodes): `if()`, `Login.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 50`** (2 nodes): `main.go`, `main()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 51`** (2 nodes): `security.go`, `SecurityHeaders()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 52`** (2 nodes): `spa.go`, `SPAHandler()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 53`** (2 nodes): `server.go`, `New()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 54`** (2 nodes): `HealthCheck()`, `health.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 55`** (2 nodes): `subscribeRequest`, `notifications.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 56`** (2 nodes): `createEventRequest`, `events.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 57`** (2 nodes): `TestColissimoCodeAndName()`, `colissimo_test.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 58`** (2 nodes): `Colissimo Carrier`, `La Poste Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 59`** (1 nodes): `vite.config.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 60`** (1 nodes): `embed.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 61`** (1 nodes): `svelte.config.js`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 62`** (1 nodes): `ParcelCard.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 63`** (1 nodes): `ParcelTimeline.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 64`** (1 nodes): `ToastContainer.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 65`** (1 nodes): `StatusBadge.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 66`** (1 nodes): `Navbar.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 67`** (1 nodes): `AddParcel.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 68`** (1 nodes): `NotFound.svelte`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 69`** (1 nodes): `sw.js`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 70`** (1 nodes): `middleware.go`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 71`** (1 nodes): `Chronopost Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 72`** (1 nodes): `Mondial Relay Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 73`** (1 nodes): `GLS Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 74`** (1 nodes): `DPD Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 75`** (1 nodes): `Colis Prive Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 76`** (1 nodes): `Relais Colis Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 77`** (1 nodes): `Manual Carrier`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Community 78`** (1 nodes): `CI/CD Pipeline`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **What connects `TrackingEvent`, `bucket`, `Config` to the rest of the system?**
  _67 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._
- **Should `Community 3` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._