(function (root, factory) {
  "use strict";
  var api = factory(root);
  if (typeof module !== "undefined" && module.exports) {
    module.exports = api;
  } else {
    root.OpenTransitOperations = api;
  }
})(typeof globalThis !== "undefined" ? globalThis : this, function (root) {
  "use strict";

  var preferenceNames = {
    filter: true,
    sort: true,
    compact: true,
    refreshInterval: true
  };
  var terminalStatuses = {
    ok: true,
    needs_review: true,
    blocked: true,
    failed: true,
    complete: true
  };

  function currentOrigin() {
    if (root.location && root.location.origin) {
      return root.location.origin;
    }
    return ["http", "://", "localhost"].join("");
  }

  function safeAdminPath(path, options) {
    var method = options && options.method ? String(options.method).toUpperCase() : "GET";
    if (!path || typeof path !== "string") {
      return "";
    }
    if (/^[a-z][a-z0-9+.-]*:/i.test(path) || path.indexOf("//") === 0) {
      return "";
    }
    var url;
    try {
      url = new URL(path, currentOrigin());
    } catch (_err) {
      return "";
    }
    if (url.origin !== currentOrigin() || url.searchParams.has("agency_id")) {
      return "";
    }
    if (method === "POST" && url.pathname === "/admin/operations/validation-health/refresh.json") {
      return url.pathname;
    }
    if (method !== "GET") {
      return "";
    }
    if (url.pathname === "/admin/operations.json") {
      return url.pathname;
    }
    if (/^\/admin\/operations\/[a-z0-9-]+\.json$/.test(url.pathname)) {
      return url.pathname;
    }
    return "";
  }

  function preferenceKey(pathname, name) {
    if (!preferenceNames[name]) {
      return "";
    }
    var path = pathname || "/admin/operations";
    if (path.indexOf("/admin/operations") !== 0 || path.indexOf("..") !== -1) {
      return "";
    }
    return "open-transit-rt.operations." + path.replace(/[^a-z0-9_-]+/gi, "_") + "." + name;
  }

  function storageAvailable() {
    try {
      if (!root.localStorage) {
        return false;
      }
      var key = "open-transit-rt.operations.__probe__";
      root.localStorage.setItem(key, "1");
      root.localStorage.removeItem(key);
      return true;
    } catch (_err) {
      return false;
    }
  }

  function readPreference(pathname, name) {
    if (!storageAvailable()) {
      return "";
    }
    var key = preferenceKey(pathname, name);
    if (!key) {
      return "";
    }
    var value = root.localStorage.getItem(key);
    return value == null ? "" : value;
  }

  function writePreference(pathname, name, value) {
    if (!storageAvailable()) {
      return false;
    }
    var key = preferenceKey(pathname, name);
    if (!key) {
      return false;
    }
    root.localStorage.setItem(key, String(value).slice(0, 80));
    return true;
  }

  function boundedRefreshInterval(value) {
    var parsed = Number.parseInt(value, 10);
    if (!Number.isFinite(parsed) || parsed <= 0) {
      return 0;
    }
    if (parsed < 15) {
      return 15;
    }
    if (parsed > 300) {
      return 300;
    }
    return parsed;
  }

  function isTerminalStatus(status) {
    return Boolean(terminalStatuses[String(status || "").toLowerCase()]);
  }

  function commandRefreshBody(csrfToken) {
    var body = new URLSearchParams();
    body.set("action", "refresh");
    if (csrfToken) {
      body.set("csrf_token", csrfToken);
    }
    return body;
  }

  function commandResultText(result) {
    if (!result || typeof result !== "object") {
      return "Private diagnostic refresh did not return a command result.";
    }
    var status = result.status ? String(result.status) : "unknown";
    var summary = result.summary ? String(result.summary) : "No private summary was returned.";
    return "Private diagnostic refresh returned " + status + ": " + summary;
  }

  function csrfTokenFromDocument(documentRef) {
    var input = documentRef ? documentRef.querySelector('input[name="csrf_token"]') : null;
    return input ? input.value : "";
  }

  function setLiveMessage(region, message) {
    if (!region) {
      return;
    }
    region.textContent = message;
  }

  function postPrivateCommand(path, fields, fetchImpl) {
    var endpoint = safeAdminPath(path, { method: "POST" });
    if (!endpoint) {
      return Promise.reject(new Error("unsupported private command endpoint"));
    }
    var body = commandRefreshBody(fields && fields.csrfToken ? fields.csrfToken : "");
    var fetcher = fetchImpl || root.fetch;
    if (typeof fetcher !== "function") {
      return Promise.reject(new Error("fetch unavailable"));
    }
    return fetcher(endpoint, {
      method: "POST",
      credentials: "same-origin",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/x-www-form-urlencoded"
      },
      body: body
    }).then(function (response) {
      if (!response || !response.ok) {
        throw new Error("private command request failed");
      }
      return response.json();
    });
  }

  function enhanceRefreshButtons(documentRef) {
    var buttons = documentRef.querySelectorAll("[data-admin-refresh]");
    buttons.forEach(function (button) {
      button.addEventListener("click", function () {
        var endpoint = button.getAttribute("data-admin-refresh") || "";
        var region = documentRef.getElementById(button.getAttribute("aria-describedby") || "");
        var csrfToken = csrfTokenFromDocument(documentRef);
        button.disabled = true;
        button.setAttribute("aria-busy", "true");
        setLiveMessage(region, "Refreshing private diagnostic summary.");
        postPrivateCommand(endpoint, { csrfToken: csrfToken }).then(function (result) {
          setLiveMessage(region, commandResultText(result));
        }).catch(function () {
          setLiveMessage(region, "Refresh did not complete. Review the private page state and retry.");
        }).finally(function () {
          button.disabled = false;
          button.removeAttribute("aria-busy");
          button.focus();
        });
      });
    });
  }

  function onReady(callback) {
    if (!root.document) {
      return;
    }
    if (root.document.readyState === "loading") {
      root.document.addEventListener("DOMContentLoaded", callback, { once: true });
    } else {
      callback();
    }
  }

  function init(doc) {
    var documentRef = doc || root.document;
    if (!documentRef || !documentRef.documentElement) {
      return false;
    }
    documentRef.documentElement.setAttribute("data-operations-enhanced", "true");
    enhanceRefreshButtons(documentRef);
    return true;
  }

  var api = {
    boundedRefreshInterval: boundedRefreshInterval,
    commandRefreshBody: commandRefreshBody,
    commandResultText: commandResultText,
    init: init,
    isTerminalStatus: isTerminalStatus,
    onReady: onReady,
    postPrivateCommand: postPrivateCommand,
    preferenceKey: preferenceKey,
    readPreference: readPreference,
    safeAdminPath: safeAdminPath,
    writePreference: writePreference
  };

  onReady(function () {
    init(root.document);
  });

  return api;
});
