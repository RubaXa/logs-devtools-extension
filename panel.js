const devChannel = chrome.runtime.connect({
    name: 'logs-devtools',
});

const canvas = document.getElementById('panel-canvas');

canvas.innerHTML = 'Ready';

devChannel.onMessage.addListener((message, sender, sendResponse) => {
	canvas.innerHTML = `${JSON.stringify(message, null, 2)}`;
});