/**
 * adds the extension ID to the event name so it's unique and matches with what
 * the host page fires.
 * @param  {String} name event name
 * @return {String}      new event name
 */
function getNamespacedEventName(name) {
	return chrome.runtime.id + '-' + name;
}

// window.addEventListener(getNamespacedEventName('object-changed'), (event) => {
// 	chrome.runtime.sendMessage({
// 		name: 'object-changed',
// 		changeList: event.detail
// 	});
// });

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
	// When panel closes, background page will tell content script to tell the host
	// page to do some clean-up
	if (message.name === 'clean-up') {
		window.dispatchEvent(new CustomEvent('clean-up'));
	}
});