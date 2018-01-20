const ports = new Set();

function addPort(port) {
	port.onDisconnect.addListener(deletePort);
	ports.add(port);

}

function deletePort(port) {
	ports.delete(port);
}

chrome.runtime.onConnect.addListener((port) => {
	if (port.name == 'logs-devtools') {
		addPort(port);
	}
});


setInterval(() => {
	for (port of ports) {
		port.postMessage({
			date: new Date().toString(),
		});
	}
}, 1000);