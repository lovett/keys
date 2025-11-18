let keyBuffer = '';
let keyTimer;

window.addEventListener('click', (e) => {
    if (e.target.nodeName !== 'A') return;

    if (e.target.classList.contains('close')) {
        e.preventDefault();
        window.dispatchEvent(new CustomEvent('app:clear'));
        return;
    }

    if (e.target.classList.contains('key')) {
        e.preventDefault();
        window.dispatchEvent(new CustomEvent('app:start'));

        // Give the start message some time to display and not flicker.
        setTimeout(() => runTrigger(e.target), 500);
        return;
    }
});

window.addEventListener('keyup', (e) => {
    if (e.key === 'Escape') {
        window.dispatchEvent(new CustomEvent('app:clear'));
        window.dispatchEvent(new CustomEvent('app:cancel'));
        return;
    }

    if (e.key === 'E') {
        const node = document.querySelector('form.edit');
        if (node) node.submit();
        return;
    }

    if (e.ctrlKey || e.shiftKey || e.altKey) return;

    if (['INPUT', 'SELECT', 'TEXTAREA'].indexOf(e.target.nodeName) > -1) return;

    if (['Shift', 'Alt', 'Meta', 'Control', 'Backspace'].indexOf(e.key) > -1) return;


    if (document.querySelector('#keys.locked')) {
        window.dispatchEvent(new CustomEvent('app:locked'));
        return;
    }

    keyBuffer += e.key;
    clearTimeout(keyTimer);

    const nodes = document.querySelectorAll(`a.key[data-keypress^='${keyBuffer}']`);
    if (nodes.length === 0) {
        keyBuffer = '';
        return;
    }

    if (nodes.length === 1 && nodes[0].dataset.keypress === keyBuffer) {
        nodes[0].click();
        keyBuffer = '';
        return;
    }

    keyTimer = setTimeout(() => {
        const key = document.querySelector(`a.key[data-keypress='${keyBuffer}']`);
        if (key) key.click();
        keyBuffer = '';
    }, 250);
});

window.addEventListener('app:cancel', () => {
    const node = document.querySelector('#cancel');
    if (!node) return;
    node.click();
});

window.addEventListener('app:clear', () => {
    const el = document.getElementById("status");
    if (!el) return;
    el.replaceChildren();
    el.className = '';
});

window.addEventListener('app:start', () => {
    setStatus('Running…', 'start');
});

window.addEventListener('app:locked', () => {
    setStatus('The keyboard is locked.', 'locked');
});

window.addEventListener('app:success', (e) => {
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
    const message = e.detail.status < 500 ? e.detail.result : 'Service Unavailable';
    setStatus(message, 'fail');
});

window.addEventListener('DOMContentLoaded', () => {
    const el = document.getElementById('save');
    if (el) {
        el.addEventListener('click', () => {
            document.querySelector('#editor form').submit();
        });
    }
});

function setStatus(message, type) {
    const container = document.getElementById('status');
    if (!container) return;

    container.className = type;

    const clone = document.querySelector('#status-message').content.cloneNode(true);
    clone.querySelector('.message').innerHTML = message;

    let icon = '';
    if (type === 'fail') icon = 'skull';
    if (type === 'success') icon = 'star';
    if (type === 'start') icon = 'wait';
    if (type === 'locked') icon = 'lock';
    clone.querySelector('.icon use').setAttribute('xlink:href', `#icon-${icon}`);

    container.replaceChildren(clone);
}

async function runTrigger(node) {
    let eventName = 'app:fail';
    let result = 'Could not connect to server';
    let status = 0;
    let state = '';
    let locked = false;

    try {
        const response = await fetch(node.href, { method: 'POST' });
        if (response.ok) {
            eventName = 'app:success';
            state = response.headers.get("X-Keys-State");
            locked = Boolean(Number.parseInt(response.headers.get("X-Keys-Locked"), 10) || 0);
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
            detail: { node, status, result, locked }
        });
        window.dispatchEvent(event);

        node.querySelector('.state').textContent = state;
    }
}
