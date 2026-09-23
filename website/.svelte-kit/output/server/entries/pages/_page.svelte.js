import "../../chunks/index-server.js";
import { a as head, b as escape_html, i as ensure_array_like, n as attr_style, r as derived, t as attr_class, y as attr } from "../../chunks/server.js";
//#region src/lib/content/site.ts
var githubUrl = "https://github.com/sentzunhat/human-ai-workflow-protocol";
var releaseLoadingVersion = "0.0.0";
var platformDownloads = [
	{
		label: "macOS Apple Silicon",
		asset: "hawp-darwin-arm64",
		icon: "⌘"
	},
	{
		label: "macOS Intel",
		asset: "hawp-darwin-amd64",
		icon: "⌘"
	},
	{
		label: "Linux x64",
		asset: "hawp-linux-amd64",
		icon: "◈"
	},
	{
		label: "Linux ARM64",
		asset: "hawp-linux-arm64",
		icon: "◈"
	},
	{
		label: "Windows x64",
		asset: "hawp-windows-amd64.exe",
		icon: "⊞"
	},
	{
		label: "Windows ARM64",
		asset: "hawp-windows-arm64.exe",
		icon: "⊞"
	}
];
function downloadUrl(asset) {
	return `${githubUrl}/releases/latest/download/${asset}`;
}
var navItems = [
	{
		label: "Protocol",
		href: "#protocol"
	},
	{
		label: "Agents",
		href: "#agents"
	},
	{
		label: "Benchmarks",
		href: "#benchmarks"
	},
	{
		label: "FAQ",
		href: "#faq"
	}
];
var problemItems = [
	[
		"01",
		"Every session starts cold.",
		"The goal and constraints get reconstructed from conversation instead of carried with the work."
	],
	[
		"02",
		"Execution drifts.",
		"Agents optimize locally while the original mission, boundaries, or expected output fade from context."
	],
	[
		"03",
		"Handoffs lose decisions.",
		"The next model, tool, or person re-derives what was already known and repeats avoidable work."
	]
];
var protocolFields = [
	[
		"01",
		"Input",
		"The request as received, before interpretation changes it."
	],
	[
		"02",
		"Context",
		"The smallest useful background the agent needs to reason correctly."
	],
	[
		"03",
		"Mission",
		"The concrete objective in one clear, execution-ready statement."
	],
	[
		"04",
		"Constraints",
		"Hard boundaries, quality bars, dependencies, and things that must not change."
	],
	[
		"05",
		"Output",
		"The observable end state that defines what “done” actually means."
	]
];
var providers = [
	"Claude Code",
	"GitHub Copilot",
	"Codex",
	"Cursor",
	"Continue"
];
var benchmarkMetrics = [
	["19%", "aggregate context-shaping token reduction across the 10-query benchmark, with dense cases saving up to 38%."],
	["95%", "downstream token reduction across the 9 successful v0.0.24 Ollama intake-shaping benchmark runs."],
	["<1 ms", "lexical search latency in the current kit benchmark; hybrid search measured 72 ms with 10/10 reported quality."]
];
var workflowSteps = [
	[
		"01 / SHAPE",
		"Lock intent before execution.",
		"Capture the request, relevant context, mission, constraints, and expected output as a work artifact."
	],
	[
		"02 / EXECUTE",
		"Give the agent bounded context.",
		"Search the kit and work index, then pack only the useful context instead of repeatedly loading everything."
	],
	[
		"03 / HAND OFF",
		"Carry the shape forward.",
		"The next agent or session reads the same portable state instead of reconstructing the project from conversation history."
	]
];
var faqs = [
	["What is HAWP?", "HAWP is an open-source Human-AI Workflow Protocol for shaping AI agent work before execution. It carries intent, context, mission, constraints, output, evidence, and handoffs in portable work artifacts."],
	["Is HAWP an AI agent framework?", "No. HAWP is a workflow protocol and CLI for shaping and carrying work. It does not replace your agent runtime, orchestrator, or model."],
	["Does HAWP require a database or cloud service?", "No. The work model is Markdown-first and repository-local. Search indexing can be generated locally, so the protocol remains a set of portable files rather than a hosted dependency."],
	["Why does HAWP use MCP?", "MCP gives connected agents a standard way to search the HAWP kit and work index without forcing the protocol into a specific editor or vendor."],
	["Can I use HAWP with more than one coding agent?", "Yes. Provider guides in the repository cover Claude Code, GitHub Copilot, Codex, Cursor, and Continue."]
];
//#endregion
//#region src/lib/components/ReleaseVersion.svelte
function ReleaseVersion($$renderer, $$props) {
	let { version, loading = false, variant = "inline" } = $$props;
	$$renderer.push(`<span${attr_class("release-version svelte-1hr2aq1", void 0, {
		"release-version-loading": loading,
		"release-version-badge": variant === "badge"
	})}><span class="release-version-text svelte-1hr2aq1">v${escape_html(version)}</span> `);
	if (loading) $$renderer.push(`<!--[0--><span class="release-skeleton-sheen svelte-1hr2aq1" aria-hidden="true"></span>`);
	else $$renderer.push("<!--[-1-->");
	$$renderer.push(`<!--]--></span>`);
}
//#endregion
//#region src/lib/components/Agents.svelte
function Agents($$renderer, $$props) {
	let { releaseVersion, releaseLoading } = $$props;
	$$renderer.push(`<section id="agents" class="border-b border-[var(--hawp-line)] py-24 sm:py-28"><div class="mx-auto grid max-w-[1280px] items-center gap-14 px-5 sm:px-8 lg:grid-cols-2 lg:px-10"><div class="overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-alt)] shadow-2xl shadow-black/20"><div class="flex gap-2 border-b border-[var(--hawp-line)] px-5 py-4"><i class="size-2.5 rounded-full bg-[var(--hawp-red)]"></i><i class="size-2.5 rounded-full bg-[var(--hawp-line-strong)]"></i><i class="size-2.5 rounded-full bg-[var(--hawp-muted-2)]"></i></div> <div class="p-6 font-mono text-sm leading-7 text-[var(--hawp-muted)] sm:p-8"><div><span class="text-[var(--hawp-clay)]">$</span> hawp work new</div> <div class="text-[var(--hawp-clay)]">✓ work item scaffolded</div> <div class="text-[var(--hawp-muted-2)]">  input       → captured</div> <div class="text-[var(--hawp-muted-2)]">  context     → bounded</div> <div class="text-[var(--hawp-muted-2)]">  mission     → explicit</div> <div class="text-[var(--hawp-muted-2)]">  constraints → preserved</div> <div class="text-[var(--hawp-muted-2)]">  output      → verifiable</div> <div class="mt-5"><span class="text-[var(--hawp-clay)]">$</span> hawp mcp</div> <div class="text-[var(--hawp-clay)]">✓ hawp_work_intake · compound shaping</div> <div class="text-[var(--hawp-clay)]">✓ hawp_work_doc · canonical work paths</div></div></div> <div><p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">Tool-independent by design</p> <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">The work survives the agent.</h2> <p class="mt-5 text-lg leading-8 text-[var(--hawp-muted)]">HAWP uses Markdown artifacts and exposes search and workflows through MCP, so the same project shape can move across coding agents without getting trapped in one vendor's memory or chat history.</p> <p class="mt-4 text-sm leading-6 text-[var(--hawp-muted-2)]">In `);
	ReleaseVersion($$renderer, {
		version: releaseVersion,
		loading: releaseLoading
	});
	$$renderer.push(`<!---->, compound intake can retrieve relevant context, reshape the request through a local LLM, and return structured mission, constraints, output specification, and done signal states without creating a work item automatically.</p> <div class="mt-7 flex flex-wrap gap-2"><!--[-->`);
	const each_array = ensure_array_like(providers);
	for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
		let provider = each_array[$$index];
		$$renderer.push(`<span class="rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-3 py-1.5 text-xs text-[var(--hawp-muted)]">${escape_html(provider)}</span>`);
	}
	$$renderer.push(`<!--]--></div></div></div></section>`);
}
//#endregion
//#region src/lib/components/Benchmarks.svelte
function Benchmarks($$renderer, $$props) {
	let { releaseVersion, releaseLoading } = $$props;
	$$renderer.push(`<section id="benchmarks" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"><div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10"><p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">`);
	ReleaseVersion($$renderer, {
		version: releaseVersion,
		loading: releaseLoading
	});
	$$renderer.push(`<!----> repository evidence</p> <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Less context. Faster retrieval. Explicit uncertainty.</h2> <p class="mt-5 max-w-3xl text-lg leading-8 text-[var(--hawp-muted)]">Current repository evidence covers search retrieval, context shaping, and structured intake shaping. These are HAWP-local measurements intended to make workflow claims inspectable.</p> <div class="mt-12 grid gap-3 md:grid-cols-3"><!--[-->`);
	const each_array = ensure_array_like(benchmarkMetrics);
	for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
		let metric = each_array[$$index];
		$$renderer.push(`<article class="rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-7"><strong class="block text-4xl font-bold tracking-[-0.04em] text-[var(--hawp-peach)]">${escape_html(metric[0])}</strong> <span class="mt-4 block text-sm leading-6 text-[var(--hawp-muted-2)]">${escape_html(metric[1])}</span></article>`);
	}
	$$renderer.push(`<!--]--></div> <p class="mt-5 text-xs text-[var(--hawp-muted-2)]">Benchmarks are project-local measurements, not universal model performance claims.</p></div></section>`);
}
//#endregion
//#region src/lib/components/Faq.svelte
function Faq($$renderer) {
	$$renderer.push(`<section id="faq" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"><div class="mx-auto max-w-4xl px-5 sm:px-8"><p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">FAQ</p> <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Small protocol, deliberate boundaries.</h2> <div class="mt-10 divide-y divide-[var(--hawp-line)] border-y border-[var(--hawp-line)]"><!--[-->`);
	const each_array = ensure_array_like(faqs);
	for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
		let faq = each_array[$$index];
		$$renderer.push(`<article class="py-6"><h3 class="text-lg font-semibold text-[var(--hawp-ink)]">${escape_html(faq[0])}</h3> <p class="mt-2 text-sm leading-6 text-[var(--hawp-muted-2)]">${escape_html(faq[1])}</p></article>`);
	}
	$$renderer.push(`<!--]--></div></div></section>`);
}
//#endregion
//#region src/lib/components/FinalCta.svelte
function FinalCta($$renderer) {
	$$renderer.push(`<section class="py-24 sm:py-28"><div class="mx-auto max-w-[1180px] px-5 sm:px-8"><div class="site-grid overflow-hidden rounded-[2rem] border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-6 py-12 text-center shadow-2xl shadow-black/20 sm:px-12 sm:py-16"><h2 class="mx-auto max-w-3xl text-4xl font-bold tracking-[-0.04em] text-[var(--hawp-ink)] sm:text-5xl">Give AI agent work a shape that survives the session.</h2> <p class="mx-auto mt-5 max-w-2xl text-lg leading-8 text-[var(--hawp-muted)]">HAWP is open source, Apache-2.0 licensed, and designed to stay readable even when the agent changes.</p> <div class="mt-8 flex flex-wrap justify-center gap-3"><a${attr("href", githubUrl)} class="rounded-xl border border-[var(--hawp-red)] bg-[var(--hawp-red-deep)] px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-[var(--hawp-red)]">Read the repository ↗</a> <a${attr("href", `${githubUrl}#install--pick-your-agent`)} class="rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] px-6 py-3.5 text-sm font-semibold text-[var(--hawp-ink)]">Get started</a></div></div></div></section>`);
}
//#endregion
//#region src/lib/components/Footer.svelte
function Footer($$renderer) {
	$$renderer.push(`<footer class="border-t border-[var(--hawp-line)] py-10"><div class="mx-auto flex max-w-[1280px] flex-col gap-5 px-5 text-sm text-[var(--hawp-muted-2)] sm:px-8 md:flex-row md:items-center md:justify-between lg:px-10"><span>HAWP · Human-AI Workflow Protocol · Sentzunhat</span> <div class="flex gap-5"><a class="hover:text-[var(--hawp-ink)]"${attr("href", githubUrl)}>GitHub</a><a class="hover:text-[var(--hawp-ink)]"${attr("href", `${githubUrl}/blob/main/LICENSE`)}>Apache-2.0</a><a class="hover:text-[var(--hawp-ink)]"${attr("href", `${githubUrl}/discussions`)}>Discussions</a></div></div></footer>`);
}
//#endregion
//#region src/lib/components/ThemeToggle.svelte
function ThemeToggle($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		const options = [
			{
				value: "light",
				label: "Light"
			},
			{
				value: "dark",
				label: "Dark"
			},
			{
				value: "system",
				label: "System"
			}
		];
		let preference = "system";
		$$renderer.push(`<div class="flex items-center rounded-full border border-[var(--hawp-line)] bg-[var(--hawp-raised)]/90 p-1 shadow-sm backdrop-blur" aria-label="Theme preference"><!--[-->`);
		const each_array = ensure_array_like(options);
		for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
			let option = each_array[$$index];
			$$renderer.push(`<button type="button"${attr_class(`flex min-h-8 items-center gap-1.5 rounded-full px-2.5 text-xs font-semibold transition sm:px-3 ${preference === option.value ? "bg-[var(--hawp-panel-soft)] text-[var(--hawp-ink)] shadow-sm" : "text-[var(--hawp-muted-2)] hover:text-[var(--hawp-ink)]"}`)}${attr("aria-pressed", preference === option.value)}${attr("aria-label", `Use ${option.label.toLowerCase()} theme`)}${attr("title", `${option.label} theme`)}>`);
			if (option.value === "light") $$renderer.push(`<!--[0--><svg aria-hidden="true" viewBox="0 0 24 24" class="size-3.5" fill="none" stroke="currentColor" stroke-width="1.8"><circle cx="12" cy="12" r="3.25"></circle><path d="M12 2.5v2.2M12 19.3v2.2M4.7 4.7l1.55 1.55M17.75 17.75l1.55 1.55M2.5 12h2.2M19.3 12h2.2M4.7 19.3l1.55-1.55M17.75 6.25l1.55-1.55"></path></svg>`);
			else if (option.value === "dark") $$renderer.push(`<!--[1--><svg aria-hidden="true" viewBox="0 0 24 24" class="size-3.5" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M20.3 15.4A8.4 8.4 0 0 1 8.6 3.7 8.4 8.4 0 1 0 20.3 15.4Z"></path></svg>`);
			else $$renderer.push(`<!--[-1--><svg aria-hidden="true" viewBox="0 0 24 24" class="size-3.5" fill="none" stroke="currentColor" stroke-width="1.8"><rect x="3.5" y="4.5" width="17" height="12" rx="2"></rect><path d="M9 20h6M12 16.5V20"></path></svg>`);
			$$renderer.push(`<!--]--> <span class="hidden lg:inline">${escape_html(option.label)}</span></button>`);
		}
		$$renderer.push(`<!--]--></div>`);
	});
}
//#endregion
//#region src/lib/components/Header.svelte
function Header($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		$$renderer.push(`<header class="sticky top-0 z-50 border-b border-[var(--hawp-line)] bg-[var(--hawp-bg)]/92 backdrop-blur-xl"><div class="mx-auto flex h-[72px] max-w-[1280px] items-center justify-between px-4 sm:px-8 lg:px-10"><a href="#top" class="flex items-center gap-2.5" aria-label="HAWP home"><span class="relative grid size-9 shrink-0 place-items-center rounded-[11px] border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel)]"><span class="size-4 rotate-45 border border-[var(--hawp-peach)]"></span> <span class="absolute size-2 rotate-45 bg-[var(--hawp-red)]"></span></span> <span class="text-xl font-bold tracking-tight text-[var(--hawp-ink)]">HAWP</span></a> <nav class="hidden items-center gap-8 md:flex" aria-label="Main navigation"><!--[-->`);
		const each_array = ensure_array_like([...navItems, {
			label: "Downloads",
			href: "#downloads"
		}]);
		for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
			let item = each_array[$$index];
			$$renderer.push(`<a${attr("href", item.href)} class="text-sm text-[var(--hawp-muted)] transition hover:text-[var(--hawp-ink)]">${escape_html(item.label)}</a>`);
		}
		$$renderer.push(`<!--]--></nav> <div class="flex items-center gap-1.5 sm:gap-3">`);
		ThemeToggle($$renderer, {});
		$$renderer.push(`<!----> <details class="group relative md:hidden"><summary class="flex min-h-10 cursor-pointer list-none items-center justify-center rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-3 text-sm font-semibold text-[var(--hawp-ink)] marker:content-none">Menu</summary> <nav class="absolute right-0 top-12 min-w-44 overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-2 shadow-2xl shadow-black/10" aria-label="Mobile navigation"><!--[-->`);
		const each_array_1 = ensure_array_like([...navItems, {
			label: "Downloads",
			href: "#downloads"
		}]);
		for (let $$index_1 = 0, $$length = each_array_1.length; $$index_1 < $$length; $$index_1++) {
			let item = each_array_1[$$index_1];
			$$renderer.push(`<a${attr("href", item.href)} class="block rounded-xl px-3 py-2.5 text-sm font-semibold text-[var(--hawp-muted)] transition hover:bg-[var(--hawp-panel)] hover:text-[var(--hawp-ink)]">${escape_html(item.label)}</a>`);
		}
		$$renderer.push(`<!--]--> <a${attr("href", githubUrl)} class="mt-1 block rounded-xl border-t border-[var(--hawp-line)] px-3 py-2.5 text-sm font-semibold text-[var(--hawp-ink)]" rel="noreferrer">GitHub ↗</a></nav></details> <a${attr("href", githubUrl)} class="hidden rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-4 py-2.5 text-sm font-semibold text-[var(--hawp-ink)] transition hover:border-[var(--hawp-red)] hover:bg-[var(--hawp-panel)] md:inline-flex" rel="noreferrer">GitHub</a></div></div></header>`);
	});
}
//#endregion
//#region src/lib/components/Hero.svelte
function Hero($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		let { releaseVersion, releaseLoading } = $$props;
		const nodePositions = [
			{
				x: 50,
				y: 10,
				label: "Input",
				detail: "request"
			},
			{
				x: 86,
				y: 39,
				label: "Context",
				detail: "background"
			},
			{
				x: 75,
				y: 88,
				label: "Mission",
				detail: "objective"
			},
			{
				x: 25,
				y: 88,
				label: "Constraints",
				detail: "limits"
			},
			{
				x: 14,
				y: 39,
				label: "Output",
				detail: "done state"
			}
		];
		let connectionPaths = [];
		let draggingIndex = -1;
		$$renderer.push(`<section id="top" class="site-grid relative overflow-hidden border-b border-[var(--hawp-line)]"><div class="pointer-events-none absolute -right-40 -top-52 size-[38rem] rounded-full bg-[var(--hawp-red-deep)]/16 blur-3xl"></div> <div class="mx-auto grid min-h-[710px] max-w-[1280px] items-center gap-14 px-5 py-20 sm:px-8 lg:grid-cols-[1.08fr_.92fr] lg:px-10 lg:py-24"><div class="relative z-10"><div class="mb-8 flex items-center gap-3 text-[11px] font-semibold tracking-[0.08em] text-[var(--hawp-peach)] uppercase"><span class="size-3 rounded-full bg-[var(--hawp-red)]"></span> `);
		ReleaseVersion($$renderer, {
			version: releaseVersion,
			loading: releaseLoading,
			variant: "badge"
		});
		$$renderer.push(`<!----> · Open protocol · Markdown-first · MCP-ready</div> <h1 class="max-w-[760px] text-5xl font-bold leading-[1.02] tracking-[-0.045em] text-[var(--hawp-ink)] sm:text-6xl lg:text-[5rem]">Shape the work before the model <span class="text-[var(--hawp-peach)]">touches it.</span></h1> <p class="mt-7 max-w-[730px] text-lg leading-7 text-[var(--hawp-muted)] sm:text-xl sm:leading-8">HAWP gives AI agent work a stable shape for intent, context, mission, constraints, and output, so agents can execute with less drift and hand off without starting from zero.</p> <div class="mt-8 flex flex-wrap gap-3"><a${attr("href", githubUrl)} class="rounded-xl border border-[var(--hawp-red)] bg-[var(--hawp-red-deep)] px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-[var(--hawp-red)]">Explore on GitHub ↗</a> <a href="#protocol" class="rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] px-6 py-3.5 text-sm font-semibold text-[var(--hawp-ink)] transition hover:border-[var(--hawp-line-strong)]">See the protocol ↓</a></div> <p class="mt-6 text-xs text-[var(--hawp-muted-2)]">Apache-2.0 · Go CLI · Static Markdown artifacts · No lock-in</p></div> <div class="protocol-visual" role="group" aria-label="Five HAWP fields converging into a stable work packet"><div class="orbit one"></div><div class="orbit two"></div><div class="orbit three"></div> <svg class="connections" viewBox="0 0 100 100" aria-hidden="true" focusable="false"><!--[-->`);
		const each_array = ensure_array_like(connectionPaths);
		for (let index = 0, $$length = each_array.length; index < $$length; index++) {
			let path = each_array[index];
			$$renderer.push(`<path class="connector-shadow"${attr("d", path)} pathLength="1"${attr_style(`--connector-delay: ${index * -.35}s`)}></path><path class="connector"${attr("d", path)} pathLength="1"${attr_style(`--connector-delay: ${index * -.35}s`)}></path>`);
		}
		$$renderer.push(`<!--]--></svg> <div class="core"><span>WORK</span></div> <!--[-->`);
		const each_array_1 = ensure_array_like(nodePositions);
		for (let index = 0, $$length = each_array_1.length; index < $$length; index++) {
			let node = each_array_1[index];
			$$renderer.push(`<button type="button"${attr_class("node", void 0, { "dragging": draggingIndex === index })}${attr_style(`--node-x: ${node.x}%; --node-y: ${node.y}%; --float-delay: ${index * -.7}s`)} aria-describedby="drag-instruction"><span class="node-content"><strong>${escape_html(node.label)}</strong>${escape_html(node.detail)}</span></button>`);
		}
		$$renderer.push(`<!--]--> <p id="drag-instruction" class="drag-hint">Drag a field to reshape the work packet</p></div></div></section>`);
	});
}
//#endregion
//#region src/lib/components/Problem.svelte
function Problem($$renderer) {
	$$renderer.push(`<section class="border-b border-[var(--hawp-line)] py-24 sm:py-28"><div class="mx-auto grid max-w-[1280px] gap-14 px-5 sm:px-8 lg:grid-cols-[.92fr_1.08fr] lg:px-10"><div><p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">The failure mode</p> <h2 class="mt-5 max-w-xl text-4xl font-bold leading-[1.08] tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">AI work often loses shape between the idea and execution.</h2> <p class="mt-6 max-w-xl text-lg leading-8 text-[var(--hawp-muted)]">The expensive part is not always generation. It is re-explaining intent, rediscovering constraints, recovering decisions, and correcting drift across sessions and tools.</p></div> <div class="divide-y divide-[var(--hawp-line)] overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)]"><!--[-->`);
	const each_array = ensure_array_like(problemItems);
	for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
		let item = each_array[$$index];
		$$renderer.push(`<article class="grid grid-cols-[46px_1fr] gap-4 p-6 sm:p-7"><span class="font-mono text-sm text-[var(--hawp-clay)]">${escape_html(item[0])}</span> <div><h3 class="text-xl font-semibold text-[var(--hawp-ink)]">${escape_html(item[1])}</h3><p class="mt-2 text-sm leading-6 text-[var(--hawp-muted-2)]">${escape_html(item[2])}</p></div></article>`);
	}
	$$renderer.push(`<!--]--></div></div></section>`);
}
//#endregion
//#region src/lib/components/Protocol.svelte
function Protocol($$renderer) {
	$$renderer.push(`<section id="protocol" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"><div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10"><p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">The protocol</p> <h2 class="mt-5 max-w-3xl text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Five fields. One portable shape for the work.</h2> <p class="mt-5 max-w-3xl text-lg leading-8 text-[var(--hawp-muted)]">HAWP keeps the protocol intentionally small. The artifact remains readable by humans, useful to agents, and portable across tools.</p> <div class="mt-12 grid gap-3 md:grid-cols-2 lg:grid-cols-5"><!--[-->`);
	const each_array = ensure_array_like(protocolFields);
	for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
		let field = each_array[$$index];
		$$renderer.push(`<article class="min-h-56 rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-6"><span class="font-mono text-xs text-[var(--hawp-clay)]">${escape_html(field[0])}</span> <h3 class="mt-8 text-xl font-semibold text-[var(--hawp-ink)]">${escape_html(field[1])}</h3> <p class="mt-3 text-sm leading-6 text-[var(--hawp-muted-2)]">${escape_html(field[2])}</p></article>`);
	}
	$$renderer.push(`<!--]--></div></div></section>`);
}
//#endregion
//#region src/lib/components/ReleaseDownloads.svelte
function ReleaseDownloads($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		let releaseLoading = true;
		function assetUrl(asset) {
			return downloadUrl(asset);
		}
		let releaseVersion = derived(() => "0.0.0");
		$$renderer.push(`<section id="downloads" class="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"><div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10"><div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between"><div><p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">Latest release</p> <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Get HAWP for your platform.</h2> <p class="mt-5 max-w-2xl text-lg leading-8 text-[var(--hawp-muted)]">Download the HAWP CLI from the latest GitHub release. The version and asset links update from the release feed when a new build is published.</p></div> <a class="shrink-0 text-sm font-semibold text-[var(--hawp-clay)] underline decoration-[var(--hawp-line-strong)] underline-offset-4 hover:text-[var(--hawp-ink)]"${attr("href", `https://github.com/sentzunhat/human-ai-workflow-protocol/releases`)}>View release notes ↗</a></div> <p class="sr-only" aria-live="polite">${escape_html("Checking GitHub for the latest HAWP release.")}</p> <div class="mt-10 grid gap-3 sm:grid-cols-2 lg:grid-cols-3"><!--[-->`);
		const each_array = ensure_array_like(platformDownloads);
		for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
			let platform = each_array[$$index];
			$$renderer.push(`<a class="flex items-center justify-between rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-5 transition hover:border-[var(--hawp-red)] hover:bg-[var(--hawp-panel)]"${attr("href", assetUrl(platform.asset))}><span class="flex items-center gap-3"><span class="grid size-9 place-items-center rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] text-sm text-[var(--hawp-clay)]" aria-hidden="true">${escape_html(platform.icon)}</span> <span><strong class="block text-sm text-[var(--hawp-ink)]">${escape_html(platform.label)}</strong> <span class="mt-1 block text-xs text-[var(--hawp-muted-2)]">`);
			ReleaseVersion($$renderer, {
				version: releaseVersion(),
				loading: releaseLoading
			});
			$$renderer.push(`<!----></span></span></span> <span class="text-lg text-[var(--hawp-clay)]" aria-hidden="true">↓</span></a>`);
		}
		$$renderer.push(`<!--]--></div></div></section>`);
	});
}
//#endregion
//#region src/lib/components/Workflow.svelte
function Workflow($$renderer) {
	$$renderer.push(`<section class="border-b border-[var(--hawp-line)] py-24 sm:py-28"><div class="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10"><p class="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">How it fits</p> <h2 class="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">Shape, execute, and hand off.</h2> <div class="mt-12 grid gap-8 md:grid-cols-3"><!--[-->`);
	const each_array = ensure_array_like(workflowSteps);
	for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
		let step = each_array[$$index];
		$$renderer.push(`<article class="border-l border-[var(--hawp-line-strong)] pl-5"><span class="font-mono text-xs text-[var(--hawp-clay)]">${escape_html(step[0])}</span> <h3 class="mt-5 text-xl font-semibold text-[var(--hawp-ink)]">${escape_html(step[1])}</h3> <p class="mt-3 text-sm leading-6 text-[var(--hawp-muted-2)]">${escape_html(step[2])}</p></article>`);
	}
	$$renderer.push(`<!--]--></div></div></section>`);
}
//#endregion
//#region src/routes/+page.svelte
function _page($$renderer, $$props) {
	$$renderer.component(($$renderer) => {
		let releaseVersion = releaseLoadingVersion;
		let releaseLoading = true;
		head("1uha8ag", $$renderer, ($$renderer) => {
			$$renderer.title(($$renderer) => {
				$$renderer.push(`<title>HAWP: Human-AI Workflow Protocol for AI Agent Work</title>`);
			});
			$$renderer.push(`<meta name="description" content="HAWP is an open-source Human-AI Workflow Protocol for shaping AI agent work, preserving context, reducing drift, and improving handoffs across tools."/> <link rel="canonical" href="https://hawp.online/"/> <meta property="og:type" content="website"/> <meta property="og:site_name" content="HAWP"/> <meta property="og:title" content="HAWP: Shape AI agent work before execution"/> <meta property="og:description" content="An open-source workflow protocol for clearer human-AI execution, evidence, handoffs, and less drift."/> <meta property="og:url" content="https://hawp.online/"/> <meta property="og:image" content="https://hawp.online/social-card.png"/> <meta property="og:image:type" content="image/png"/> <meta property="og:image:width" content="1200"/> <meta property="og:image:height" content="630"/> <meta name="twitter:card" content="summary_large_image"/> <meta name="twitter:title" content="HAWP: Human-AI Workflow Protocol"/> <meta name="twitter:description" content="Shape AI agent work before execution. Preserve intent, constraints, evidence, and handoffs."/> <meta name="twitter:image" content="https://hawp.online/social-card.png"/> <meta name="application-name" content="HAWP"/> <meta name="author" content="Sentzunhat"/> <meta name="subject" content="Human-AI workflow protocol and AI agent work"/>`);
		});
		Header($$renderer, {});
		$$renderer.push(`<!----> <main>`);
		Hero($$renderer, {
			releaseVersion,
			releaseLoading
		});
		$$renderer.push(`<!----> `);
		Problem($$renderer, {});
		$$renderer.push(`<!----> `);
		Protocol($$renderer, {});
		$$renderer.push(`<!----> `);
		Agents($$renderer, {
			releaseVersion,
			releaseLoading
		});
		$$renderer.push(`<!----> `);
		Benchmarks($$renderer, {
			releaseVersion,
			releaseLoading
		});
		$$renderer.push(`<!----> `);
		Workflow($$renderer, {});
		$$renderer.push(`<!----> `);
		ReleaseDownloads($$renderer, {});
		$$renderer.push(`<!----> `);
		Faq($$renderer, {});
		$$renderer.push(`<!----> `);
		FinalCta($$renderer, {});
		$$renderer.push(`<!----></main> `);
		Footer($$renderer, {});
		$$renderer.push(`<!---->`);
	});
}
//#endregion
export { _page as default };
