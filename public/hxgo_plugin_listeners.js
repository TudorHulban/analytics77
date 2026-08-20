document.addEventListener('DOMContentLoaded', () => {
    const HX = window.hx.constants;

    function applyShowHide(targetSelector, timeToShow, cssTransitionClass, isShow) {
        const targetElement = document.querySelector(targetSelector);
        if (!targetElement) return;

        if (cssTransitionClass) targetElement.classList.add(cssTransitionClass);

        if (isShow) {
            targetElement.style.display = 'block';
            if (timeToShow > 0) {
                setTimeout(() => {
                    targetElement.style.display = 'none';
                    if (cssTransitionClass) targetElement.classList.remove(cssTransitionClass);
                }, timeToShow);
            }
        } else {
            targetElement.style.display = 'none';
            if (cssTransitionClass) targetElement.classList.remove(cssTransitionClass);
        }
    }

    function handleHxShowHide(element, attributeName) {
        const attributeValue = element.getAttribute(attributeName);
        if (!attributeValue) return;

        const targets = attributeValue.split(',').map(target => target.trim());
        targets.forEach(targetSpec => {
            const parts = targetSpec.split(':');
            if (parts.length >= 1) {
                const targetSelector = parts[0];
                const timeToShow = parts.length > 1 ? parseInt(parts[1], 10) : 0;
                const cssTransitionClass = parts.length > 2 ? parts[2] : '';
                applyShowHide(targetSelector, timeToShow, cssTransitionClass, attributeName === HX.SHOW);
            }
        });
    }

    function reattachEventListeners(element) {
        if (!element) return;

        // Show/Hide handlers
        if (element.matches(`[${HX.SHOW}], [${HX.HIDE}]`)) {
            handleHxShowHide(element, HX.SHOW);
            handleHxShowHide(element, HX.HIDE);
        }

        // Button handlers (delegated by core)
        if (element.matches(`button[${HX.GET}], button[${HX.POST}], button[${HX.UPLOAD}], a[${HX.GET}], a[${HX.POST}], a[${HX.UPLOAD}]`)) {
            // Core already handles these via delegation
        }

        // Select handlers (delegated by core)
        if (element.matches(`select[${HX.GET}], select[${HX.POST}]`)) {
            // Core already handles these via delegation
        }

        // Search clear handler
        if (element.id === 'items-search') {
            element.addEventListener('dblclick', handleDblClickClear);
        }

        // Recursive for children
        element.querySelectorAll(`[${HX.SHOW}], [${HX.HIDE}]`).forEach(el => {
            handleHxShowHide(el, HX.SHOW);
            handleHxShowHide(el, HX.HIDE);
        });

        const searchItem = element.querySelector('#items-search');
        if (searchItem) searchItem.addEventListener('dblclick', handleDblClickClear);
    }

    function handleDblClickClear(event) {
        if (event.target.tagName.toLowerCase() === 'button' || event.target.getElementsByTagName('button').length > 0) return;
        const searchItem = event.currentTarget;
        const elInputs = searchItem.querySelectorAll('input, select');
        for (let input of elInputs) {
            if (input.tagName.toLowerCase() === 'input') input.value = '';
            else if (input.tagName.toLowerCase() === 'select') input.selectedIndex = 0;
        }
    }

    // Handle reattachment after swaps
    document.addEventListener('hx:reattachListeners', (e) => {
        reattachEventListeners(e.detail.element);
    });

    // Initial setup
    document.addEventListener('hx:init', () => {
        document.querySelectorAll(`[${HX.SHOW}], [${HX.HIDE}]`).forEach(el => {
            handleHxShowHide(el, HX.SHOW);
            handleHxShowHide(el, HX.HIDE);
        });

        const searchItems = document.getElementById('items-search');
        if (searchItems) searchItems.addEventListener('dblclick', handleDblClickClear);
    });
});