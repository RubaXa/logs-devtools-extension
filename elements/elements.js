import {toCamelCase} from '../utils.js';

const R_HAS_EXPR = /\$\{[^}]+\}/;
const R_EXPR = /\$\{([^}]+)\}/g;
const R_ATTRS = /\battrs\.([a-zA-Z-]+)/g;
const R_THIS = /\bthis\./g;

const BOOL_ATTRS = {
	disabled: true,
	checked: true,
	readonly: true,
	autofocus: true,
};

const ATTRS_PROPS = {
	class: 'className',
	disabled: 'disabled',
	readonly: 'readOnly',
	autofocus: 'autoFocus',
	nodeValue: 'nodeValue',
	innerHTML: 'innerHTML',
};

const DEF_FACTORY = ((XElement) => class extends XElement {});

let __XElements__ = {};

function getImportDocument() {
	const doc = window['__importDocument__']
	window['__importDocument__'] = null
	return doc;
}

export function define(factory = DEF_FACTORY, options = {}) {
	const doc = getImportDocument();
	const view = createTemplate(doc.querySelector('template'));
	const XElement = createXClass(view, options);
	const name = doc.baseURI.match(/([a-z0-9-]+)\.html/)[1];
	const XElementClass = factory(XElement);

	XElementClass.meta = {
		name,
		view,
		options,
		document: doc,
		class: XElementClass,
	};

	__XElements__[name] = XElementClass;
	customElements.define(name, XElementClass);
	console.info('xdefine:', XElementClass.meta);
}

function createXClass(template, options) {
	const HTMLElement = options.HTMLElement || window.HTMLElement;
	const XElement = class XElement extends HTMLElement {
		constructor() {
			super();
			this.refs = {};
		}

		getShadowRoot() {
			if (!this.__shadowRoot) {
				this.__shadowRoot = this.attachShadow({mode: 'open'});
			}

			return this.__shadowRoot;
		}

		connectedCallback() {
			if (!this.view) {
				const shadowRoot = this.getShadowRoot();
				this.view = template.factory(this).appendTo(shadowRoot);
			}
		}

		attributeChangedCallback(attr) {
			if (this.view) {
				this.view.update();
			}
		}
	}

	Object.defineProperty(XElement, 'observedAttributes', {
		get: () => template.observedAttributes,
	});

	template.observedAttributes.forEach(name => {
		const prop = toCamelCase(name);

		if (XElement.prototype[prop]) {
			return;
		}

		Object.defineProperty(XElement.prototype, prop, {
			get() {
				return BOOL_ATTRS.hasOwnProperty(name)
					? this.hasAttribute(name)
					: this.getAttribute(name)
				;
			},

			set(value) {
				if (BOOL_ATTRS.hasOwnProperty(name)) {
					this[value ? 'setAttribute' : 'removeAttribute'](name, true);
				} else {
					this.setAttribute(name, value);
				}
			},
		});
	});

	return XElement;
}

function toExpr(value, observedAttributes) {
	return value.split(R_EXPR)
		.map((part, i) => {
			if (i % 2) {
				return part
					.replace(R_ATTRS, (_, name) => {
						!observedAttributes.includes(name) && observedAttributes.push(name);
						
						if (BOOL_ATTRS.hasOwnProperty(name)) {
							return `__this__.${name}`;
						} 
						
						return `TO_STR(__this__.${toCamelCase(name)})`;
					})
					.replace(R_THIS, '__this__.')
				;
			} else {
				return JSON.stringify(part)
			}

		})
		.filter(v => v !== '""')
		.join(' + ');
}

function TO_STR(v) {
	return (v == null) ? '' : v;
}

function createTemplate(template) {
	const setters = {};
	const content = template.content;

	let observedAttributes = [];

	if (template.hasAttribute('observed')) {
		observedAttributes = template.getAttribute('observed').split(' ');
	}

	function tryAdd(node, path, name, value) {
		if (R_HAS_EXPR.test(value) || name === 'ref') {
			const key = path.join('.');
			const list = setters[key] || (setters[key] = []);
			
			list.name = node.nodeName;
			list.push({
				name,
				expr: toExpr(value, observedAttributes),
			});
		}
	}
	
	(function next(node, path) {
		if (node) {
			if (node.nodeType === node.TEXT_NODE) {
				tryAdd(node, path, 'nodeValue', node.nodeValue);
			} else {
				[].forEach.call(node.attributes || [], (attr) => {
					tryAdd(node, path, attr.name, attr.value);
				});

				next(node.firstChild, path.concat('firstChild'));
			}

			next(node.nextSibling, path.concat('nextSibling'));
		}
	})(content, []);

	const testSource = [];
	const initialSource = [];
	const updateSource = []

	Object.keys(setters).forEach((key, i) => {
		const elem = `ctx._${i}`;

		initialSource.unshift(`${elem} = root.${key};`);
		
		const XElement = __XElements__[setters[key].name.toLowerCase()];
		let testIdx = -1;

		// if (XElement) {
		// 	const xoptions = XElement.meta.options;
		// 	if (xoptions.test) {
		// 		testIdx = testSource.indexOf(xoptions.test);
				
		// 		if (testIdx === -1) {
		// 			testIdx = testSource.push(xoptions.test) - 1;
		// 		}

		// 		updateSource.push(`if (__test${testIdx}(${elem})) {`);
		// 	}
		// }

		setters[key].forEach(({name, expr}) => {
			if (name === 'ref') {
				initialSource.push(`__this__.refs[${expr}] = ${elem};`)
			} else if (ATTRS_PROPS.hasOwnProperty(name)) {
				name = ATTRS_PROPS[name];
				updateSource.push(`${elem}.${name} = (${expr});`)
			} else {
				updateSource.push(`${elem}.setAttribute("${name}", ${expr});`)
			}
		});

		// if (testIdx > -1) {
		// 	updateSource.push('}');
		// }
	});
	
	const factory = Function('__this__, root, ctx, TO_STR', `
		"use strict";
		var lock;
		${testSource.map((test, i) => `var __test${i} = ${test.toString()};`).join('\n')}
		${initialSource.join('\n')}
		function update(__attr__) {
			lock = false;
			${updateSource.join('\n')}
		}
		update();
		return function () {
			if (lock) return;
			lock = true;
			requestAnimationFrame(update);
		};
	`);

	// console.log(factory.toString())

	return {
		observedAttributes,

		factory(elem) {
			let fragment = content.cloneNode(true);
			const update = factory(elem, fragment, {}, TO_STR);
	
			return {
				appendTo(el) {
					el.appendChild(fragment);
					return this;
				},
				update,
				fragment,
			};
		},
	}
}

export function detachChildren(root) {
	const children = document.createDocumentFragment();
		
	while (root.firstChild) {
		children.appendChild(root.firstChild);
	}

	return children;
}