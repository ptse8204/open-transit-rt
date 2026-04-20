# Conversation Summary

## Overall conclusion

The best opening product is **not** a full Swiftly/Passio/Connexionz replacement. The best wedge is:
- low-cost or BYOD GPS devices
- static GTFS import or lightweight GTFS authoring
- very good GTFS-RT Vehicle Positions
- simple hosting and onboarding
- Trip Updates treated as a later, pluggable module

## Major takeaways from the discussion

### 1. Free turnkey options
The old GTFS apps list did not clearly show a free turnkey stack that combines a free tracking app, GTFS-RT generation, and free hosted publishing.

### 2. How agencies really do GTFS-RT
Agencies usually do **not** buy GTFS-RT hosting as a standalone product. GTFS-RT is usually one output of a broader CAD/AVL, prediction, or passenger information platform.

### 3. Pricing
The broad pricing picture we discussed was:
- small/simple deployments: tens of thousands of dollars per year
- richer mid-size systems: low six figures per year
- full CAD/AVL replacements: millions to tens of millions

### 4. Market gap
Many agencies still do not publish active GTFS-RT Vehicle Positions feeds, which means there is still a real market gap for “get live vehicle positions online at all.”

### 5. Vendor comparison
- Passio looked strongest as a lower-cost turnkey bundle
- Swiftly looked strongest as a software overlay on existing hardware
- Connexionz looked most like a broader full-stack ITS deployment
- self-hosted open source looked cheapest on software but highest-risk operationally

### 6. California / Cal-ITP requirements
We separated:
- technical/data compliance requirements
- procurement/marketplace vendor requirements

Open-source can meet the first set. It does not automatically satisfy the second.

### 7. Open source usefulness
Yes, an open-source product like this would be useful, but the first useful version is narrower than a full enterprise replacement.

### 8. Why Trip Updates are hard
Trip Updates are not just trip completion detection. They require:
- correct trip instance matching
- stop-level ETAs
- block/interlining handling
- handling cancellations/additions/reroutes
- freshness and consistency across the system

### 9. AirTag idea
AirTags are not a good fit for GTFS-RT because they are consumer item-finders built around the Find My network, not deterministic backend fleet telemetry.

### 10. What is actually hard
The hard part is not collecting coordinates. The hard part is **transit inference**:
- identity and assignment
- matching GPS points to the right trip
- block handling
- ETA generation
- feed validation and reliability

### 11. Merge existing tools?
Yes, but as a platform with modules, not one giant monolith rewrite.

### 12. Would people use it before Trip Updates?
Yes. Vehicle Positions plus easier deployment is already useful. Trip Updates make it much more competitive.

### 13. Recommended module layout
- agency config
- GTFS Studio
- GTFS importer/publisher
- telemetry ingest
- state engine
- Vehicle Positions publisher
- prediction adapter
- Trip Updates publisher
- Alerts publisher
- validation/monitoring
- admin UI

### 14. GTFS Studio
GTFS Studio should support both:
- interactive authoring
- bring-your-own-GTFS import

### 15. Final direction
The implementation direction we settled on was:
- modular monorepo
- mostly Go
- Postgres/PostGIS
- Trip Updates as pluggable
- GTFS Studio in front of realtime
- first value centered on Vehicle Positions and usability for small agencies
