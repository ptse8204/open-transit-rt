# Malformed GTFS Fixture

This fixture intentionally omits `trips.txt` and references an unknown stop from `stop_times.txt`.

Later GTFS import phases should assert that this fixture fails validation and cannot activate a published feed version.
