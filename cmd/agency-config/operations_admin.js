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
    return true;
  }

  var api = {
    init: init,
    onReady: onReady,
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
