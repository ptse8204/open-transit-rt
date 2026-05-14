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

test("command refresh body posts explicit refresh action and form CSRF", () => {
  const body = operations.commandRefreshBody("token-value");
  assert.equal(body.get("action"), "refresh");
  assert.equal(body.get("csrf_token"), "token-value");
  assert.equal(body.has("validator_id"), false);
  assert.equal(body.has("path"), false);
});

test("private command fetch is same-origin and form-encoded", async () => {
  const calls = [];
  const result = await operations.postPrivateCommand(
    "/admin/operations/validation-health/refresh.json",
    { csrfToken: "token-value" },
    async (url, init) => {
      calls.push({ url, init });
      return { ok: true, json: async () => ({ status: "ok", summary: "Local records refreshed." }) };
    }
  );
  assert.equal(result.status, "ok");
  assert.equal(calls.length, 1);
  assert.equal(calls[0].url, "/admin/operations/validation-health/refresh.json");
  assert.equal(calls[0].init.method, "POST");
  assert.equal(calls[0].init.credentials, "same-origin");
  assert.equal(calls[0].init.headers.Accept, "application/json");
  assert.equal(calls[0].init.headers["Content-Type"], "application/x-www-form-urlencoded");
  assert.equal(calls[0].init.body.get("action"), "refresh");
  assert.equal(calls[0].init.body.get("csrf_token"), "token-value");

  await assert.rejects(
    () => operations.postPrivateCommand("/public/gtfsrt/vehicle_positions.json", {}, async () => ({ ok: true })),
    /unsupported private command endpoint/
  );
});

test("status and refresh interval helpers stay bounded", () => {
  assert.equal(operations.isTerminalStatus("ok"), true);
  assert.equal(operations.isTerminalStatus("needs_review"), true);
  assert.equal(operations.isTerminalStatus("running"), false);
  assert.equal(operations.boundedRefreshInterval("0"), 0);
  assert.equal(operations.boundedRefreshInterval("4"), 15);
  assert.equal(operations.boundedRefreshInterval("120"), 120);
  assert.equal(operations.boundedRefreshInterval("999"), 300);
});

test("command result text stays diagnostic", () => {
  const text = operations.commandResultText({ status: "needs_review", summary: "Validator records are stale." });
  assert.match(text, /Private diagnostic refresh returned needs_review/);
  assert.doesNotMatch(text.toLowerCase(), /compliant|accepted|production ready|certified/);
});

test("review row filtering keeps local diagnostic language", () => {
  assert.equal(operations.rowMatchesFilterText("Vehicle Positions stale next action", "stale", "needs_action", ""), true);
  assert.equal(operations.rowMatchesFilterText("Schedule configured current", "ok", "needs_action", ""), false);
  assert.equal(operations.rowMatchesFilterText("Alerts missing artifact", "missing", "missing", "alerts"), true);
  assert.equal(operations.rowMatchesFilterText("Alerts missing artifact", "missing", "missing", "vehicle"), false);
});

test("review row sorting puts needs-action rows first", () => {
  const rows = operations.sortRowModels([
    { name: "Schedule", status: "ok", text: "configured current" },
    { name: "Alerts", status: "missing", text: "missing artifact" },
    { name: "Vehicle Positions", status: "stale", text: "stale diagnostic" }
  ], "needs_action");
  assert.deepEqual(rows.map((row) => row.name), ["Alerts", "Vehicle Positions", "Schedule"]);
  const byName = operations.sortRowModels(rows, "name");
  assert.deepEqual(byName.map((row) => row.name), ["Alerts", "Schedule", "Vehicle Positions"]);
});
