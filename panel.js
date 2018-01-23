import {cfg} from './config.js';
import {call} from './utils.js';

// MAIN
(async () => {
	const canvas = document.getElementById('panel-canvas')
	
	try {
		await boot({
			id: chrome.devtools.inspectedWindow.tabId,
			canvas
		});
	} catch (err) {
		canvas.innerHTML = `<pre>${err.stack}</pre>`;
	}
})();

async function boot({id, canvas}) {
	const version = await call('version')
	const devChannel = chrome.runtime.connect({
		name: 'logs-devtools',
	});

	canvas.innerHTML = 'Interactive';
	// devChannel.onMessage.addListener((message, sender, sendResponse) => {
	// 	canvas.innerHTML = `v${version}:${JSON.stringify(message, null, 2)}`;
	// });

	const socket = await connect({
		id,
		onopen() {
			canvas.innerHTML = `Connected`;
		},
		onmessage(data) {
			canvas.innerHTML += `<br/>v${version}:${JSON.stringify(data, null, 2)}`;
		},
	});
	const logs = await setup(id);
	canvas.innerHTML += `<br/>v${version}:${JSON.stringify(logs, null, 2)}`;
}

async function connect({id, onopen, onmessage}) {
	return new Promise((resolve, reject) => {
		const socket = new WebSocket(`ws://${cfg.host}/ws/${id}`);
		
		socket.onopen = () => {
			onopen();
			resolve({
				postMessage() {
				}
			});
		};

		socket.onclose = () => {
			// todo: реконнект
		};

		socket.onmessage = ({data}) => {
			onmessage(data)
		};

		socket.onerror = ({message}) => {
			reject(message);
		};
	});
}

async function setup(id) {
	return await call('setup', [
		`id=${id}`,
		'log=10 /Users/k.lebedev/Developer/logs-devtools-extension/go-service/test_log.txt',
	]);
}