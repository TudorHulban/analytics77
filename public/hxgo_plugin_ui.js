document.addEventListener('DOMContentLoaded', () => {
    const CONFIG = {
        MS_LOADING_INDICATOR_DELAY: 100,
        MS_DURATION_DISPLAY_POPOVER: 3000,
        MS_DURATION_DISPLAY_ERROR_ALERT: 7000
    };

    // Safety net only – must be longer than the core's request timeout.
    const BUTTON_SAFETY_MS = (window.hx?.config?.WS_REQUEST_TIMEOUT || 30000) + 5000;

    let loadingIndicatorTimeout;

    const loadingIndicator = document.createElement('div');
    loadingIndicator.className = 'loading-indicator';
    loadingIndicator.style.cssText = 'position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:10000;font-family:Material Symbols Outlined;font-size:48px;color:rgba(0,0,0,.8);padding:10px 20px;background-color:hsl(48,85%,26%);border-radius:5px;box-shadow:0 0 10px rgba(0,0,0,.5);display:none;';
    loadingIndicator.textContent = 'sync';
    document.body.appendChild(loadingIndicator);

    function showLoadingIndicator() {
        loadingIndicatorTimeout = setTimeout(() => {
            loadingIndicator.style.display = 'block';
        }, CONFIG.MS_LOADING_INDICATOR_DELAY);
    }

    function hideLoadingIndicator() {
        clearTimeout(loadingIndicatorTimeout);
        loadingIndicator.style.display = 'none';
    }

    function showPopover(element, message) {
        const popover = document.createElement('div');
        popover.className = 'popover';
        popover.textContent = message;
        document.body.appendChild(popover);

        const rect = element.getBoundingClientRect();
        popover.style.left = `${rect.left + window.scrollX}px`;
        popover.style.top = `${rect.bottom + window.scrollY}px`;
        popover.style.zIndex = '1000';

        setTimeout(() => {
            if (popover.parentNode) document.body.removeChild(popover);
        }, CONFIG.MS_DURATION_DISPLAY_POPOVER);
    }

    function showErrorAlert(message) {
        const errorAlertBox = document.createElement('div');
        errorAlertBox.className = 'error-alert';
        errorAlertBox.textContent = message;
        document.body.appendChild(errorAlertBox);

        setTimeout(() => {
            if (errorAlertBox.parentNode) document.body.removeChild(errorAlertBox);
        }, CONFIG.MS_DURATION_DISPLAY_ERROR_ALERT);
    }

    // ---------- Loading indicator ----------
    // Only show once a request/upload is actually underway — i.e. after the
    // beforeRequest/beforeUpload cancellation checkpoint has already cleared.
    // This makes the show/hide pairing independent of plugin listener order:
    // a cancelling plugin (validation, cache) simply prevents requestSent /
    // uploadStarted from ever firing, so there's nothing here to race.
    document.addEventListener('hx:requestSent', () => showLoadingIndicator());
    document.addEventListener('hx:uploadStarted', () => showLoadingIndicator());

    document.addEventListener('hx:afterResponse', (e) => {
        if (e.detail.pendingCount === 0) hideLoadingIndicator();
    });
    document.addEventListener('hx:requestCancelled', hideLoadingIndicator);
    document.addEventListener('hx:timeout', hideLoadingIndicator);
    document.addEventListener('hx:connectionError', hideLoadingIndicator);
    document.addEventListener('hx:uploadComplete', hideLoadingIndicator);
    document.addEventListener('hx:uploadCancelled', hideLoadingIndicator);
    document.addEventListener('hx:error', hideLoadingIndicator);

    // ---------- Button disable / re-enable (centralized) ----------
    // element → { safetyTimer }
    const pendingButtons = new Map();

    const REENABLE_EVENTS = [
        'hx:afterResponse',
        'hx:requestCancelled',
        'hx:timeout',
        'hx:connectionError',
        'hx:uploadComplete',
        'hx:uploadCancelled',
        'hx:error'
    ];

    function reenableButton(element) {
        if (!element || !pendingButtons.has(element)) return;
        element.disabled = false;
        clearTimeout(pendingButtons.get(element).safetyTimer);
        pendingButtons.delete(element);
    }

    // Single set of permanent listeners
    REENABLE_EVENTS.forEach(name => {
        document.addEventListener(name, (ev) => {
            const el = ev.detail?.element;
            if (el) reenableButton(el);
        });
    });

    function disableButton(element) {
        if (!element) return;

        element.disabled = true;

        // Replace any previous entry for this element
        if (pendingButtons.has(element)) {
            clearTimeout(pendingButtons.get(element).safetyTimer);
        }

        const safetyTimer = setTimeout(() => {
            reenableButton(element);
        }, BUTTON_SAFETY_MS);

        pendingButtons.set(element, { safetyTimer });
    }

    // Disable only once a request/upload is actually underway (same
    // post-checkpoint events as the loading indicator, same reasoning).
    document.addEventListener('hx:requestSent', (e) => disableButton(e.detail.element));
    document.addEventListener('hx:uploadStarted', (e) => disableButton(e.detail.element));

    // ---------- Error handling ----------
    document.addEventListener('hx:error', (e) => {
        if (e.detail.message) showErrorAlert(e.detail.message);
        else if (e.detail.error) showErrorAlert('Error: ' + e.detail.error.message);
    });

    document.addEventListener('hx:timeout', () => {
        showErrorAlert('Request timed out');
    });

    document.addEventListener('hx:connectionError', () => {
        showErrorAlert('Not connected to server. Reconnecting...');
    });

    // ---------- Public API ----------
    window.hx.ui = {
        showPopover,
        showErrorAlert,
        showLoadingIndicator,
        hideLoadingIndicator
    };
});