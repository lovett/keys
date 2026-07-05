let keyBuffer = '';
let keyTimer = 0;

window.addEventListener('click', (e) => {
    const target = e.target;
    if (target instanceof HTMLAnchorElement === false) return;

    if (target.classList.contains('close')) {
        e.preventDefault();
        window.dispatchEvent(new CustomEvent('app:clear'));
    }

    if (target.classList.contains('key')) {
        e.preventDefault();
        window.dispatchEvent(new CustomEvent('app:start'));

        // Give the start message some time to display and not flicker.
        setTimeout(() => runTrigger(target), 500);
    }
});

window.addEventListener('keyup', (e) => {
    if (e.key === 'Escape') {
        window.dispatchEvent(new CustomEvent('app:clear'));
        window.dispatchEvent(new CustomEvent('app:cancel'));
        return;
    }

    if (e.ctrlKey || e.altKey) {
        return;
    }

    if (['Control', 'Alt', 'Tab'].indexOf(e.key) > -1) return;

    if (e.key === 'E') {
        const node = document.querySelector('form.edit');
        if (node instanceof HTMLFormElement) node.submit();
        return;
    }

    if (e.target instanceof HTMLElement) {
        const tag = e.target.nodeName;
        if (['INPUT', 'SELECT', 'TEXTAREA'].indexOf(tag) > -1) return;
    }

    if (document.querySelector('#keys.locked')) return;

    keyBuffer += e.key;
    clearTimeout(keyTimer);

    const nodes = document.querySelectorAll(`a.key[data-keypress^='${keyBuffer}']`);
    if (nodes.length === 0) {
        keyBuffer = '';
        return;
    }

    if (nodes.length === 1 && nodes[0] instanceof HTMLAnchorElement && nodes[0].dataset.keypress === keyBuffer) {
        nodes[0].click();
        keyBuffer = '';
        return;
    }

    keyTimer = setTimeout(() => {
        const key = document.querySelector(`a.key[data-keypress='${keyBuffer}']`);
        if (key instanceof HTMLAnchorElement) key.click();
        keyBuffer = '';
    }, 250);
});

window.addEventListener('app:cancel', () => {
    const node = document.querySelector('#cancel');
    if (node instanceof HTMLAnchorElement) node.click();
});

window.addEventListener('app:clear', () => {
    const el = document.getElementById("status");
    if (el instanceof HTMLDivElement) {
        el.replaceChildren();
        el.className = '';
    }
});

window.addEventListener('app:start', () => {
    setStatus('Running…', 'start');
});

window.addEventListener('app:success', (e) => {
    if (e instanceof CustomEvent === false) return;
    const message = e.detail.result ? e.detail.result : 'Done!';
    setStatus(message, 'success');

    const configNode = document.querySelector('#config-locked');
    const keysNode = document.querySelector('#keysNode');
    if (e.detail.locked) {
        if (configNode) configNode.classList.remove('hidden');
        if (keysNode) keysNode.classList.add('locked');

    } else {
        if (configNode) configNode.classList.add('hidden');
        if (keysNode) keysNode.classList.remove('locked');
    }
});

window.addEventListener('app:fail', (e) => {
    if (e instanceof CustomEvent === false) return;
    const message = e.detail.status < 500 ? e.detail.result : 'Service Unavailable';
    setStatus(message, 'fail');
});

window.addEventListener('DOMContentLoaded', () => {
    const el = document.getElementById('save');
    if (el instanceof HTMLButtonElement === false) return;

    el.addEventListener('click', () => {
        const node = document.querySelector('#editor form')
        if (node instanceof HTMLFormElement) node.submit();
    });
});

/**
 * @param {string} message
 * @param {string} type
 */
function setStatus(message, type) {
    const container = document.getElementById('status');
    if (container instanceof HTMLDivElement === false) return;

    container.className = type;

    const statusEl = document.querySelector('#status-message')

    if (statusEl instanceof HTMLTemplateElement === false) return;

    container.replaceChildren(statusEl.content.cloneNode(true));
    const messageEl = container.querySelector('.message')
    if (messageEl instanceof HTMLElement === false) return;
    messageEl.innerHTML = message;

    let icon = '';
    if (type === 'fail') icon = 'skull';
    if (type === 'success') icon = 'star';
    if (type === 'start') icon = 'wait';

    container.querySelector('.icon use')?.setAttribute('xlink:href', `#icon-${icon}`);
}

/**
 * @param {HTMLAnchorElement} el
 */
async function runTrigger(el) {
    let eventName = 'app:fail';
    let result = 'Could not connect to server';
    let status = 0;
    let state = '';
    let locked = false;

    try {
        const response = await fetch(el.href, { method: 'POST' });
        if (response.ok) {
            eventName = 'app:success';
            state = response.headers.get("X-Keys-State") || "";
            locked = Boolean(Number.parseInt(response.headers.get("X-Keys-Locked") || "", 10) || 0);
        }

        if (response.headers.get("Content-Type") === "text/html") {
            const parser = new DOMParser()
            const doc = parser.parseFromString(await response.text(), "text/html")
            const body = doc.querySelector('body');
            result = (body) ? body.innerHTML : '<em>Response cannot be shown.</em>';
        } else {
            result = await response.text();
        }

        status = response.status;
    } finally {
        const event = new CustomEvent(eventName, {
            detail: { node: el, status, result, locked }
        });
        window.dispatchEvent(event);

        if (locked) {
            document.getElementById('keys')?.classList.add('locked');
            document.getElementById('config-locked')?.classList.add('locked');
        } else {
            document.getElementById('keys')?.classList.remove('locked');
            document.getElementById('config-locked')?.classList.remove('locked');
        }
        const stateEl = el.querySelector('.state');
        if (stateEl) {
            stateEl.textContent = state;
        }
    }
}
