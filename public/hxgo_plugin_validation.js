document.addEventListener('DOMContentLoaded', () => {
    const HX = window.hx.constants;
    let isInitialLoad = true;

    const validateLength = (element, showPopup = true) => {
        const minLength = element.getAttribute(HX.MIN) ? parseInt(element.getAttribute(HX.MIN), 10) : null;
        const maxLength = element.getAttribute(HX.MAX) ? parseInt(element.getAttribute(HX.MAX), 10) : null;
        const targetDisableId = element.getAttribute(HX.VDISABLE);
        const targetElement = targetDisableId ? document.querySelector(targetDisableId) : null;
        const valueLength = element.value.length;

        let isValid = true;

        if (minLength !== null && valueLength < minLength) {
            if (showPopup && !isInitialLoad) {
                const elementName = element.name ? element.name.charAt(0).toUpperCase() + element.name.slice(1) : "Element";
                window.hx.ui.showPopover(element, `${elementName}: Minimum length is ${minLength} characters.`);
            }
            isValid = false;
        }
        if (maxLength !== null && valueLength > maxLength) {
            if (showPopup && !isInitialLoad) {
                const elementName = element.name ? element.name.charAt(0).toUpperCase() + element.name.slice(1) : "Element";
                window.hx.ui.showPopover(element, `${elementName}: Maximum length is ${maxLength} characters.`);
            }
            isValid = false;
        }

        if (targetElement) targetElement.disabled = !isValid;
        return isValid;
    };

    const validatePasswords = (elementPw1, elementPw2, targetElement, showMatchPopup = false) => {
        const validLengthPw1 = validateLength(elementPw1, false);
        const validLengthPw2 = validateLength(elementPw2, false);
        const passwordsMatch = elementPw1.value === elementPw2.value;
        const bothFilled = elementPw1.value.length > 0 && elementPw2.value.length > 0;

        if (showMatchPopup && bothFilled && !passwordsMatch && !isInitialLoad) {
            window.hx.ui.showPopover(elementPw2, 'Passwords must match');
            targetElement.disabled = true;
            return false;
        }

        const isValid = validLengthPw1 && validLengthPw2 && passwordsMatch && bothFilled;
        targetElement.disabled = !isValid;
        return isValid;
    };

    const validateRequirements = (element, form) => {
        const requireAttr = element.getAttribute(HX.REQUIRE);
        if (requireAttr) {
            const requiredIds = requireAttr.split(',');
            for (let id of requiredIds) {
                const requiredElement = document.querySelector(id.trim());
                if (requiredElement && !requiredElement.value) {
                    if (!isInitialLoad) {
                        window.hx.ui.showPopover(requiredElement, 'This field is required.');
                    }
                    return false;
                }
            }
        }

        let allValid = true;
        const lengthElements = form.querySelectorAll(`[${HX.MIN}], [${HX.MAX}]`);
        lengthElements.forEach(el => {
            if (!validateLength(el, true)) allValid = false;
        });

        return allValid;
    };

    // ---------- Wiring (initial load + reattachment after swaps) ----------

    // Finds matches for `selector` within root, including root itself —
    // querySelectorAll only searches descendants, so root needs a separate check.
    function findInRoot(root, selector) {
        const matches = Array.from(root.querySelectorAll(selector));
        if (root.matches?.(selector)) matches.unshift(root);
        return matches;
    }

    function setupLengthValidation(root) {
        findInRoot(root, `[${HX.MIN}], [${HX.MAX}]`).forEach(element => {
            element.addEventListener('change', () => {
                validateLength(element);
                isInitialLoad = false;
            });
        });
    }

    function setupPasswordValidation(root) {
        // Scoped to this call only, so re-running for a swapped-in subtree
        // doesn't skip a group just because an earlier call already saw it.
        const processedGroups = new Set();

        findInRoot(root, `[${HX.PDISABLE}]`).forEach(element => {
            const pdisableValue = element.getAttribute(HX.PDISABLE);
            if (processedGroups.has(pdisableValue)) return;
            processedGroups.add(pdisableValue);

            const ids = pdisableValue.split(',').map(id => id.trim());
            if (ids.length !== 3) return;

            const [pw1Id, pw2Id, targetId] = ids;
            const pw1Element = document.querySelector(pw1Id);
            const pw2Element = document.querySelector(pw2Id);
            const targetElement = document.querySelector(targetId);

            if (!pw1Element || !pw2Element || !targetElement) return;

            targetElement.disabled = true;
            const validateBoth = (showPopup) => {
                validatePasswords(pw1Element, pw2Element, targetElement, showPopup);
                isInitialLoad = false;
            };

            pw1Element.addEventListener('input', () => validateBoth(false));
            pw2Element.addEventListener('input', () => validateBoth(false));
            pw1Element.addEventListener('blur', () => validateBoth(true));
            pw2Element.addEventListener('blur', () => validateBoth(true));
        });
    }

    function setupValidation(root) {
        setupLengthValidation(root);
        setupPasswordValidation(root);
    }

    // Hook into request lifecycle
    document.addEventListener('hx:beforeRequest', (e) => {
        const { element, form } = e.detail;
        if (element.hasAttribute(HX.POST) && !validateRequirements(element, form)) {
            e.preventDefault(); // Cancel request if validation fails
        }
    });

    document.addEventListener('hx:beforeUpload', (e) => {
        const { element, form } = e.detail;
        if (!validateRequirements(element, form)) {
            e.preventDefault(); // Cancel upload if validation fails
        }
    });

    // Re-bind validation on content swapped in after initial load — without
    // this, hx-min/hx-max/hx-pdisable on swapped-in forms are silently inert.
    document.addEventListener('hx:reattachListeners', (e) => {
        if (e.detail.element) setupValidation(e.detail.element);
    });

    // Initial validation setup
    document.addEventListener('hx:init', () => {
        setupValidation(document);
    });
});