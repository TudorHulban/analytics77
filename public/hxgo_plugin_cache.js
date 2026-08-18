/**
 * Caching plugin for hxgo
 * Adds client-side caching with Time-To-Live (TTL) in seconds
 *
 * Usage:
 *   <div hx-get="/api/data" hx-cache="true" hx-cache-ttl="60">
 *   <div hx-get="/api/user" hx-cache="true" hx-cache-ttl="300">
 *   <button hx-post="/refresh" hx-cache-invalidate="/api/data">
 *   <button hx-post="/refresh" hx-cache-invalidate="/api/data,/api/user">
 */

document.addEventListener('DOMContentLoaded', () => {
    class CacheManager {
        constructor(maxEntries = 100, defaultTTL = 300) {
            this.cache = new Map();
            this.maxEntries = maxEntries;
            this.defaultTTL = defaultTTL;
        }

        get(key) {
            const entry = this.cache.get(key);
            if (!entry) return null;

            if (Date.now() > entry.expiresAt) {
                this.cache.delete(key);
                return null;
            }

            return entry.value;
        }

        set(key, value, ttlSeconds = this.defaultTTL) {
            // FIFO eviction (insertion order), not LRU: the oldest-inserted
            // entry goes first, regardless of how recently it was read.
            if (this.cache.size >= this.maxEntries) {
                const oldestKey = this.cache.keys().next().value;
                this.cache.delete(oldestKey);
            }

            this.cache.set(key, {
                value,
                expiresAt: Date.now() + (ttlSeconds * 1000)
            });
        }

        invalidate(key) {
            this.cache.delete(key);
        }

        invalidatePrefix(prefix) {
            for (const key of this.cache.keys()) {
                if (key.startsWith(prefix)) {
                    this.cache.delete(key);
                }
            }
        }

        clear() {
            this.cache.clear();
        }
    }

    // window.hx is guaranteed to exist by this point: core's own setup also
    // runs inside DOMContentLoaded, and script tags fire that event only
    // after every synchronous <script> has finished executing.
    const cacheManager = new CacheManager(
        window.hx.config?.CACHE_MAX_ENTRIES ?? 100,
        window.hx.config?.CACHE_DEFAULT_TTL ?? 300
    );

    // ---------- Helpers ----------

    function collectFormParams(element, form) {
        const formData = new FormData(form);

        // Mirror core's hx-send handling so cache keys stay consistent
        const hxSend = element.getAttribute(window.hx.constants.SEND || 'hx-send');
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

        const params = new URLSearchParams();
        for (const [key, value] of formData.entries()) {
            if (value instanceof File) continue;
            params.append(key, value);
        }
        return params;
    }

    function buildCacheKey(method, endpoint, element, form) {
        const params = collectFormParams(element, form);
        return `${method}:${endpoint}:${params.toString()}`;
    }

    function applyCachedResponse(html) {
        const temp = document.createElement('div');
        temp.innerHTML = html;

        const redirectElement = temp.querySelector('[hx-redirect]');
        if (redirectElement) {
            const redirectUrl = redirectElement.getAttribute('hx-redirect');
            const redirectEvent = window.hx.fireEvent('redirect', {
                url: redirectUrl,
                fromCache: true,
                cancelable: true
            });
            if (!redirectEvent.defaultPrevented) {
                window.location.href = redirectUrl;
            }
            return;
        }

        const elements = temp.querySelectorAll('[id]');
        if (elements.length === 0) return;

        // requestId = null for cache hits (same as core path)
        elements.forEach(el => {
            window.hx.swapElementById(el.id, el.outerHTML, null);
        });
    }

    // ---------- Event handlers ----------

    document.addEventListener('hx:beforeRequest', (e) => {
        const { element, form, requestId } = e.detail;
        if (!element) return;

        const shouldCache = element.getAttribute('hx-cache') === 'true';
        if (!shouldCache) return;

        const endpoint = element.getAttribute('hx-get') || element.getAttribute('hx-post');
        if (!endpoint) return;

        const method = element.hasAttribute('hx-post') ? 'POST' : 'GET';
        const cacheKey = buildCacheKey(method, endpoint, element, form || document.createElement('form'));

        const cachedResponse = cacheManager.get(cacheKey);
        if (!cachedResponse) return;

        // Fire a proper cacheHit event
        window.hx.fireEvent('cacheHit', {
            cacheKey,
            response: cachedResponse,
            element,
            endpoint,
            method,
            requestId
        });

        // Apply using the same swap path as the core
        applyCachedResponse(cachedResponse);

        // Synthesize afterResponse so UI plugin hides the loading indicator
        // and any other afterResponse listeners still run
        window.hx.fireEvent('afterResponse', {
            requestId: requestId,
            pendingCount: 0,          // critical for loading indicator
            response: cachedResponse,
            element,
            endpoint,
            method,
            fromCache: true
        });

        // Cancel the real request + stop any later beforeRequest listeners
        e.preventDefault();
        e.stopImmediatePropagation();
    });

    document.addEventListener('hx:afterResponse', (e) => {
        // Skip the synthetic event we just fired ourselves
        if (e.detail.fromCache) return;

        const { element, response, endpoint, method } = e.detail;
        if (!element || !response) return;

        const shouldCache = element.getAttribute('hx-cache') === 'true';
        if (!shouldCache) return;

        const ttl = parseInt(element.getAttribute('hx-cache-ttl'), 10) || cacheManager.defaultTTL;
        const form = element.closest('form') || document.createElement('form');
        const resolvedEndpoint = endpoint || element.getAttribute('hx-get') || element.getAttribute('hx-post');
        const resolvedMethod = method || (element.hasAttribute('hx-post') ? 'POST' : 'GET');

        if (!resolvedEndpoint) return;

        const cacheKey = buildCacheKey(resolvedMethod, resolvedEndpoint, element, form);
        cacheManager.set(cacheKey, response, ttl);
    });

    document.addEventListener('hx:afterRequest', (e) => {
        const { element } = e.detail;
        if (!element) return;

        const invalidateAttr = element.getAttribute('hx-cache-invalidate');
        if (!invalidateAttr) return;

        invalidateAttr.split(',').forEach(raw => {
            const key = raw.trim();
            if (!key) return;

            // Exact match (in case someone puts a full key)
            cacheManager.invalidate(key);

            // Endpoint-style invalidation (the common case)
            ['GET', 'POST'].forEach(m => {
                // Match both "GET:/api/data" and "GET:/api/data:params..."
                cacheManager.invalidatePrefix(`${m}:${key}`);
            });
        });
    });

    // ---------- Public API ----------

    window.hx.cache = {
        get: (key) => cacheManager.get(key),
        set: (key, value, ttl) => cacheManager.set(key, value, ttl),
        invalidate: (key) => cacheManager.invalidate(key),
        invalidatePrefix: (prefix) => cacheManager.invalidatePrefix(prefix),
        clear: () => cacheManager.clear()
    };
});