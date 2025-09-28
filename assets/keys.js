async function runTrigger(node) {
    const response = await fetch(node.href, {method: 'POST'});

    let eventName = 'trigger:success';

    if (!response.ok) {
        eventName = 'trigger:fail';
    }

    const event = new CustomEvent(eventName, {detail: {node, result: await response.text()}});
    window.dispatchEvent(event);
}

window.addEventListener('DOMContentLoaded',  () => {
    const saveButton = document.getElementById('save');
    if (saveButton) {
        saveButton.addEventListener('click', () => {
            document.querySelector('main.editor form').submit();
        });
    }

    window.addEventListener('click', (e) => {
        if (e.target.nodeName !== 'A') return;

        if (e.target.classList.contains('key')) {
            e.preventDefault();
            runTrigger(e.target);
        }
    });

    window.addEventListener('trigger:success', (e) => {
        const node = e.detail.node;
        const status = document.getElementById('status');
        status.classList.add('success');
        status.classList.remove('fail');

        if (e.detail.result) {
            status.textContent = e.detail.result;
        } else {
            status.textContent = `Pressed ${node.dataset.keypress}`;
        }
    });

    window.addEventListener('trigger:fail', (e) => {
        const status = document.getElementById('status');
        status.classList.add('fail');
        status.classList.remove('success');
        status.textContent = e.detail.result;
    });
});
