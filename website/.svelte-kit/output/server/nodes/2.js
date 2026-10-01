

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/2.B5iuQP9w.js","_app/immutable/chunks/VT3rjbcm.js","_app/immutable/chunks/xihTtKlq.js","_app/immutable/chunks/DhLJT252.js"];
export const stylesheets = ["_app/immutable/assets/2.BjUMGSP7.css"];
export const fonts = [];
