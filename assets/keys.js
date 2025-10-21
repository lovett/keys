async function runTrigger(node) {
    let eventName = 'trigger:fail';
    let result = 'Could not connect to server';
    let status = 0;
    let state = '';

    try {
        const response = await fetch(node.href, { method: 'POST' });
        if (response.ok) {
            eventName = 'trigger:success';
            state = response.headers.get("X-Keys-State");
        }
        result = await response.text()
        status = response.status;
    } finally {
        const event = new CustomEvent(eventName, {
            detail: { node, status, result }
        });
        window.dispatchEvent(event);

        node.querySelector('.state').textContent = state;
    }
}

window.addEventListener('DOMContentLoaded', () => {
    let keyBuffer = '';
    let keyTimer;

    const saveButton = document.getElementById('save');
    if (saveButton) {
        saveButton.addEventListener('click', () => {
            document.querySelector('main.editor form').submit();
        });
    }

    window.addEventListener('keyup', (e) => {
        const formTags = ['INPUT', 'SELECT', 'TEXTAREA'];
        const specials = ['Shift', 'Alt', 'Meta', 'Control', 'Backspace'];
        if (formTags.indexOf(e.target.nodeName) > -1) {
            return;
        }

        if (specials.indexOf(e.key) > -1) {
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

        keyBuffer += e.key;
        clearTimeout(keyTimer);

        const keyNodes = document.querySelectorAll(`a.key[data-keypress^='${keyBuffer}']`);
        if (keyNodes.length === 0) {
            keyBuffer = '';
            return;
        }

        if (keyNodes.length === 1) {
            keyNodes[0].click();
            keyBufer = '';
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

    window.addEventListener('trigger:start', (e) => {
        setStatus('Running…', 'start');
    });

    window.addEventListener('trigger:success', (e) => {
        const node = e.detail.node;
        const message = e.detail.result ? e.detail.result : 'Done!';
        setStatus(message, 'success');
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
    const el = document.getElementById('status');
    if (!el) return;

    let icon = '';

    if (type === 'fail') {
        el.className = 'fail';
        icon = 'skull';
    }

    if (type === 'success') {
        el.className = 'success';
        icon = 'star';
    }

    if (type === 'start') {
        el.className = 'start';
        icon = 'wait';
    }


    el.innerHTML = `<svg class="icon"><use xlink:href="#icon-${icon}"></use></svg>
    <div>${message}</div>
    <a href="#" class="close"><svg class="icon"><use xlink:href="#icon-close"></use></svg></a>`;
}
