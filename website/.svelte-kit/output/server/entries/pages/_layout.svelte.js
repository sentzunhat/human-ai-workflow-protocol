import { l as slot, p as escape_html } from "../../chunks/index-server.js";
//#region src/lib/components/AnalyticsConsent.svelte
function AnalyticsConsent($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		$$renderer.push("<!--[-1-->");
		$$renderer.push(`<!--]--> <button class="analytics-settings svelte-1x0139m" type="button" aria-label="Analytics settings">${escape_html("Analytics settings")}</button>`);
	});
}
//#endregion
//#region src/routes/+layout.svelte
function _layout($$renderer, $$props) {
	$$renderer.push(`<!--[-->`);
	slot($$renderer, $$props, "default", {}, null);
	$$renderer.push(`<!--]--> `);
	AnalyticsConsent($$renderer, {});
	$$renderer.push(`<!---->`);
}
//#endregion
export { _layout as default };
