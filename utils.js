import {cfg} from './config.js';

export async function call(method, query) {
	const qs = (query ? [].concat(query) : [])
		.map(v => {
			v = v.split('=');
			return `${v[0]}=${encodeURIComponent(v.slice(1))}`;
		})
		.join('&');

	const resp = await fetch(`http://${cfg.host}/${method}/?${qs}`);
	const json = await resp.json();
	return json;
}

const R_CAMELCASE = /-(.)/g;

function camelCaseReplacer(_, chr) {
	return chr.toCamelCase()
}

export function toCamelCase(v) {
	return v.replace(R_CAMELCASE, camelCaseReplacer);
}