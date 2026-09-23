
import root from '../root.js';
import { set_building, set_prerendering } from '$app/env/internal';
import { set_assets } from '$app/paths/internal/server';
import { set_manifest, set_read_implementation } from '__sveltekit/server';
import { set_private_env, set_public_env } from '../../../node_modules/@sveltejs/kit/src/runtime/shared-server.js';
import error from '../shared/error-template.js';

export const options = {
	app_template_contains_nonce: false,
	async: false,
	csp: {"mode":"auto","directives":{"upgrade-insecure-requests":false,"block-all-mixed-content":false},"reportOnly":{"upgrade-insecure-requests":false,"block-all-mixed-content":false}},
	csrf_check_origin: true,
	csrf_trusted_origins: [],
	embedded: false,
	env_public_prefix: 'PUBLIC_',
	env_private_prefix: '',
	hash_routing: false,
	hooks: null, // added lazily, via `get_hooks`
	preload_strategy: "modulepreload",
	root,
	service_worker: false,
	service_worker_options: undefined,
	server_error_boundaries: false,
	templates: {
		app: ({ head, body, assets, nonce, env }) => "<!doctype html>\n<html lang=\"en\">\n  <head>\n    <meta charset=\"utf-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n    <meta name=\"theme-color\" content=\"#151112\" />\n    <link rel=\"icon\" href=\"" + assets + "/favicon.svg\" />\n    <link rel=\"describedby\" href=\"" + assets + "/llms.txt\" />\n    <script>\n      (() => {\n        try {\n          const preference = localStorage.getItem('hawp-theme') || 'system';\n          const dark = preference === 'dark' || (preference === 'system' && matchMedia('(prefers-color-scheme: dark)').matches);\n          document.documentElement.dataset.preference = preference;\n          document.documentElement.dataset.theme = dark ? 'dark' : 'light';\n          document.querySelector('meta[name=\"theme-color\"]').content = dark ? '#151112' : '#f5eee8';\n        } catch {\n          const dark = matchMedia('(prefers-color-scheme: dark)').matches;\n          document.documentElement.dataset.preference = 'system';\n          document.documentElement.dataset.theme = dark ? 'dark' : 'light';\n        }\n      })();\n    </script>\n    <script type=\"application/ld+json\">\n      {\n        \"@context\": \"https://schema.org\",\n        \"@graph\": [\n          {\n            \"@type\": \"WebSite\",\n            \"name\": \"HAWP: Human-AI Workflow Protocol\",\n            \"url\": \"https://hawp.online/\",\n            \"description\": \"An open-source, Markdown-first workflow protocol and Go CLI for shaping AI agent work before execution.\",\n            \"inLanguage\": \"en\",\n            \"publisher\": {\n              \"@type\": \"Organization\",\n              \"name\": \"Sentzunhat\",\n              \"url\": \"https://github.com/sentzunhat\"\n            }\n          },\n          {\n            \"@type\": \"SoftwareSourceCode\",\n            \"name\": \"HAWP: Human-AI Workflow Protocol\",\n            \"description\": \"An open-source protocol and CLI for shaping AI agent work before execution, preserving context, reducing drift, and improving handoffs.\",\n            \"url\": \"https://hawp.online/\",\n            \"codeRepository\": \"https://github.com/sentzunhat/human-ai-workflow-protocol\",\n            \"license\": \"https://www.apache.org/licenses/LICENSE-2.0\",\n            \"programmingLanguage\": [\"Go\", \"Markdown\"],\n            \"keywords\": \"AI agents, human-AI workflow, context engineering, MCP, Markdown\"\n          },\n          {\n            \"@type\": \"WebPage\",\n            \"name\": \"HAWP: Human-AI Workflow Protocol for AI Agent Work\",\n            \"description\": \"HAWP is an open-source workflow protocol and Go CLI for shaping AI agent work, preserving context, reducing drift, and improving handoffs across tools.\",\n            \"url\": \"https://hawp.online/\",\n            \"isPartOf\": { \"@type\": \"WebSite\", \"url\": \"https://hawp.online/\" },\n            \"about\": [\n              { \"@type\": \"Thing\", \"name\": \"AI agent workflow\" },\n              { \"@type\": \"Thing\", \"name\": \"Context engineering\" },\n              { \"@type\": \"Thing\", \"name\": \"Human-AI collaboration\" }\n            ]\n          },\n          {\n            \"@type\": \"FAQPage\",\n            \"mainEntity\": [\n              {\n                \"@type\": \"Question\",\n                \"name\": \"What is HAWP?\",\n                \"acceptedAnswer\": { \"@type\": \"Answer\", \"text\": \"HAWP is an open-source Human-AI Workflow Protocol for shaping AI agent work before execution. It carries intent, context, mission, constraints, output, evidence, and handoffs in portable work artifacts.\" }\n              },\n              {\n                \"@type\": \"Question\",\n                \"name\": \"Is HAWP an AI agent framework?\",\n                \"acceptedAnswer\": { \"@type\": \"Answer\", \"text\": \"No. HAWP is a workflow protocol and CLI for shaping and carrying work. It does not replace your agent runtime, orchestrator, or model.\" }\n              },\n              {\n                \"@type\": \"Question\",\n                \"name\": \"Does HAWP require a database or cloud service?\",\n                \"acceptedAnswer\": { \"@type\": \"Answer\", \"text\": \"No. The work model is Markdown-first and repository-local. Search indexing can be generated locally, so the protocol remains a set of portable files rather than a hosted dependency.\" }\n              },\n              {\n                \"@type\": \"Question\",\n                \"name\": \"Why does HAWP use MCP?\",\n                \"acceptedAnswer\": { \"@type\": \"Answer\", \"text\": \"MCP gives connected agents a standard way to search the HAWP kit and work index without forcing the protocol into a specific editor or vendor.\" }\n              },\n              {\n                \"@type\": \"Question\",\n                \"name\": \"Can I use HAWP with more than one coding agent?\",\n                \"acceptedAnswer\": { \"@type\": \"Answer\", \"text\": \"Yes. Provider guides in the repository cover Claude Code, GitHub Copilot, Codex, Cursor, and Continue.\" }\n              }\n            ]\n          }\n        ]\n      }\n    </script>\n    " + head + "\n  </head>\n  <body data-sveltekit-preload-data=\"hover\">\n    <div style=\"display: contents\">" + body + "</div>\n  </body>\n</html>\n",
		error
	},
	version_hash: "1vdakwl"
};

export async function get_hooks() {
	let handle;
	let handleFetch;
	let handleError;
	let handleValidationError;
	let init;
	

	let reroute;
	let transport;
	

	return {
		handle,
		handleFetch,
		handleError,
		handleValidationError,
		init,
		reroute,
		transport
	};
}

export { set_assets, set_building, set_manifest, set_prerendering, set_private_env, set_public_env, set_read_implementation };
