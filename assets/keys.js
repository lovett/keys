async function runTrigger(node) {
    let eventName = 'trigger:fail';
    let result = 'Could not connect to server';
    let status = 0;

    try {
        const response = await fetch(node.href, {method: 'POST'});
        if (response.ok) {
            eventName = 'trigger:success';
            result = await response.text()
            status = response.status;
        }
    } finally {
        const event = new CustomEvent(eventName, {
            detail: {node, status, result }
        });
        window.dispatchEvent(event);
    }

}

window.addEventListener('DOMContentLoaded',  () => {
    const saveButton = document.getElementById('save');
    if (saveButton) {
        saveButton.addEventListener('click', () => {
            document.querySelector('main.editor form').submit();
        });
    }

    window.addEventListener('keyup', (e) => {
        if (e.key === 'Escape') {
            const el = document.getElementById("status");
            if (!el) return;
            el.innerHTML = '';
            el.className = '';
        }
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
            runTrigger(e.target);
        }
    });

    window.addEventListener('trigger:success', (e) => {
        const node = e.detail.node;

        let message = `Pressed ${node.dataset.keypress}`;
        if (e.detail.result) {
            message = e.detail.result;
        }

        setStatus(message, 'success');
    });

    window.addEventListener('trigger:fail', (e) => {
        let message = e.detail.result;

        switch (e.detail.status) {
        case 503:
            message = 'Service Unavailable';
            break;
        }

        setStatus(messge, 'fail');
    });
});

function setStatus(message, type) {
    const el = document.getElementById('status');
    if (!el) return;

    let icon = '';

    if (type === 'fail') {
        el.classList.add('fail');
        el.classList.remove('success');
        icon = 'skull';
    }

    if (type === 'success') {
        el.classList.add('success');
        el.classList.remove('fail');
        icon = 'star';
    }

    el.innerHTML = `<svg class="icon"><use xlink:href="#icon-${icon}"></use></svg>
    <div>${message}</div>
    <a href="#" class="close"><svg class="icon"><use xlink:href="#icon-close"></use></svg></a>`;
}
