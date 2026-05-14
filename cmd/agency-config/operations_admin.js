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

  function rowMatchesFilterText(text, status, filter, search) {
    var haystack = String((status || "") + " " + (text || "")).toLowerCase();
    var selected = String(filter || "all").toLowerCase();
    var query = String(search || "").trim().toLowerCase();
    if (query && haystack.indexOf(query) === -1) {
      return false;
    }
    if (selected === "all") {
      return true;
    }
    if (selected === "needs_action") {
      return /(needs|blocked|missing|stale|unknown|failed|warning|not available|not run)/.test(haystack);
    }
    if (selected === "fresh") {
      return /(fresh|ok|recorded|configured|passed|current|available)/.test(haystack);
    }
    return haystack.indexOf(selected.replace("_", " ")) !== -1 || haystack.indexOf(selected) !== -1;
  }

  function needsActionScore(text, status) {
    return rowMatchesFilterText(text, status, "needs_action", "") ? 0 : 1;
  }

  function sortRowModels(rows, sort) {
    var mode = String(sort || "needs_action").toLowerCase();
    return rows.map(function (row, index) {
      return { row: row, index: index };
    }).sort(function (a, b) {
      var aName = String(a.row.name || "").toLowerCase();
      var bName = String(b.row.name || "").toLowerCase();
      if (mode === "name") {
        return aName.localeCompare(bName) || a.index - b.index;
      }
      if (mode === "status") {
        var statusOrder = String(a.row.status || "").localeCompare(String(b.row.status || ""));
        return statusOrder || aName.localeCompare(bName) || a.index - b.index;
      }
      var score = needsActionScore(a.row.text, a.row.status) - needsActionScore(b.row.text, b.row.status);
      return score || aName.localeCompare(bName) || a.index - b.index;
    }).map(function (entry) {
      return entry.row;
    });
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

  function enhanceCopyValues(documentRef) {
    var values = documentRef.querySelectorAll(".copy-value[data-copy-value]");
    values.forEach(function (value, index) {
      if (value.getAttribute("data-copy-ready") === "true") {
        return;
      }
      value.setAttribute("data-copy-ready", "true");
      var button = documentRef.createElement("button");
      var status = documentRef.createElement("span");
      var statusID = "copy-status-" + index;
      button.type = "button";
      button.className = "copy-action";
      button.textContent = "Copy configured URL";
      button.setAttribute("aria-describedby", statusID);
      status.id = statusID;
      status.className = "review-status";
      status.setAttribute("aria-live", "polite");
      status.textContent = "Configured value is visible for manual copy.";
      value.insertAdjacentElement("afterend", status);
      value.insertAdjacentElement("afterend", button);
      button.addEventListener("click", function () {
        var text = value.getAttribute("data-copy-value") || value.textContent || "";
        if (!text) {
          status.textContent = "This value is not marked safe to copy.";
          return;
        }
        var write = root.navigator && root.navigator.clipboard && root.navigator.clipboard.writeText
          ? root.navigator.clipboard.writeText(text)
          : Promise.reject(new Error("clipboard unavailable"));
        write.then(function () {
          status.textContent = "Copied configured value.";
        }).catch(function () {
          status.textContent = "Copy unavailable. Select the visible configured value manually.";
        });
      });
    });
  }

  function enhanceReviewTools(documentRef) {
    var tools = documentRef.querySelectorAll("[data-review-tools]");
    tools.forEach(function (tool) {
      var targetID = tool.getAttribute("data-review-target");
      var target = targetID ? documentRef.getElementById(targetID) : null;
      if (!target) {
        return;
      }
      var rows = Array.prototype.slice.call(target.querySelectorAll("[data-review-row]"));
      var filter = tool.querySelector("[data-review-filter]");
      var search = tool.querySelector("[data-review-search]");
      var sort = tool.querySelector("[data-review-sort]");
      var reset = tool.querySelector("[data-review-reset]");
      var remember = tool.querySelector("[data-review-remember]");
      var status = tool.querySelector("[data-review-status]");
      var scopePath = root.location && root.location.pathname ? root.location.pathname : "/admin/operations";
      var storedFilter = readPreference(scopePath, "filter");
      var storedSort = readPreference(scopePath, "sort");
      if (storedFilter && filter) {
        filter.value = storedFilter;
        if (remember) {
          remember.checked = true;
        }
      }
      if (storedSort && sort) {
        sort.value = storedSort;
        if (remember) {
          remember.checked = true;
        }
      }
      function applyReviewState() {
        var filterValue = filter ? filter.value : "all";
        var searchValue = search ? search.value : "";
        var sortValue = sort ? sort.value : "needs_action";
        var models = rows.map(function (row) {
          return {
            element: row,
            name: row.getAttribute("data-review-name") || "",
            status: row.getAttribute("data-review-status") || "",
            text: row.textContent || ""
          };
        });
        sortRowModels(models, sortValue).forEach(function (model) {
          target.appendChild(model.element);
        });
        var shown = 0;
        rows.forEach(function (row) {
          var show = rowMatchesFilterText(row.textContent || "", row.getAttribute("data-review-status") || "", filterValue, searchValue);
          row.hidden = !show;
          if (show) {
            shown += 1;
          }
        });
        if (status) {
          status.textContent = "Showing " + shown + " of " + rows.length + " private diagnostic rows.";
        }
        if (remember && remember.checked) {
          writePreference(scopePath, "filter", filterValue);
          writePreference(scopePath, "sort", sortValue);
        }
      }
      [filter, search, sort].forEach(function (control) {
        if (control) {
          control.addEventListener("input", applyReviewState);
          control.addEventListener("change", applyReviewState);
        }
      });
      if (reset) {
        reset.addEventListener("click", function () {
          if (filter) {
            filter.value = "all";
          }
          if (search) {
            search.value = "";
          }
          if (sort) {
            sort.value = "needs_action";
          }
          applyReviewState();
        });
      }
      applyReviewState();
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
    enhanceCopyValues(documentRef);
    enhanceReviewTools(documentRef);
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
    rowMatchesFilterText: rowMatchesFilterText,
    safeAdminPath: safeAdminPath,
    sortRowModels: sortRowModels,
    writePreference: writePreference
  };

  onReady(function () {
    init(root.document);
  });

  return api;
});
