// This content script runs on all pages. It can't be conditionally injected because we can't
// miss out on any events by the time the script gets injected. On the other hand manifest.json allows
// us to name content scripts to be run at 'document_start' which seemed like the ideal time.
(function () {
	// Used to style console.log messages.
	var messageStyle = 'color: blue; font-size: 15px;';
	var timeStyle = 'color: green; font-size: 13px';

	function getFormattedDate() {
		var d = new Date();
		return `${d.getFullYear()}-${d.getMonth() + 1}-${d.getDate()} ${d.getHours()}:${d.getMinutes()}:${d.getSeconds()}:${d.getMilliseconds()}`;
	}

	// console.log(`%cDocument start: %c${getFormattedDate()}`, messageStyle, timeStyle);
})();