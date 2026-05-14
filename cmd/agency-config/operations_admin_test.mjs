import test from "node:test";
import assert from "node:assert/strict";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const operations = require("./operations_admin.js");

test("safeAdminPath allows only relative private Operations JSON", () => {
  assert.equal(operations.safeAdminPath("/admin/operations.json"), "/admin/operations.json");
  assert.equal(operations.safeAdminPath("/admin/operations/feed-health.json"), "/admin/operations/feed-health.json");
  assert.equal(operations.safeAdminPath("/admin/operations/validation-health/refresh.json", { method: "POST" }), "/admin/operations/validation-health/refresh.json");

  assert.equal(operations.safeAdminPath("https://example.test/admin/operations.json"), "");
  assert.equal(operations.safeAdminPath("//example.test/admin/operations.json"), "");
  assert.equal(operations.safeAdminPath("/public/gtfsrt/vehicle_positions.json"), "");
  assert.equal(operations.safeAdminPath("/v1/events"), "");
  assert.equal(operations.safeAdminPath("/admin/operations/feed-health.json?agency_id=other"), "");
  assert.equal(operations.safeAdminPath("/admin/validation/run"), "");
});

test("preferenceKey is scoped to UI-only preference names", () => {
  assert.equal(
    operations.preferenceKey("/admin/operations/feed-health", "filter"),
    "open-transit-rt.operations._admin_operations_feed-health.filter"
  );
  assert.equal(operations.preferenceKey("/admin/operations/feed-health", "csrf_token"), "");
  assert.equal(operations.preferenceKey("/admin/operations/feed-health", "bearer"), "");
  assert.equal(operations.preferenceKey("/admin/operations/feed-health", "device_token"), "");
  assert.equal(operations.preferenceKey("/public/operations", "filter"), "");
});
