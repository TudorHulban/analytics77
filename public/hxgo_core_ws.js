document.addEventListener('DOMContentLoaded', (event) => {
    const CONFIG = Object.freeze({
        WS_REQUEST_TIMEOUT: 30000,
        WS_HEARTBEAT_INTERVAL: 30000,
        WS_RECONNECT_MAX_BACKOFF: 30000,
        WS_RECONNECT_JITTER_MAX: 5000,
        MS_FILE_UPLOAD_RESET_TIMEOUT: 1500,
        CACHE_MAX_ENTRIES: 100,
        CACHE_DEFAULT_TTL: 300
    });

    const HX = Object.freeze({
        GET: 'hx-get',
        POST: 'hx-post',
        UPLOAD: 'hx-upload',
        REDIRECT: 'hx-redirect',
        SWAP: 'hx-swap',
        SEND: 'hx-send',
        REQUIRE: 'hx-require',
        MIN: 'hx-min',
        MAX: 'hx-max',
        VDISABLE: 'hx-vdisable',
        PDISABLE: 'hx-pdisable',
        SHOW: 'hx-show',
        HIDE: 'hx-hide',
        SHOW_ONLOAD: 'hx-show-onload',
        ONCHANGE_ENABLE: 'hx-onchange-enable'
    });

    const VERSION = "0.1.3-core-ws";
    console.log("version %s", VERSION);

    // --- WebSocket State ---
    const WS_URL = (window.location.protocol === 'https:' ? 'wss://' : 'ws://') + window.location.host + '/ws';
    let ws = null;
    let wsBackoff = 1000;
    let wsReconnectTimer = null;
    let wsReconnectAttempts = 0;
    let heartbeatInterval = null;
    let pendingRequests = new Map();
    let requestIdCounter = 0;

    // --- Event System ---
    const fireEvent = (name, detail = {}, options = {}) => {
        const event = new CustomEvent(`hx:${name}`, {
            detail,
            cancelable: options.cancelable !== false,
            bubbles: true
        });
        return document.dispatchEvent(event);
    };

    // --- Auth Token Management ---
    function getAuthToken() {
        const metaToken = document.querySelector('meta[name="auth-token"]');
        if (metaToken && metaToken.content) return metaToken.content;
        const match = document.cookie.match(/(?:^|;\s*)auth_token=([^;]+)/);
        return match ? match[1] : null;
    }

    function getCSRFToken() {
        const metaToken = document.querySelector('meta[name="csrf-token"]');
        if (metaToken && metaToken.content) return metaToken.content;
        const match = document.cookie.match(/(?:^|;\s*)csrf_token=([^;]+)/);
        return match ? match[1] : null;
    }

    // --- Request Management ---
    function retireRequest(id, responseData = null, extra = {}) {
        if (!id || !pendingRequests.has(id)) return;
        const entry = pendingRequests.get(id);
        clearTimeout(entry.timerId);
        pendingRequests.delete(id);

        fireEvent('afterResponse', {
            requestId: id,
            pendingCount: pendingRequests.size,
            response: responseData,
            element: entry.element,
            endpoint: entry.endpoint,
            method: entry.method,
            ...extra
        });
    }

    // --- DOM Swap Helper ---
    // Single, consistent beforeSwap / afterSwap per element.
    // requestId is optional (present for live WS responses, null for cache hits).
    function swapElementById(id, outerHTML, requestId = null) {
        const target = document.getElementById(id);
        if (!target) return;

        fireEvent('beforeSwap', { id, target, content: outerHTML, requestId });
        target.outerHTML = outerHTML;
        const fresh = document.getElementById(id);

        if (fresh) {
            fireEvent('afterSwap', { id, target: fresh, requestId });
            // Let plugins handle listener reattachment
            fireEvent('reattachListeners', { element: fresh });
        }
    }

    // --- WebSocket Core ---
    function initWs() {
        if (ws?.readyState === WebSocket.CONNECTING || ws?.readyState === WebSocket.OPEN) return;

        // Always re-read the current auth token (supports token refresh)
        const authToken = getAuthToken();
        let wsUrl = WS_URL;
        if (authToken) {
            wsUrl += (wsUrl.includes('?') ? '&' : '?') + 'token=' + encodeURIComponent(authToken);
        }

        try {
            ws = new WebSocket(wsUrl);
        } catch (e) {
            fireEvent('connectionError', { error: e });
            scheduleReconnect();
            return;
        }

        ws.onopen = () => {
            wsBackoff = 1000;
            wsReconnectAttempts = 0;
            fireEvent('connected');

            if (heartbeatInterval) clearInterval(heartbeatInterval);
            heartbeatInterval = setInterval(() => {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    ws.send('ping');
                }
            }, CONFIG.WS_HEARTBEAT_INTERVAL);
        };

        ws.onmessage = (event) => {
            let html = event.data;

            if (html === 'pong') return;

            let responseId = null;
            const commentMatch = html.match(/^<!--\s*_hx_req_id:\s*(\S+)\s*-->\s*/);
            if (commentMatch) {
                responseId = commentMatch[1];
                html = html.slice(commentMatch[0].length);
            }

            if (html.startsWith('unknown endpoint:')) {
                retireRequest(responseId, null, { error: html });
                fireEvent('error', { message: html, requestId: responseId });
                return;
            }

            const temp = document.createElement('div');
            temp.innerHTML = html;

            // Check for redirect
            const redirectElement = temp.querySelector('[hx-redirect]');
            if (redirectElement) {
                const redirectUrl = redirectElement.getAttribute('hx-redirect');

                // Fire redirect first so plugins can cancel before we retire
                const redirectEvent = fireEvent('redirect', {
                    url: redirectUrl,
                    requestId: responseId,
                    response: html,
                    cancelable: true
                });

                // Always retire so pendingCount stays accurate
                retireRequest(responseId, html, { redirected: true, redirectUrl });

                if (!redirectEvent.defaultPrevented) {
                    window.location.href = redirectUrl;
                }
                return;
            }

            const elements = temp.querySelectorAll('[id]');
            if (elements.length === 0) {
                // Empty / non-swap response – still retire cleanly
                retireRequest(responseId, html);
                return;
            }

            // No batch beforeSwap – only the per-element event inside swapElementById
            elements.forEach(el => swapElementById(el.id, el.outerHTML, responseId));

            retireRequest(responseId, html);
        };

        ws.onerror = (err) => {
            fireEvent('error', { error: err });
        };

        ws.onclose = (event) => {
            // Notify every in-flight request so listeners (loading indicator, etc.) can clean up
            const pendingSnapshot = Array.from(pendingRequests.entries());
            pendingRequests.clear();

            pendingSnapshot.forEach(([id, entry]) => {
                clearTimeout(entry.timerId);
                fireEvent('timeout', {
                    requestId: id,
                    element: entry.element,
                    endpoint: entry.endpoint,
                    method: entry.method,
                    reason: 'connection_closed',
                    code: event.code,
                    closeReason: event.reason
                });
                fireEvent('afterResponse', {
                    requestId: id,
                    pendingCount: 0,
                    response: null,
                    element: entry.element,
                    endpoint: entry.endpoint,
                    method: entry.method,
                    connectionClosed: true
                });
            });

            fireEvent('disconnected', { code: event.code, reason: event.reason });

            if (heartbeatInterval) {
                clearInterval(heartbeatInterval);
                heartbeatInterval = null;
            }
            scheduleReconnect();
        };
    }

    function scheduleReconnect() {
        if (wsReconnectTimer) return;

        const baseDelay = Math.min(wsBackoff, CONFIG.WS_RECONNECT_MAX_BACKOFF);
        const jitter = Math.random() * CONFIG.WS_RECONNECT_JITTER_MAX;
        const delay = baseDelay + jitter;

        fireEvent('reconnecting', { delay, attempt: wsReconnectAttempts + 1 });

        wsReconnectTimer = setTimeout(() => {
            wsReconnectTimer = null;
            wsReconnectAttempts++;
            initWs();
        }, delay);
        wsBackoff = Math.min(wsBackoff * 2, CONFIG.WS_RECONNECT_MAX_BACKOFF);
    }

    /**
     * Force a clean reconnect (useful after auth-token refresh).
     * Closes the current socket (if any) and immediately starts a new connection
     * with the latest token from meta/cookie.
     */
    function reconnect() {
        if (wsReconnectTimer) {
            clearTimeout(wsReconnectTimer);
            wsReconnectTimer = null;
        }
        if (heartbeatInterval) {
            clearInterval(heartbeatInterval);
            heartbeatInterval = null;
        }
        if (ws) {
            // Avoid the normal onclose reconnect logic racing with us
            const oldWs = ws;
            ws = null;
            try {
                oldWs.onopen = null;
                oldWs.onclose = null;
                oldWs.close(1000, 'Manual reconnect');
            } catch (_) { /* ignore */ }
        }
        wsBackoff = 1000;
        wsReconnectAttempts = 0;
        initWs();
    }

    // --- Core Request Sending ---
    const sendWsAction = (element, form) => {
        const id = 'hx' + (++requestIdCounter);

        // Fire before request - plugins can cancel
        const beforeEvent = fireEvent('beforeRequest', {
            element,
            form,
            requestId: id,
            cancelable: true
        });

        if (beforeEvent.defaultPrevented) {
            fireEvent('requestCancelled', {
                element,
                requestId: id
            });

            // Also fire afterResponse so UI knows to clean up
            fireEvent('afterResponse', {
                requestId: id,
                pendingCount: pendingRequests.size,
                response: null,
                element,
                endpoint: element.getAttribute(HX.GET) || element.getAttribute(HX.POST),
                method: element.hasAttribute(HX.POST) ? 'POST' : 'GET',
                cancelled: true  // ← Flag to indicate cancellation
            });

            return;// do not proceed with actual request
        }

        if (!ws || ws.readyState !== WebSocket.OPEN) {
            fireEvent('connectionError', { element });
            return;
        }

        const endpoint = element.getAttribute(HX.GET) || element.getAttribute(HX.POST);
        const isPost = element.hasAttribute(HX.POST);
        const params = new URLSearchParams();

        // Collect form data + hx-send values for BOTH GET and POST
        const formData = new FormData(form);

        const hxSend = element.getAttribute(HX.SEND);
        if (hxSend) {
            hxSend.split(',').map(id => id.trim()).forEach(id => {
                const el = document.querySelector(id);
                if (!el) return;
                const value = ['INPUT', 'TEXTAREA', 'SELECT'].includes(el.tagName)
                    ? el.value
                    : (el.innerText || el.textContent);
                formData.append(el.id || id, value);
            });
        }

        for (const [key, value] of formData.entries()) {
            if (value instanceof File) continue;
            params.append(key, value);
        }

        const csrf = getCSRFToken();
        if (csrf) params.append('_csrf', csrf);

        const timerId = setTimeout(() => {
            if (pendingRequests.has(id)) {
                const entry = pendingRequests.get(id);
                pendingRequests.delete(id);
                clearTimeout(entry.timerId);

                fireEvent('timeout', {
                    requestId: id,
                    element: entry.element,
                    endpoint: entry.endpoint,
                    method: entry.method
                });

                // Always emit afterResponse so UI plugins can hide spinners etc.
                fireEvent('afterResponse', {
                    requestId: id,
                    pendingCount: pendingRequests.size,
                    response: null,
                    element: entry.element,
                    endpoint: entry.endpoint,
                    method: entry.method,
                    timedOut: true
                });
            }
        }, CONFIG.WS_REQUEST_TIMEOUT);

        pendingRequests.set(id, {
            timerId,
            element,
            endpoint,
            method: isPost ? 'POST' : 'GET'
        });
        params.append('_hx_req_id', id);

        const wire = `${isPost ? 'POST' : 'GET'} ${endpoint}\n${params.toString()}`;

        try {
            ws.send(wire);
            fireEvent('requestSent', {
                requestId: id,
                element,
                endpoint,
                method: isPost ? 'POST' : 'GET'
            });
            // Emit afterRequest so cache invalidation (and other plugins) can react
            fireEvent('afterRequest', {
                requestId: id,
                element,
                endpoint,
                method: isPost ? 'POST' : 'GET'
            });
        } catch (e) {
            if (pendingRequests.has(id)) {
                clearTimeout(pendingRequests.get(id).timerId);
                pendingRequests.delete(id);
            }
            fireEvent('error', { error: e, element });
            // Ensure afterResponse still fires for consistency
            fireEvent('afterResponse', {
                requestId: id,
                pendingCount: pendingRequests.size,
                response: null,
                element,
                endpoint,
                method: isPost ? 'POST' : 'GET',
                error: e
            });
        }
    };

    // --- Core Upload (HTTP) ---
    const handleUpload = (element, form, file = null) => {
        const beforeEvent = fireEvent('beforeUpload', {
            element,
            form,
            file,
            cancelable: true
        });

        if (beforeEvent.defaultPrevented) {
            fireEvent('uploadCancelled', { element });
            return;
        }

        const endpoint = element.getAttribute(HX.UPLOAD);
        const targetSelectors = element.getAttribute(HX.SWAP);
        const targetElements = targetSelectors
            ? targetSelectors.split(',').map(selector => document.querySelector(selector.trim()))
            : [];
        const redirectUrl = element.getAttribute(HX.REDIRECT);

        if (!file) return;

        const formData = new FormData(form);
        formData.append('file', file);

        const csrfToken = getCSRFToken();
        if (csrfToken) formData.append('_csrf', csrfToken);

        const xhr = new XMLHttpRequest();
        xhr.open('POST', endpoint, true);

        xhr.upload.addEventListener('progress', (event) => {
            if (event.lengthComputable) {
                const percentComplete = (event.loaded / event.total) * 100;
                fireEvent('uploadProgress', {
                    element,
                    progress: percentComplete
                });
            }
        });

        xhr.onload = () => {
            fireEvent('uploadComplete', { element, status: xhr.status });

            if (xhr.status >= 200 && xhr.status < 300) {
                if (redirectUrl) {
                    fireEvent('redirect', { url: redirectUrl, element });
                } else {
                    const data = xhr.responseText;
                    if (data) {
                        const extractedHTML = parseString(data, targetElements);
                        targetElements.forEach(targetElement => {
                            if (!targetElement) return;
                            const responseElement = extractedHTML.get(targetElement.id);
                            if (responseElement) {
                                swapElementById(targetElement.id, responseElement);
                            }
                        });
                    }
                }
            } else {
                fireEvent('error', {
                    message: xhr.responseText,
                    element,
                    status: xhr.status
                });
            }
        };

        xhr.onerror = () => {
            fireEvent('error', {
                error: xhr.statusText,
                element
            });
        };

        fireEvent('uploadStarted', { element, file });
        xhr.send(formData);
    };

    // --- Event Handlers (Core) ---
    const handleButtonClick = (event) => {
        event.preventDefault();
        const element = event.currentTarget;
        const form = element.closest('form') || document.createElement('form');

        if (element.hasAttribute(HX.UPLOAD)) {
            handleUploadClick(event);
        } else {
            sendWsAction(element, form);
        }
    };

    const handleSelectChange = (event) => {
        const element = event.currentTarget;
        const form = document.createElement('form');
        const hiddenInput = document.createElement('input');
        hiddenInput.type = 'hidden';
        hiddenInput.name = 'selectedOption';
        hiddenInput.value = element.value;
        form.appendChild(hiddenInput);
        sendWsAction(element, form);
    };

    const handleUploadClick = (event) => {
        event.preventDefault();
        const element = event.currentTarget;
        const form = element.closest('form') || document.createElement('form');

        let fileInput = form.querySelector('input[type="file"]');
        if (!fileInput) {
            fileInput = document.createElement('input');
            fileInput.type = 'file';
            fileInput.name = 'file';
            fileInput.style.display = 'none';
            form.appendChild(fileInput);
        }

        fileInput.value = '';
        fileInput.addEventListener('change', function onChange() {
            const fileUploaded = fileInput.files[0];
            handleUpload(element, form, fileUploaded);
            setTimeout(() => {
                fileInput.value = '';
            }, CONFIG.MS_FILE_UPLOAD_RESET_TIMEOUT);
        }, { once: true });

        fileInput.click();
    };

    // --- Utility Functions ---
    function parseString(inputString, elements) {
        let parts = inputString.split('|');
        let matchedElementsMap = new Map();

        for (let part of parts) {
            for (let element of elements) {
                if (element && element.id && part.includes(`id="${element.id}"`)) {
                    matchedElementsMap.set(element.id, part);
                    break;
                }
            }
        }
        return matchedElementsMap;
    }

    // --- Event Delegation for Dynamic Elements ---
    // This ensures buttons and selects work even after DOM swaps
    document.addEventListener('click', (event) => {
        const element = event.target.closest(
            `button[${HX.GET}], button[${HX.POST}], button[${HX.UPLOAD}], a[${HX.GET}], a[${HX.POST}], a[${HX.UPLOAD}]`
        );
        if (element) {
            event.preventDefault();
            handleButtonClick({ currentTarget: element, preventDefault: () => { } });
        }
    });

    document.addEventListener('change', (event) => {
        const element = event.target.closest(`select[${HX.GET}], select[${HX.POST}]`);
        if (element) {
            handleSelectChange({ currentTarget: element });
        }
    });

    // --- Lifecycle Events ---
    window.addEventListener('beforeunload', () => {
        fireEvent('beforeUnload');

        if (heartbeatInterval) {
            clearInterval(heartbeatInterval);
            heartbeatInterval = null;
        }
        if (ws) {
            ws.close(1000, 'Page unloading');
            ws = null;
        }
        if (wsReconnectTimer) {
            clearTimeout(wsReconnectTimer);
            wsReconnectTimer = null;
        }
    });

    document.addEventListener('visibilitychange', () => {
        if (document.visibilityState === 'visible') {
            if (!ws || ws.readyState === WebSocket.CLOSED || ws.readyState === WebSocket.CLOSING) {
                fireEvent('visibilityRestored');
                initWs();
            }
        }
    });

    window.addEventListener('online', () => {
        fireEvent('online');
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            initWs();
        }
    });

    window.addEventListener('offline', () => {
        fireEvent('offline');
    });

    // --- Expose API for plugins ---
    window.hx = {
        version: VERSION,
        fireEvent,
        sendWsAction,
        handleUpload,
        swapElementById,
        getAuthToken,
        getCSRFToken,
        reconnect,          // force reconnect with fresh token
        config: CONFIG,
        constants: HX
    };

    // Initialize
    fireEvent('init');
    initWs();
});