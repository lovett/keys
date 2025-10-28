async function runTrigger(node) {
    let eventName = 'trigger:fail';
    let result = 'Could not connect to server';
    let status = 0;
    let state = '';
    let locked = false;

    try {
        const response = await fetch(node.href, { method: 'POST' });
        if (response.ok) {
            eventName = 'trigger:success';
            state = response.headers.get("X-Keys-State");
            locked = Boolean(parseInt(response.headers.get("X-Keys-Locked"), 10) || 0);
        }
        result = await response.text()
        status = response.status;
    } finally {
        const event = new CustomEvent(eventName, {
            detail: { node, status, result, locked }
        });
        window.dispatchEvent(event);

        node.querySelector('.state').textContent = state;
    }
}

window.addEventListener('DOMContentLoaded', () => {
    let keyBuffer = '';
    let keyTimer;
    const keyList = document.getElementById('keys')

    const saveButton = document.getElementById('save');
    if (saveButton) {
        saveButton.addEventListener('click', () => {
            document.querySelector('main.editor form').submit();
        });
    }

    window.addEventListener('keyup', (e) => {
        if (e.ctrlKey || e.shiftKey || e.altKey) return;

        if (['INPUT', 'SELECT', 'TEXTAREA'].indexOf(e.target.nodeName) > -1) {
            return;
        }

        if (['Shift', 'Alt', 'Meta', 'Control', 'Backspace'].indexOf(e.key) > -1) {
            return;
        }

        if (e.key === 'Escape') {
            const el = document.getElementById("status");
            if (el) {
                el.innerHTML = '';
                el.className = '';
            }
            return;
        }

        if (keyList && keyList.classList.contains('locked')) {
            window.dispatchEvent(new CustomEvent('trigger:locked'));
            return;
        }


        keyBuffer += e.key;
        clearTimeout(keyTimer);

        const keyNodes = document.querySelectorAll(`a.key[data-keypress^='${keyBuffer}']`);
        if (keyNodes.length === 0) {
            keyBuffer = '';
            return;
        }

        if (keyNodes.length === 1 && keyNodes[0].dataset.keypress === keyBuffer) {
            keyNodes[0].click();
            keyBuffer = '';
            return;
        }

        keyTimer = setTimeout(() => {
            const key = document.querySelector(`a.key[data-keypress='${keyBuffer}']`);
            if (key) key.click();
            keyBuffer = '';
        }, 250);
    });

    window.addEventListener('click', (e) => {
        if (e.target.nodeName !== 'A') return;

        if (e.target.classList.contains('close')) {
            const event = new KeyboardEvent('keyup', {
                key: "Escape",
                code: "Escape",
                bubbles: true,
            });
            window.dispatchEvent(event);
            return;
        }

        if (e.target.classList.contains('key')) {
            e.preventDefault();
            window.dispatchEvent(new CustomEvent('trigger:start'));

            // Give the trigger:start message some time to display and not flicker.
            setTimeout(() => runTrigger(e.target), 500);
        }
    });

    window.addEventListener('trigger:start', () => {
        setStatus('Running…', 'start');
    });

    window.addEventListener('trigger:locked', () => {
        setStatus('The keyboard is locked.', 'locked');
    });

    window.addEventListener('trigger:success', (e) => {
        const message = e.detail.result ? e.detail.result : 'Done!';
        setStatus(message, 'success');

        if (e.detail.command) {
            console.log(e.detail.command);
        }

        const configNode = document.getElementById('config-locked');
        if (configNode) {
            const className = 'hidden';
            if (e.detail.locked) {
                configNode.classList.remove(className);
            } else {
                configNode.classList.add(className);
            }
        }

        if (keyList) {
            const className = 'locked';
            if (e.detail.locked) {
                keyList.classList.add(className);
            } else {
                keyList.classList.remove(className);
            }
        }
    });

    window.addEventListener('trigger:fail', (e) => {
        let message = e.detail.result;

        if (e.detail.status === 503) {
            message = 'Service Unavailable';
        }

        setStatus(message, 'fail');
    });
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
